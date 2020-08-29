package chatservice

import "github.com/golang-friends/slack-clone/chat/protos/chatpb"

type Connection struct {
	server *chatpb.ChatService_SubscribeServer
	meta ConnectionMeta
}

type ConnectionMeta struct {
	workspaceID string
}
