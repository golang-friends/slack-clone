package inmemoryrepo

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang-friends/slack-clone/chat/protos/chatpb"
	"github.com/hashicorp/go-uuid"
	"sync"
)

// InMemoryRepo is the implementation of MessageRepository.
type InMemoryRepo struct {
	messages []*chatpb.Message
	sync.RWMutex
}

// GetAllMessages returns all messages in the memory.
func (i *InMemoryRepo) GetAllMessages() ([]*chatpb.Message, error) {
	i.RLock()
	defer i.RUnlock()
	return i.messages, nil
}

// AddNewMessage will create a new message when received.
func (i *InMemoryRepo) AddNewMessage(message *chatpb.Message) (*chatpb.Message, error) {
	randomID, _ := uuid.GenerateUUID()
	message.MessageId = randomID

	i.Lock()
	defer i.Unlock()

	i.messages = append(i.messages, message)
	return message, nil
}

// UpdateMessage will update the message based on the message get id.
func (i *InMemoryRepo) UpdateMessage(message *chatpb.Message) (*chatpb.Message, error) {
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

// NewInMemoryRepo is the factory for InMemoryRepo.
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{}
}
