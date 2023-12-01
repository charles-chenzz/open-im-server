package api

import (
	"github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type Auth rpcclient.Auth

func NewAuth(discov discoveryregistry.SvcDiscoveryRegistry) Auth {
	return Auth(*rpcclient.NewAuth(discov))
}

func (o *Auth) UserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserToken, o.Client, c)
}

func (o *Auth) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, o.Client, c)
}

func (o *Auth) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, o.Client, c)
}
