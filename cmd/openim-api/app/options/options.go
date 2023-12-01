package options

import (
	"errors"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/apiresp"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/tokenverify"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/logger"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prom"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var DefaultRpcGroups = []string{"User", "Friend", "Msg", "Push", "MessageGateway", "Group", "Auth", "Conversation", "Third"}

type ServerRunOptions struct {
	*logger.Options
	monitorOptions *prom.Options
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Options: &logger.Options{
			StorageLocation:     "../logs/",
			RotationTime:        24,
			RemainRotationCount: 2,
			RemainLogLevel:      6,
			IsStdout:            false,
			IsJson:              true,
			WithStack:           false,
		},
		monitorOptions: &prom.Options{
			Enable:                        false,
			PrometheusUrl:                 "",
			ApiPrometheusPort:             []int{},
			UserPrometheusPort:            []int{},
			FriendPrometheusPort:          []int{},
			MessagePrometheusPort:         []int{},
			MessageGatewayPrometheusPort:  []int{},
			GroupPrometheusPort:           []int{},
			AuthPrometheusPort:            []int{},
			PushPrometheusPort:            []int{},
			ConversationPrometheusPort:    []int{},
			RtcPrometheusPort:             []int{},
			MessageTransferPrometheusPort: []int{},
			ThirdPrometheusPort:           []int{},
		},
	}
}

// WithCors setting up the cors
func WithCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar",
		) // Cross-domain key settings allow browsers to resolve.
		c.Header(
			"Access-Control-Max-Age",
			"172800",
		) // Cache request information in seconds.
		c.Header(
			"Access-Control-Allow-Credentials",
			"false",
		) //  Whether cross-domain requests need to carry cookie information, the default setting is true.
		c.Header(
			"content-type",
			"application/json",
		) // Set the return format to json.
		// Release all option pre-requests
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
			return
		}
		c.Next()
	}
}

// WithOperationId setting up operationId for trace
func WithOperationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			operationID := c.Request.Header.Get(constant.OperationID)
			if operationID == "" {
				err := errors.New("header must have operationID")
				apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error()))
				c.Abort()
				return
			}
			c.Set(constant.OperationID, operationID)
		}
		c.Next()
	}
}

// WithRecovery simply wrap gin.Recovery
func WithRecovery() gin.HandlerFunc {
	return gin.Recovery()
}

// WithToken setting up the token verify for services like user
func WithToken(rdb redis.UniversalClient) gin.HandlerFunc {
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
