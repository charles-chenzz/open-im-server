package api

import (
	"github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type Conversation rpcclient.Conversation

func NewConversation(discov discoveryregistry.SvcDiscoveryRegistry) Conversation {
	return Conversation(*rpcclient.NewConversation(discov))
}

func (o *Conversation) GetAllConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetAllConversations, o.Client, c)
}

func (o *Conversation) GetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversation, o.Client, c)
}

func (o *Conversation) GetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversations, o.Client, c)
}

func (o *Conversation) SetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversations, o.Client, c)
}

func (o *Conversation) GetConversationOfflinePushUserIDs(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversationOfflinePushUserIDs, o.Client, c)
}
