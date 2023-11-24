package discoveryregister

import (
	"errors"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister/provider"
)

func NewDiscoveryRegister(envType string) (discoveryregistry.SvcDiscoveryRegistry, error) {
	var client discoveryregistry.SvcDiscoveryRegistry
	var err error
	switch envType {
	case "zookeeper":
		client, err = provider.NewZookeeperDiscovery()
	case "k8s":
		client, err = provider.NewK8sDiscoveryRegister()
	default:
		client = nil
		err = errors.New("envType not correct")
	}
	return client, err
}
