package api

import (
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/go-playground/validator/v10"
)

// RPC embed all the rpc client
type RPC struct {
	User
	Friend
	Conversation
	Auth
	Group
	Message
	Third
}

func NewRPC(discov discoveryregistry.SvcDiscoveryRegistry) RPC {
	return RPC{
		User:         NewUser(discov),
		Friend:       NewFriend(discov),
		Conversation: NewConversation(discov),
		Auth:         NewAuth(discov),
		Group:        NewGroup(discov),
		Message:      NewMessage(discov),
		Third:        NewThird(discov),
	}
}

func RequiredIf(fl validator.FieldLevel) bool {
	sessionType := fl.Parent().FieldByName("SessionType").Int()
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		if fl.FieldName() == "RecvID" {
			return fl.Field().String() != ""
		}
	case constant.GroupChatType, constant.SuperGroupChatType:
		if fl.FieldName() == "GroupID" {
			return fl.Field().String() != ""
		}
	default:
		return true
	}
	return true
}
