package mongodbrepo

import (
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
)

// MongoDBRepo is the implementation of MessageRepository.
type MongoDBRepo struct {
}

// GetAllMessages not implemented.
func (m *MongoDBRepo) GetAllMessages() ([]*chatpb.Message, error) {
	panic("implement me")
}

// AddNewMessage not implemented.
func (m *MongoDBRepo) AddNewMessage(message *chatpb.Message) (*chatpb.Message, error) {
	panic("implement me")
}

// UpdateMessage not implemented.
func (m *MongoDBRepo) UpdateMessage(message *chatpb.Message) (*chatpb.Message, error) {
	panic("implement me")
}

// NewMongoDBRepo is the factory function for MongoDBRepo.
func NewMongoDBRepo() *MongoDBRepo {
	return &MongoDBRepo{}
}
