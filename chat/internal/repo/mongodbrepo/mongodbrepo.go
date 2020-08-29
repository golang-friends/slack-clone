package mongodbrepo

import (
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
)

type mongoDbRepo struct {

}

func (m *mongoDbRepo) GetAllMessage() ([]*chatpb.Message, error) {
	panic("implement me")
}

func (m *mongoDbRepo) AddNewMessage(message *chatpb.Message) (*chatpb.Message, error) {
	panic("implement me")
}

func (m *mongoDbRepo) UpdateMessage(message *chatpb.Message) (*chatpb.Message, error) {
	panic("implement me")
}

func NewMongoDbRepo() *mongoDbRepo {
	return &mongoDbRepo{}
}
