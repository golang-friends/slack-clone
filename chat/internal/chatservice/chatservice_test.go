package chatservice

import (
	"context"
	"github.com/golang-friends/slack-clone/chat/internal/repo/inmemoryrepo"
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"sync"
	"testing"
	"time"
)

func TestChatService_PostMessage(t *testing.T) {
	server := newTestServer()

	_, err := server.PostMessage(context.Background(), &chatpb.PostMessageRequest{
		Token: "",
		Message: &chatpb.Message{
			UserId:      "kkweon",
			WorkspaceId: "workspace-id",
			ChannelId:   "channel-id",
			Message:     "hello world",
		},
	})

	assert.NoError(t, err)
}

func TestChatService_Subscribe(t *testing.T) {
	server := newTestServer()

	mocksub := &chatsubscribeServer{}

	mocksub.wg.Add(1)

	// in order to not to way forever
	go func() {
		time.Sleep(time.Second * 5)
		mocksub.wg.Done()
	}()

	// Call subscribe to the server
	assert.NoError(t, server.Subscribe(&chatpb.SubscribeRequest{
		Token:       "token",
		WorkspaceId: "workspace-id",
	}, mocksub))

	// Send a message
	_, err := server.PostMessage(context.Background(), &chatpb.PostMessageRequest{
		Token: "",
		Message: &chatpb.Message{
			UserId:      "kkweon",
			WorkspaceId: "workspace-id",
			ChannelId:   "channel-id",
			Message:     "hello world",
		},
	})
	assert.NoError(t, err)

	mocksub.wg.Wait()

	actualMessage := mocksub.responses[0]
	assert.NotEmpty(t, actualMessage.GetMessage().GetMessageId())
	assert.Equal(t, "hello world", actualMessage.GetMessage().GetMessage())
}

type chatsubscribeServer struct {
	responses []*chatpb.SubscribeResponse
	wg        sync.WaitGroup
}

func (c *chatsubscribeServer) Send(response *chatpb.SubscribeResponse) error {
	defer c.wg.Done()
	c.responses = append(c.responses, response)
	return nil
}

func (c *chatsubscribeServer) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (c *chatsubscribeServer) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (c *chatsubscribeServer) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (c *chatsubscribeServer) Context() context.Context {
	panic("implement me")
}

func (c *chatsubscribeServer) SendMsg(m interface{}) error {
	panic("implement me")
}

func (c *chatsubscribeServer) RecvMsg(m interface{}) error {
	panic("implement me")
}

func newTestServer() *ChatService {
	return &ChatService{
		repo:        inmemoryrepo.NewInMemoryRepo(),
		connections: nil,
	}
}
