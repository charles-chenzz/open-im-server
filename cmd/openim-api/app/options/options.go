package options

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/logger"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prom"
)

var DefaultRpcGroups = []string{"User", "Friend", "Msg", "Push", "MessageGateway", "Group", "Auth", "Conversation", "Third"}

type ServerRunOptions struct {
	*logger.Logs
	Monitor *prom.Monitor
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		Logs: &logger.Logs{
			StorageLocation:     "../logs/",
			RotationTime:        24,
			RemainRotationCount: 2,
			RemainLogLevel:      6,
			IsStdout:            false,
			IsJson:              true,
			WithStack:           false,
		},
		Monitor: &prom.Monitor{
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
