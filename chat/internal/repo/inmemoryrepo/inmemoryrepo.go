package inmemoryrepo

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
	"github.com/hashicorp/go-uuid"
	"sync"
)

type inmemoryRepo struct {
	messages []*chatpb.Message
	sync.RWMutex
}

func (i *inmemoryRepo) GetAllMessage() ([]*chatpb.Message, error) {
	i.RLock()
	defer i.RUnlock()
	return i.messages, nil
}

func (i *inmemoryRepo) AddNewMessage(message *chatpb.Message) (*chatpb.Message, error) {
	randomID, _ := uuid.GenerateUUID()
	message.MessageId = randomID

	i.Lock()
	defer i.Unlock()

	i.messages = append(i.messages, message)
	return message, nil
}

func (i *inmemoryRepo) UpdateMessage(message *chatpb.Message) (*chatpb.Message, error) {
	i.Lock()
	defer i.Unlock()
	for idx, msg := range i.messages {
		if msg.GetMessageId() == message.GetMessageId() {
			proto.Merge(msg, message)
			i.messages[idx] = msg
			return msg, nil
		}
	}
	return nil, fmt.Errorf("unable to ffnd the message: %+v", message)
}

func NewInmemoryRepo() *inmemoryRepo {
	return &inmemoryRepo{}
}
