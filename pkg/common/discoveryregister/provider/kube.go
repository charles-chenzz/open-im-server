package provider

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"google.golang.org/grpc"
)

type K8sDR struct {
	options         []grpc.DialOption
	rpcRegisterAddr string
}

func NewK8sDiscoveryRegister() (discoveryregistry.SvcDiscoveryRegistry, error) {
	return &K8sDR{}, nil
}

func (cli *K8sDR) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	cli.rpcRegisterAddr = serviceName
	return nil
}
func (cli *K8sDR) UnRegister() error {

	return nil
}
func (cli *K8sDR) CreateRpcRootNodes(serviceNames []string) error {

	return nil
}
func (cli *K8sDR) RegisterConf2Registry(key string, conf []byte) error {

	return nil
}

func (cli *K8sDR) GetConfFromRegistry(key string) ([]byte, error) {

	return nil, nil
}
func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {

	conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
	return []*grpc.ClientConn{conn}, err
}
func (cli *K8sDR) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	return grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
}
func (cli *K8sDR) GetSelfConnTarget() string {

	return cli.rpcRegisterAddr
}
func (cli *K8sDR) AddOption(opts ...grpc.DialOption) {
	cli.options = append(cli.options, opts...)
}
func (cli *K8sDR) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}

// do not use this method for call rpc
func (cli *K8sDR) GetClientLocalConns() map[string][]*grpc.ClientConn {
	fmt.Println("should not call this function!!!!!!!!!!!!!!!!!!!!!!!!!")
	return nil
}
func (cli *K8sDR) Close() {
	return
}
