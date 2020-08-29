package application

import (
	"fmt"
	"github.com/golang-friends/slack-clone/chat/internal/chatservice"
	"github.com/golang-friends/slack-clone/chat/internal/config"
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// Application is the Application struct to deal with gRPC server.
type Application struct {
	cfg  *config.Config
	repo chatservice.MessageRepository
}

// Start starts the server and listens to the port.
func (a *Application) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", a.cfg.Port))
	if err != nil {
		logrus.WithError(err).Fatal("unable to listen the server")
	}

	s := grpc.NewServer()
	chatpb.RegisterChatServiceServer(s, chatservice.NewChatService(a.repo))

	reflection.Register(s)

	logrus.Infof("serving at 0.0.0.0:%v", a.cfg.Port)
	return s.Serve(lis)
}

// NewApplication is the main factory function for the Application.
func NewApplication(config *config.Config, repo chatservice.MessageRepository) *Application {
	return &Application{cfg: config, repo: repo}
}
