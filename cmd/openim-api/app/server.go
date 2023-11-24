package app

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/internal/api"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	ginProm "github.com/openimsdk/open-im-server/v3/pkg/common/ginprometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/spf13/cobra"
	"net"
	"strconv"
)

func NewAPIRouterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "api",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfgFolderPath, _ := cmd.Flags().GetString(constant.FlagConf)
			err := config.InitConfig(cfgFolderPath)
			if err != nil {
				return err
			}
			// apply default option(logger)
			err = applyLogger()
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			port, err := cmd.Flags().GetInt(constant.FlagPort)
			if err != nil {
				fmt.Println("Error getting port flag:", err)
			}
			promPort, _ := cmd.Flags().GetInt(constant.FlagPrometheusPort)
			if promPort == 0 {
				promPort = config.Config.Prometheus.ApiPrometheusPort[0]
			}
			return run(port, promPort)
		},
	}

	cmd.Flags().StringP(constant.FlagConf, "c", "", "path to config file folder")
	cmd.Flags().IntP(constant.FlagPort, "p", 10002, "server listen port")
	return cmd
}

func run(port int, proPort int) error {
	log.ZInfo(context.Background(), "Openim api port:", "port", port, "proPort", proPort)

	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		log.ZError(context.Background(), err, nil)

		return fmt.Errorf(err)
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize Redis", err)

		return err
	}
	log.ZInfo(context.Background(), "api start init discov client")

	var client discoveryregistry.SvcDiscoveryRegistry

	// Determine whether zk is passed according to whether it is a clustered deployment
	client, err = kdisc.NewDiscoveryRegister("")
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize discovery register", err)

		return err
	}
	if err = client.CreateRpcRootNodes(config.Config.GetServiceNames()); err != nil {
		log.ZError(context.Background(), "Failed to create RPC root nodes", err)

		return err
	}
	log.ZInfo(context.Background(), "api register public config to discov")
	if err = client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, config.Config.EncodeConfig()); err != nil {
		log.ZError(context.Background(), "Failed to register public config to discov", err)

		return err
	}
	log.ZInfo(context.Background(), "api register public config to discov success")
	router := api.NewGinRouter(client, rdb)

	if config.Config.Prometheus.Enable {
		p := ginProm.NewPrometheus("app", prommetrics.GetGinCusMetrics("Api"))
		p.SetListenAddress(fmt.Sprintf(":%d", proPort))
		p.Use(router)
	}

	log.ZInfo(context.Background(), "api init router success")
	var address string
	if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	log.ZInfo(context.Background(), "start api server", "address", address, "OpenIM version", config.Version)

	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "api run failed", err, "address", address)

		return err
	}

	return nil
}

// this is not done yet
func applyLogger() error {
	logConfig := config.Config.Log

	return log.InitFromConfig(
		"",
		"",
		logConfig.RemainLogLevel,
		logConfig.IsStdout,
		logConfig.IsJson,
		logConfig.StorageLocation,
		logConfig.RemainRotationCount,
		logConfig.RotationTime,
	)
}
