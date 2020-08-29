package chatservice

import (
	"context"
	"fmt"
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
	"github.com/sirupsen/logrus"
	"sync"
)

type MessageRepository interface {
	GetAllMessage() ([]*chatpb.Message, error)
	AddNewMessage(message *chatpb.Message) (*chatpb.Message, error)
	UpdateMessage(message *chatpb.Message) (*chatpb.Message, error)
}

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

func NewChatService(repo MessageRepository) *ChatService {
	return &ChatService{
		repo:        repo,
		connections: nil}
}

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

func (c *ChatService) UpdateMessage(_ context.Context, _ *chatpb.UpdateMessageRequest) (*chatpb.EmptyResponse, error) {
	panic("implement me")
}

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
