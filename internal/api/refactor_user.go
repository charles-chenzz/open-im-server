package api

import (
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msggateway"
	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type User rpcclient.User

func NewUser(discov discoveryregistry.SvcDiscoveryRegistry) User {
	return User(*rpcclient.NewUser(discov))
}

func (u *User) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, u.Client, c)
}

func (u *User) UpdateUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.UpdateUserInfo, u.Client, c)
}

func (u *User) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.Call(user.UserClient.SetGlobalRecvMessageOpt, u.Client, c)
}

func (u *User) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.Client, c)
}

func (u *User) GetAllUsersID(c *gin.Context) {
	a2r.Call(user.UserClient.GetAllUserID, u.Client, c)
}

func (u *User) AccountCheck(c *gin.Context) {
	a2r.Call(user.UserClient.AccountCheck, u.Client, c)
}

func (u *User) GetUsers(c *gin.Context) {
	a2r.Call(user.UserClient.GetPaginationUsers, u.Client, c)
}

// GetUsersOnlineStatus Get user online status.
func (u *User) GetUsersOnlineStatus(c *gin.Context) {
	var req msggateway.GetUsersOnlineStatusReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conns, err := u.Discov.GetConns(c, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var wsResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	flag := false

	// Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.GetUsersOnlineStatus(c, &req)
		if err != nil {
			log.ZWarn(c, "GetUsersOnlineStatus rpc err", err)

			parseError := apiresp.ParseError(err)
			if parseError.ErrCode == errs.NoPermissionError {
				apiresp.GinError(c, err)
				return
			}
		} else {
			wsResult = append(wsResult, reply.SuccessResult...)
		}
	}
	// Traversing the userIDs in the api request body
	for _, v1 := range req.UserIDs {
		flag = false
		res := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		// Iterate through the online results fetched from various gateways
		for _, v2 := range wsResult {
			// If matches the above description on the line, and vice versa
			if v2.UserID == v1 {
				flag = true
				res.UserID = v1
				res.Status = constant.OnlineStatus
				res.DetailPlatformStatus = append(res.DetailPlatformStatus, v2.DetailPlatformStatus...)
				break
			}
		}
		if !flag {
			res.UserID = v1
			res.Status = constant.OfflineStatus
		}
		respResult = append(respResult, res)
	}
	apiresp.GinSuccess(c, respResult)
}

func (u *User) UserRegisterCount(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegisterCount, u.Client, c)
}

// GetUsersOnlineTokenDetail Get user online token details.
func (u *User) GetUsersOnlineTokenDetail(c *gin.Context) {
	var wsResult []*msggateway.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*msggateway.SingleDetail
	flag := false
	var req msggateway.GetUsersOnlineStatusReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conns, err := u.Discov.GetConns(c, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	// Online push message
	for _, v := range conns {
		msgClient := msggateway.NewMsgGatewayClient(v)
		reply, err := msgClient.GetUsersOnlineStatus(c, &req)
		if err != nil {
			log.ZWarn(c, "GetUsersOnlineStatus rpc  err", err)
			continue
		} else {
			wsResult = append(wsResult, reply.SuccessResult...)
		}
	}

	for _, v1 := range req.UserIDs {
		m := make(map[string][]string, 10)
		flag = false
		temp := new(msggateway.SingleDetail)
		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				temp.UserID = v1
				temp.Status = constant.OnlineStatus
				for _, status := range v2.DetailPlatformStatus {
					if v, ok := m[status.Platform]; ok {
						m[status.Platform] = append(v, status.Token)
					} else {
						m[status.Platform] = []string{status.Token}
					}
				}
			}
		}
		for p, tokens := range m {
			t := new(msggateway.SinglePlatformToken)
			t.Platform = p
			t.Token = tokens
			t.Total = int32(len(tokens))
			temp.SinglePlatformToken = append(temp.SinglePlatformToken, t)
		}

		if flag {
			respResult = append(respResult, temp)
		}
	}

	apiresp.GinSuccess(c, respResult)
}

// SubscriberStatus Presence status of subscribed users.
func (u *User) SubscriberStatus(c *gin.Context) {
	a2r.Call(user.UserClient.SubscribeOrCancelUsersStatus, u.Client, c)
}

// GetUserStatus Get the online status of the user.
func (u *User) GetUserStatus(c *gin.Context) {
	a2r.Call(user.UserClient.GetUserStatus, u.Client, c)
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func (u *User) GetSubscribeUsersStatus(c *gin.Context) {
	a2r.Call(user.UserClient.GetSubscribeUsersStatus, u.Client, c)
}
