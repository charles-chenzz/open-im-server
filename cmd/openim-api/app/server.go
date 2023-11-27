package app

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/cmd/openim-api/app/options"
	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/logger"
	"github.com/spf13/cobra"
	"net"
)

func NewAPIRouterCommand() *cobra.Command {
	s := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "api",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// activate logger here
			err := logger.Apply(s.Logs)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config.Init()
			return Run()
		},
	}

	cmd.Flags().String(constant.FlagPrometheusPort, "20100", "prom port")
	cmd.Flags().StringP(constant.FlagConf, "c", "./config/", "path to config file folder")
	cmd.Flags().IntP(constant.FlagPort, "p", 10002, "server listen port")
	return cmd
}

func Run() error {
	/*	if port == 0 || proPort == 0 {
		return fmt.Errorf("port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort))
	}*/
	fmt.Println("start")
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	log.ZInfo(context.Background(), "api start init discov client")

	var client discoveryregistry.SvcDiscoveryRegistry

	// Determine whether zk is passed according to whether it is a clustered deployment
	client, err = kdisc.NewDiscoveryRegister("zookeeper")
	if err != nil {
		fmt.Println("failed to register")
		return err
	}
	if err = client.CreateRpcRootNodes(options.DefaultRpcGroups); err != nil {
		fmt.Println("failed to create root nodes")
		return err
	}

	// todo can we get rid of this?
	if err = client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, config.Config.EncodeConfig()); err != nil {
		return err
	}
	log.ZInfo(context.Background(), "api register public config to discov success")
	router := api.NewGinRouter(client, rdb)

	// safe to disable the code, we won't enable the prom until we finish refactor.
	/*	if config.Config.Prometheus.Enable {
		p := ginProm.NewPrometheus("app", prommetrics.GetGinCusMetrics("Api"))
		p.SetListenAddress(fmt.Sprintf(":%d", proPort))
		p.Use(router)
	}*/

	log.ZInfo(context.Background(), "api init router success")
	//var address string
	/*if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}*/

	// in container, we're safe to run this by default
	address := net.JoinHostPort("0.0.0.0", "10002")
	err = router.Run(address)
	if err != nil {
		return err
	}

	return nil
}
