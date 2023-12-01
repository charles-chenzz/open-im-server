package api

import (
	"github.com/OpenIMSDK/protocol/friend"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type Friend rpcclient.Friend

func NewFriend(discov discoveryregistry.SvcDiscoveryRegistry) Friend {
	return Friend(*rpcclient.NewFriend(discov))
}

func (o *Friend) ApplyToAddFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.ApplyToAddFriend, o.Client, c)
}

func (o *Friend) RespondFriendApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.RespondFriendApply, o.Client, c)
}

func (o *Friend) DeleteFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.DeleteFriend, o.Client, c)
}

func (o *Friend) GetFriendApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyTo, o.Client, c)
}

func (o *Friend) GetDesignatedFriendsApply(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetDesignatedFriendsApply, o.Client, c)
}

func (o *Friend) GetSelfApplyList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriendsApplyFrom, o.Client, c)
}

func (o *Friend) GetFriendList(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationFriends, o.Client, c)
}

func (o *Friend) GetDesignatedFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetDesignatedFriends, o.Client, c)
}

func (o *Friend) SetFriendRemark(c *gin.Context) {
	a2r.Call(friend.FriendClient.SetFriendRemark, o.Client, c)
}

func (o *Friend) AddBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.AddBlack, o.Client, c)
}

func (o *Friend) GetPaginationBlacks(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetPaginationBlacks, o.Client, c)
}

func (o *Friend) RemoveBlack(c *gin.Context) {
	a2r.Call(friend.FriendClient.RemoveBlack, o.Client, c)
}

func (o *Friend) ImportFriends(c *gin.Context) {
	a2r.Call(friend.FriendClient.ImportFriends, o.Client, c)
}

func (o *Friend) IsFriend(c *gin.Context) {
	a2r.Call(friend.FriendClient.IsFriend, o.Client, c)
}

func (o *Friend) GetFriendIDs(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetFriendIDs, o.Client, c)
}

func (o *Friend) GetSpecifiedFriendsInfo(c *gin.Context) {
	a2r.Call(friend.FriendClient.GetSpecifiedFriendsInfo, o.Client, c)
}
