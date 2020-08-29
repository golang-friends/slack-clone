package chatservice

import (
	"context"
	"fmt"
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
	"github.com/sirupsen/logrus"
	"sync"
)

// MessageRepository represents Repository that ChatService expects.
// Each repository should implement the following methods.
type MessageRepository interface {
	GetAllMessages() ([]*chatpb.Message, error)
	AddNewMessage(message *chatpb.Message) (*chatpb.Message, error)
	UpdateMessage(message *chatpb.Message) (*chatpb.Message, error)
}

// ChatService is the implementation of gRPC ChatService.
type ChatService struct {
	repo MessageRepository

	// TODO: We need additional event driven architecture if we need more than
	// 1 instance of chat service.
	//
	// For example,
	//
	//   PostMessage in 1 server
	//   -> create a MessageCreated event
	//   -> pub/sub service
	//   -> flows to other instances.
	connections []*Connection

	sync.Mutex
}

// NewChatService returns a new ChatService.
func NewChatService(repo MessageRepository) *ChatService {
	return &ChatService{
		repo:        repo,
		connections: nil}
}

// PostMessage is invoked when a user sends a message. Hence, it will create a new message in the repository.
// In addition, it will broadcast the message to other connected users.
func (c *ChatService) PostMessage(_ context.Context, request *chatpb.PostMessageRequest) (*chatpb.EmptyResponse, error) {
	// TODO: Verify the token and retrieve the token
	newMessage := request.GetMessage()
	if newMessage == nil {
		return nil, fmt.Errorf("no message was received")
	}
	newMessage, err := c.repo.AddNewMessage(newMessage)
	if err != nil {
		return nil, err
	}
	go broadcastMessageToOtherConnections(c, newMessage, newMessage.GetWorkspaceId())
	return &chatpb.EmptyResponse{}, nil
}

// UpdateMessage is invoked when a user modifies the message.
func (c *ChatService) UpdateMessage(_ context.Context, _ *chatpb.UpdateMessageRequest) (*chatpb.EmptyResponse, error) {
	panic("implement me")
}

// Subscribe is the first function to be invoked when a user joins the workspace or the room.
// It will then store the connection and then later will be invoked when broadcasting messages.
func (c *ChatService) Subscribe(request *chatpb.SubscribeRequest, server chatpb.ChatService_SubscribeServer) error {
	// TODO: Verify the token
	workspaceID := request.GetWorkspaceId()
	go addNewConnection(c, server, workspaceID)
	return nil
}

func addNewConnection(c *ChatService, server chatpb.ChatService_SubscribeServer, workspaceID string) {
	c.Lock()
	defer c.Unlock()
	newConnection := &Connection{
		server: &server,
		meta: ConnectionMeta{
			workspaceID: workspaceID,
		},
	}
	c.connections = append(c.connections, newConnection)
}

func broadcastMessageToOtherConnections(c *ChatService, newMessage *chatpb.Message, workspaceID string) {
	c.Lock()
	defer c.Unlock()

	var validConnections []*Connection

	for _, conn := range c.connections {
		if conn.server == nil || conn.meta.workspaceID != workspaceID {
			logrus.Info("skipping connection either (server is nil) or (workspace is not matched)")
			continue
		}
		err := (*conn.server).Send(&chatpb.SubscribeResponse{Message: newMessage})
		if err != nil {
			logrus.WithField("message", newMessage).WithError(err).Warn(
				"failed to send message")
			continue
		}
		validConnections = append(validConnections, conn)
	}

	c.connections = validConnections
}
