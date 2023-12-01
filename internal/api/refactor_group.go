package api

import (
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type Group rpcclient.Group

func NewGroup(discov discoveryregistry.SvcDiscoveryRegistry) Group {
	return Group(*rpcclient.NewGroup(discov))
}

func (o *Group) CreateGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CreateGroup, o.Client, c)
}

func (o *Group) SetGroupInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupInfo, o.Client, c)
}

func (o *Group) JoinGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.JoinGroup, o.Client, c)
}

func (o *Group) QuitGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.QuitGroup, o.Client, c)
}

func (o *Group) ApplicationGroupResponse(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupApplicationResponse, o.Client, c)
}

func (o *Group) TransferGroupOwner(c *gin.Context) {
	a2r.Call(group.GroupClient.TransferGroupOwner, o.Client, c)
}

func (o *Group) GetRecvGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupApplicationList, o.Client, c)
}

func (o *Group) GetUserReqGroupApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetUserReqApplicationList, o.Client, c)
}

func (o *Group) GetGroupUsersReqApplicationList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupUsersReqApplicationList, o.Client, c)
}

func (o *Group) GetGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupsInfo, o.Client, c)
}

func (o *Group) KickGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.KickGroupMember, o.Client, c)
}

func (o *Group) GetGroupMembersInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMembersInfo, o.Client, c)
}

func (o *Group) GetGroupMemberList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMemberList, o.Client, c)
}

func (o *Group) InviteUserToGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.InviteUserToGroup, o.Client, c)
}

func (o *Group) GetJoinedGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedGroupList, o.Client, c)
}

func (o *Group) DismissGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.DismissGroup, o.Client, c)
}

func (o *Group) MuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroupMember, o.Client, c)
}

func (o *Group) CancelMuteGroupMember(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroupMember, o.Client, c)
}

func (o *Group) MuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.MuteGroup, o.Client, c)
}

func (o *Group) CancelMuteGroup(c *gin.Context) {
	a2r.Call(group.GroupClient.CancelMuteGroup, o.Client, c)
}

func (o *Group) SetGroupMemberInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.SetGroupMemberInfo, o.Client, c)
}

func (o *Group) GetGroupAbstractInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupAbstractInfo, o.Client, c)
}

func (o *Group) GetJoinedSuperGroupList(c *gin.Context) {
	a2r.Call(group.GroupClient.GetJoinedSuperGroupList, o.Client, c)
}

func (o *Group) GetSuperGroupsInfo(c *gin.Context) {
	a2r.Call(group.GroupClient.GetSuperGroupsInfo, o.Client, c)
}

func (o *Group) GroupCreateCount(c *gin.Context) {
	a2r.Call(group.GroupClient.GroupCreateCount, o.Client, c)
}

func (o *Group) GetGroups(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroups, o.Client, c)
}

func (o *Group) GetGroupMemberUserIDs(c *gin.Context) {
	a2r.Call(group.GroupClient.GetGroupMemberUserIDs, o.Client, c)
}
