// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/openimsdk/open-im-server/v3/cmd/openim-api/app/options"
	"net/http"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/tokenverify"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mw"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

func NewGinRouter(discov discoveryregistry.SvcDiscoveryRegistry, rdb redis.UniversalClient) *gin.Engine {
	discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("required_if", RequiredIf)
	}

	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	// init rpc client here
	userRpc := rpcclient.NewUser(discov)
	groupRpc := rpcclient.NewGroup(discov)
	friendRpc := rpcclient.NewFriend(discov)
	messageRpc := rpcclient.NewMessage(discov)
	conversationRpc := rpcclient.NewConversation(discov)
	authRpc := rpcclient.NewAuth(discov)
	thirdRpc := rpcclient.NewThird(discov)

	u := NewUserApi(*userRpc)
	m := NewMessageApi(messageRpc, userRpc)
	ParseToken := GinParseToken(rdb)
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", ParseToken, u.UpdateUserInfo)
		userRouterGroup.POST("/set_global_msg_recv_opt", ParseToken, u.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", ParseToken, u.GetUsersPublicInfo)
		userRouterGroup.POST("/get_all_users_uid", ParseToken, u.GetAllUsersID)
		userRouterGroup.POST("/account_check", ParseToken, u.AccountCheck)
		userRouterGroup.POST("/get_users", ParseToken, u.GetUsers)
		userRouterGroup.POST("/get_users_online_status", ParseToken, u.GetUsersOnlineStatus)
		userRouterGroup.POST("/get_users_online_token_detail", ParseToken, u.GetUsersOnlineTokenDetail)
		userRouterGroup.POST("/subscribe_users_status", ParseToken, u.SubscriberStatus)
		userRouterGroup.POST("/get_users_status", ParseToken, u.GetUserStatus)
		userRouterGroup.POST("/get_subscribe_users_status", ParseToken, u.GetSubscribeUsersStatus)
	}
	// friend routing group
	friendRouterGroup := r.Group("/friend", ParseToken)
	{
		f := NewFriendApi(*friendRpc)
		friendRouterGroup.POST("/delete_friend", f.DeleteFriend)
		friendRouterGroup.POST("/get_friend_apply_list", f.GetFriendApplyList)
		friendRouterGroup.POST("/get_designated_friend_apply", f.GetDesignatedFriendsApply)
		friendRouterGroup.POST("/get_self_friend_apply_list", f.GetSelfApplyList)
		friendRouterGroup.POST("/get_friend_list", f.GetFriendList)
		friendRouterGroup.POST("/get_designated_friends", f.GetDesignatedFriends)
		friendRouterGroup.POST("/add_friend", f.ApplyToAddFriend)
		friendRouterGroup.POST("/add_friend_response", f.RespondFriendApply)
		friendRouterGroup.POST("/set_friend_remark", f.SetFriendRemark)
		friendRouterGroup.POST("/add_black", f.AddBlack)
		friendRouterGroup.POST("/get_black_list", f.GetPaginationBlacks)
		friendRouterGroup.POST("/remove_black", f.RemoveBlack)
		friendRouterGroup.POST("/import_friend", f.ImportFriends)
		friendRouterGroup.POST("/is_friend", f.IsFriend)
		friendRouterGroup.POST("/get_friend_id", f.GetFriendIDs)
		friendRouterGroup.POST("/get_specified_friends_info", f.GetSpecifiedFriendsInfo)
	}
	g := NewGroupApi(*groupRpc)
	groupRouterGroup := r.Group("/group", ParseToken)
	{
		groupRouterGroup.POST("/create_group", g.CreateGroup)
		groupRouterGroup.POST("/set_group_info", g.SetGroupInfo)
		groupRouterGroup.POST("/join_group", g.JoinGroup)
		groupRouterGroup.POST("/quit_group", g.QuitGroup)
		groupRouterGroup.POST("/group_application_response", g.ApplicationGroupResponse)
		groupRouterGroup.POST("/transfer_group", g.TransferGroupOwner)
		groupRouterGroup.POST("/get_recv_group_applicationList", g.GetRecvGroupApplicationList)
		groupRouterGroup.POST("/get_user_req_group_applicationList", g.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_group_users_req_application_list", g.GetGroupUsersReqApplicationList)
		groupRouterGroup.POST("/get_groups_info", g.GetGroupsInfo)
		groupRouterGroup.POST("/kick_group", g.KickGroupMember)
		groupRouterGroup.POST("/get_group_members_info", g.GetGroupMembersInfo)
		groupRouterGroup.POST("/get_group_member_list", g.GetGroupMemberList)
		groupRouterGroup.POST("/invite_user_to_group", g.InviteUserToGroup)
		groupRouterGroup.POST("/get_joined_group_list", g.GetJoinedGroupList)
		groupRouterGroup.POST("/dismiss_group", g.DismissGroup) //
		groupRouterGroup.POST("/mute_group_member", g.MuteGroupMember)
		groupRouterGroup.POST("/cancel_mute_group_member", g.CancelMuteGroupMember)
		groupRouterGroup.POST("/mute_group", g.MuteGroup)
		groupRouterGroup.POST("/cancel_mute_group", g.CancelMuteGroup)
		groupRouterGroup.POST("/set_group_member_info", g.SetGroupMemberInfo)
		groupRouterGroup.POST("/get_group_abstract_info", g.GetGroupAbstractInfo)
		groupRouterGroup.POST("/get_groups", g.GetGroups)
		groupRouterGroup.POST("/get_group_member_user_id", g.GetGroupMemberUserIDs)
	}
	superGroupRouterGroup := r.Group("/super_group", ParseToken)
	{
		superGroupRouterGroup.POST("/get_joined_group_list", g.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", g.GetSuperGroupsInfo)
	}
	// certificate
	authRouterGroup := r.Group("/auth")
	{
		a := NewAuthApi(*authRpc)
		authRouterGroup.POST("/user_token", a.UserToken)
		authRouterGroup.POST("/parse_token", a.ParseToken)
		authRouterGroup.POST("/force_logout", ParseToken, a.ForceLogout)
	}
	// Third service
	thirdGroup := r.Group("/third", ParseToken)
	{
		thirdGroup.GET("/prometheus", GetPrometheus)
		t := NewThirdApi(*thirdRpc)
		thirdGroup.POST("/fcm_update_token", t.FcmUpdateToken)
		thirdGroup.POST("/set_app_badge", t.SetAppBadge)

		logs := thirdGroup.Group("/logs")
		logs.POST("/upload", t.UploadLogs)
		logs.POST("/delete", t.DeleteLogs)
		logs.POST("/search", t.SearchLogs)

		objectGroup := r.Group("/object", ParseToken)

		objectGroup.POST("/part_limit", t.PartLimit)
		objectGroup.POST("/part_size", t.PartSize)
		objectGroup.POST("/initiate_multipart_upload", t.InitiateMultipartUpload)
		objectGroup.POST("/auth_sign", t.AuthSign)
		objectGroup.POST("/complete_multipart_upload", t.CompleteMultipartUpload)
		objectGroup.POST("/access_url", t.AccessURL)
		objectGroup.GET("/*name", t.ObjectRedirect)
	}
	// Message
	msgGroup := r.Group("/msg", ParseToken)
	{
		msgGroup.POST("/newest_seq", m.GetSeq)
		msgGroup.POST("/search_msg", m.SearchMsg)
		msgGroup.POST("/send_msg", m.SendMessage)
		msgGroup.POST("/send_business_notification", m.SendBusinessNotification)
		msgGroup.POST("/pull_msg_by_seq", m.PullMsgBySeqs)
		msgGroup.POST("/revoke_msg", m.RevokeMsg)
		msgGroup.POST("/mark_msgs_as_read", m.MarkMsgsAsRead)
		msgGroup.POST("/mark_conversation_as_read", m.MarkConversationAsRead)
		msgGroup.POST("/get_conversations_has_read_and_max_seq", m.GetConversationsHasReadAndMaxSeq)
		msgGroup.POST("/set_conversation_has_read_seq", m.SetConversationHasReadSeq)

		msgGroup.POST("/clear_conversation_msg", m.ClearConversationsMsg)
		msgGroup.POST("/user_clear_all_msg", m.UserClearAllMsg)
		msgGroup.POST("/delete_msgs", m.DeleteMsgs)
		msgGroup.POST("/delete_msg_phsical_by_seq", m.DeleteMsgPhysicalBySeq)
		msgGroup.POST("/delete_msg_physical", m.DeleteMsgPhysical)

		msgGroup.POST("/batch_send_msg", m.BatchSendMsg)
		msgGroup.POST("/check_msg_is_send_success", m.CheckMsgIsSendSuccess)
		msgGroup.POST("/get_server_time", m.GetServerTime)
	}
	// Conversation
	conversationGroup := r.Group("/conversation", ParseToken)
	{
		c := NewConversationApi(*conversationRpc)
		conversationGroup.POST("/get_all_conversations", c.GetAllConversations)
		conversationGroup.POST("/get_conversation", c.GetConversation)
		conversationGroup.POST("/get_conversations", c.GetConversations)
		conversationGroup.POST("/set_conversations", c.SetConversations)
		conversationGroup.POST("/get_conversation_offline_push_user_ids", c.GetConversationOfflinePushUserIDs)
	}

	statisticsGroup := r.Group("/statistics", ParseToken)
	{
		statisticsGroup.POST("/user/register", u.UserRegisterCount)
		statisticsGroup.POST("/user/active", m.GetActiveUser)
		statisticsGroup.POST("/group/create", g.GroupCreateCount)
		statisticsGroup.POST("/group/active", m.GetActiveGroup)
	}
	return r
}

func GinParseToken(rdb redis.UniversalClient) gin.HandlerFunc {
	dataBase := controller.NewAuthDatabase(
		cache.NewMsgCacheModel(rdb),
		config.Config.Secret,
		config.Config.TokenPolicy.Expire,
	)
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost:
			token := c.Request.Header.Get(constant.Token)
			if token == "" {
				log.ZWarn(c, "header get token error", errs.ErrArgs.Wrap("header must have token"))
				apiresp.GinError(c, errs.ErrArgs.Wrap("header must have token"))
				c.Abort()
				return
			}
			claims, err := tokenverify.GetClaimFromToken(token, authverify.Secret())
			if err != nil {
				log.ZWarn(c, "jwt get token error", errs.ErrTokenUnknown.Wrap())
				apiresp.GinError(c, errs.ErrTokenUnknown.Wrap())
				c.Abort()
				return
			}
			m, err := dataBase.GetTokensWithoutError(c, claims.UserID, claims.PlatformID)
			if err != nil {
				log.ZWarn(c, "cache get token error", errs.ErrTokenNotExist.Wrap())
				apiresp.GinError(c, errs.ErrTokenNotExist.Wrap())
				c.Abort()
				return
			}
			if len(m) == 0 {
				log.ZWarn(c, "cache do not exist token error", errs.ErrTokenNotExist.Wrap())
				apiresp.GinError(c, errs.ErrTokenNotExist.Wrap())
				c.Abort()
				return
			}
			if v, ok := m[token]; ok {
				switch v {
				case constant.NormalToken:
				case constant.KickedToken:
					log.ZWarn(c, "cache kicked token error", errs.ErrTokenKicked.Wrap())
					apiresp.GinError(c, errs.ErrTokenKicked.Wrap())
					c.Abort()
					return
				default:
					log.ZWarn(c, "cache unknown token error", errs.ErrTokenUnknown.Wrap())
					apiresp.GinError(c, errs.ErrTokenUnknown.Wrap())
					c.Abort()
					return
				}
			} else {
				apiresp.GinError(c, errs.ErrTokenNotExist.Wrap())
				c.Abort()
				return
			}
			c.Set(constant.OpUserPlatform, constant.PlatformIDToName(claims.PlatformID))
			c.Set(constant.OpUserID, claims.UserID)
			c.Next()
		}
	}
}

func NewRouter(discov discoveryregistry.SvcDiscoveryRegistry, rdb redis.UniversalClient) *gin.Engine {
	discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("required_if", RequiredIf)
	}

	r.Use(options.WithCors(), options.WithOperationId(), options.WithRecovery())

	rpc := NewRPC(discov)
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", rpc.UserRegister)
		userRouterGroup.POST("/update_user_info", options.WithToken(rdb), rpc.UpdateUserInfo)
		userRouterGroup.POST("/set_global_msg_recv_opt", options.WithToken(rdb), rpc.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", options.WithToken(rdb), rpc.GetUsersPublicInfo)
		userRouterGroup.POST("/get_all_users_uid", options.WithToken(rdb), rpc.GetAllUsersID)
		userRouterGroup.POST("/account_check", options.WithToken(rdb), rpc.AccountCheck)
		userRouterGroup.POST("/get_users", options.WithToken(rdb), rpc.GetUsers)
		userRouterGroup.POST("/get_users_online_status", options.WithToken(rdb), rpc.GetUsersOnlineStatus)
		userRouterGroup.POST("/get_users_online_token_detail", options.WithToken(rdb), rpc.GetUsersOnlineTokenDetail)
		userRouterGroup.POST("/subscribe_users_status", options.WithToken(rdb), rpc.SubscriberStatus)
		userRouterGroup.POST("/get_users_status", options.WithToken(rdb), rpc.GetUserStatus)
		userRouterGroup.POST("/get_subscribe_users_status", options.WithToken(rdb), rpc.GetSubscribeUsersStatus)
	}
	// friend routing group
	friendRouterGroup := r.Group("/friend", options.WithToken(rdb))
	{
		friendRouterGroup.POST("/delete_friend", rpc.DeleteFriend)
		friendRouterGroup.POST("/get_friend_apply_list", rpc.GetFriendApplyList)
		friendRouterGroup.POST("/get_designated_friend_apply", rpc.GetDesignatedFriendsApply)
		friendRouterGroup.POST("/get_self_friend_apply_list", rpc.GetSelfApplyList)
		friendRouterGroup.POST("/get_friend_list", rpc.GetFriendList)
		friendRouterGroup.POST("/get_designated_friends", rpc.GetDesignatedFriends)
		friendRouterGroup.POST("/add_friend", rpc.ApplyToAddFriend)
		friendRouterGroup.POST("/add_friend_response", rpc.RespondFriendApply)
		friendRouterGroup.POST("/set_friend_remark", rpc.SetFriendRemark)
		friendRouterGroup.POST("/add_black", rpc.AddBlack)
		friendRouterGroup.POST("/get_black_list", rpc.GetPaginationBlacks)
		friendRouterGroup.POST("/remove_black", rpc.RemoveBlack)
		friendRouterGroup.POST("/import_friend", rpc.ImportFriends)
		friendRouterGroup.POST("/is_friend", rpc.IsFriend)
		friendRouterGroup.POST("/get_friend_id", rpc.GetFriendIDs)
		friendRouterGroup.POST("/get_specified_friends_info", rpc.GetSpecifiedFriendsInfo)
	}

	groupRouterGroup := r.Group("/group", options.WithToken(rdb))
	{
		groupRouterGroup.POST("/create_group", rpc.CreateGroup)
		groupRouterGroup.POST("/set_group_info", rpc.SetGroupInfo)
		groupRouterGroup.POST("/join_group", rpc.JoinGroup)
		groupRouterGroup.POST("/quit_group", rpc.QuitGroup)
		groupRouterGroup.POST("/group_application_response", rpc.ApplicationGroupResponse)
		groupRouterGroup.POST("/transfer_group", rpc.TransferGroupOwner)
		groupRouterGroup.POST("/get_recv_group_applicationList", rpc.GetRecvGroupApplicationList)
		groupRouterGroup.POST("/get_user_req_group_applicationList", rpc.GetUserReqGroupApplicationList)
		groupRouterGroup.POST("/get_group_users_req_application_list", rpc.GetGroupUsersReqApplicationList)
		groupRouterGroup.POST("/get_groups_info", rpc.GetGroupsInfo)
		groupRouterGroup.POST("/kick_group", rpc.KickGroupMember)
		groupRouterGroup.POST("/get_group_members_info", rpc.GetGroupMembersInfo)
		groupRouterGroup.POST("/get_group_member_list", rpc.GetGroupMemberList)
		groupRouterGroup.POST("/invite_user_to_group", rpc.InviteUserToGroup)
		groupRouterGroup.POST("/get_joined_group_list", rpc.GetJoinedGroupList)
		groupRouterGroup.POST("/dismiss_group", rpc.DismissGroup)
		groupRouterGroup.POST("/mute_group_member", rpc.MuteGroupMember)
		groupRouterGroup.POST("/cancel_mute_group_member", rpc.CancelMuteGroupMember)
		groupRouterGroup.POST("/mute_group", rpc.MuteGroup)
		groupRouterGroup.POST("/cancel_mute_group", rpc.CancelMuteGroup)
		groupRouterGroup.POST("/set_group_member_info", rpc.SetGroupMemberInfo)
		groupRouterGroup.POST("/get_group_abstract_info", rpc.GetGroupAbstractInfo)
		groupRouterGroup.POST("/get_groups", rpc.GetGroups)
		groupRouterGroup.POST("/get_group_member_user_id", rpc.GetGroupMemberUserIDs)
	}
	superGroupRouterGroup := r.Group("/super_group", options.WithToken(rdb))
	{
		superGroupRouterGroup.POST("/get_joined_group_list", rpc.GetJoinedSuperGroupList)
		superGroupRouterGroup.POST("/get_groups_info", rpc.GetSuperGroupsInfo)
	}
	// certificate
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/user_token", rpc.UserToken)
		authRouterGroup.POST("/parse_token", rpc.ParseToken)
		authRouterGroup.POST("/force_logout", options.WithToken(rdb), rpc.ForceLogout)
	}
	// Third service
	thirdGroup := r.Group("/third", options.WithToken(rdb))
	{
		thirdGroup.GET("/prometheus", GetPrometheus)

		thirdGroup.POST("/fcm_update_token", rpc.FcmUpdateToken)
		thirdGroup.POST("/set_app_badge", rpc.SetAppBadge)

		logs := thirdGroup.Group("/logs")
		logs.POST("/upload", rpc.UploadLogs)
		logs.POST("/delete", rpc.DeleteLogs)
		logs.POST("/search", rpc.SearchLogs)

		objectGroup := r.Group("/object", options.WithToken(rdb))

		objectGroup.POST("/part_limit", rpc.PartLimit)
		objectGroup.POST("/part_size", rpc.PartSize)
		objectGroup.POST("/initiate_multipart_upload", rpc.InitiateMultipartUpload)
		objectGroup.POST("/auth_sign", rpc.AuthSign)
		objectGroup.POST("/complete_multipart_upload", rpc.CompleteMultipartUpload)
		objectGroup.POST("/access_url", rpc.AccessURL)
		objectGroup.GET("/*name", rpc.ObjectRedirect)
	}
	// Message
	msgGroup := r.Group("/msg", options.WithToken(rdb))
	{
		msgGroup.POST("/newest_seq", rpc.GetSeq)
		msgGroup.POST("/search_msg", rpc.SearchMsg)
		msgGroup.POST("/send_msg", rpc.SendMessage)
		msgGroup.POST("/send_business_notification", rpc.SendBusinessNotification)
		msgGroup.POST("/pull_msg_by_seq", rpc.PullMsgBySeqs)
		msgGroup.POST("/revoke_msg", rpc.RevokeMsg)
		msgGroup.POST("/mark_msgs_as_read", rpc.MarkMsgsAsRead)
		msgGroup.POST("/mark_conversation_as_read", rpc.MarkConversationAsRead)
		msgGroup.POST("/get_conversations_has_read_and_max_seq", rpc.GetConversationsHasReadAndMaxSeq)
		msgGroup.POST("/set_conversation_has_read_seq", rpc.SetConversationHasReadSeq)

		msgGroup.POST("/clear_conversation_msg", rpc.ClearConversationsMsg)
		msgGroup.POST("/user_clear_all_msg", rpc.UserClearAllMsg)
		msgGroup.POST("/delete_msgs", rpc.DeleteMsgs)
		msgGroup.POST("/delete_msg_phsical_by_seq", rpc.DeleteMsgPhysicalBySeq)
		msgGroup.POST("/delete_msg_physical", rpc.DeleteMsgPhysical)

		msgGroup.POST("/batch_send_msg", rpc.BatchSendMsg)
		msgGroup.POST("/check_msg_is_send_success", rpc.CheckMsgIsSendSuccess)
		msgGroup.POST("/get_server_time", rpc.GetServerTime)
	}
	// Conversation
	conversationGroup := r.Group("/conversation", options.WithToken(rdb))
	{
		conversationGroup.POST("/get_all_conversations", rpc.GetAllConversations)
		conversationGroup.POST("/get_conversation", rpc.GetConversation)
		conversationGroup.POST("/get_conversations", rpc.GetConversations)
		conversationGroup.POST("/set_conversations", rpc.SetConversations)
		conversationGroup.POST("/get_conversation_offline_push_user_ids", rpc.GetConversationOfflinePushUserIDs)
	}

	statisticsGroup := r.Group("/statistics", options.WithToken(rdb))
	{
		statisticsGroup.POST("/user/register", rpc.UserRegisterCount)
		statisticsGroup.POST("/user/active", rpc.GetActiveUser)
		statisticsGroup.POST("/group/create", rpc.GroupCreateCount)
		statisticsGroup.POST("/group/active", rpc.GetActiveGroup)
	}
	return r
}
