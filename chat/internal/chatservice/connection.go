package chatservice

import "github.com/golang-friends/slack-clone/chat/protos/chatpb"

// Connection is one user connection.
type Connection struct {
	server *chatpb.ChatService_SubscribeServer
	meta ConnectionMeta
}

// ConnectionMeta contains meta data of the Connection.
type ConnectionMeta struct {
	workspaceID string
}
