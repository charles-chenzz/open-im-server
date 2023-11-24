package provider

import (
	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"time"
)

// todo move openkeeper to this file
func NewZookeeperDiscovery() (discoveryregistry.SvcDiscoveryRegistry, error) {
	client, err := openkeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openkeeper.WithFreq(time.Hour), openkeeper.WithUserNameAndPassword(
			config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password,
		), openkeeper.WithRoundRobin(), openkeeper.WithTimeout(10), openkeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return nil, err
	}
	return client, nil
}
