package prom

import (
	"fmt"
	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/openimsdk/open-im-server/v3/cmd/openim-api/app/options"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"time"
)

var metricsPath = "/metrics"

/*
RequestCounterURLLabelMappingFn is a function which can be supplied to the middleware to control
the cardinality of the request counter's "url" label, which might be required in some contexts.
For instance, if for a "/customer/:name" route you don't want to generate a time series for every
possible customer name, you could use this function:

	func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "name" {
				url = strings.Replace(url, p.Value, ":name", 1)
				break
			}
		}
		return url
	}

which would map "/customer/alice" and "/customer/bob" to their template "/customer/:name".
*/
type RequestCounterURLLabelMappingFn func(c *gin.Context) string

type Options struct {
	Enable                        bool   `yaml:"enable"`
	PrometheusUrl                 string `yaml:"prometheusUrl"`
	ApiPrometheusPort             []int  `yaml:"apiPrometheusPort"`
	UserPrometheusPort            []int  `yaml:"userPrometheusPort"`
	FriendPrometheusPort          []int  `yaml:"friendPrometheusPort"`
	MessagePrometheusPort         []int  `yaml:"messagePrometheusPort"`
	MessageGatewayPrometheusPort  []int  `yaml:"messageGatewayPrometheusPort"`
	GroupPrometheusPort           []int  `yaml:"groupPrometheusPort"`
	AuthPrometheusPort            []int  `yaml:"authPrometheusPort"`
	PushPrometheusPort            []int  `yaml:"pushPrometheusPort"`
	ConversationPrometheusPort    []int  `yaml:"conversationPrometheusPort"`
	RtcPrometheusPort             []int  `yaml:"rtcPrometheusPort"`
	MessageTransferPrometheusPort []int  `yaml:"messageTransferPrometheusPort"`
	ThirdPrometheusPort           []int  `yaml:"thirdPrometheusPort"`
}

// Metrics is a wrapper of prometheus collector and params that we customize
type Metrics struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

// Prom aggregate all metrics together
type Prom struct {
	reqCnt        *prometheus.CounterVec
	reqDur        *prometheus.HistogramVec
	reqSz, resSz  prometheus.Summary
	router        *gin.Engine
	listenAddress string
	GatewayConfig PushGateway

	// List is metrics list
	List []*Metrics
	// Path is metrics path
	Path string

	mappingFn func(c *gin.Context) string

	// gin.Context string to use as a prometheus URL label
	URLLabelFromContext string
}

// PushGateway contains config for pushing to Prom push gateway
type PushGateway struct {
	// Push interval in seconds
	Interval time.Duration

	// Push Gateway URL in format http://domain:port
	// where JOB NAME can be any string of your choice
	URL string

	// Local metrics URL where metrics are fetched from, this could be omitted in the future
	// if implemented using prometheus common/expfmt instead
	MetricsURL string

	// push gateway job name, defaults to "gin"
	Job string
}

func NewProm(subSystem string, customMetrics ...[]*Metrics) *Prom {
	if subSystem == "" {
		subSystem = "app"
	}

	var list []*Metrics

	if len(customMetrics) > 1 {
		panic("Too many args. NewPrometheus( string, <optional []*Metric> ).")
	} else if len(customMetrics) == 1 {
		list = customMetrics[0]
	}
	list = append(list, standardMetrics...)

	p := &Prom{
		List: list,
		Path: metricsPath,
		mappingFn: func(c *gin.Context) string {
			return c.FullPath()
		},
	}

	p.register(subSystem)

	return p
}

func (p *Prom) register(subSystem string) {
	for _, value := range p.List {
		metric := applyMetricsOf(subSystem, value)
		if err := prometheus.Register(metric); err != nil {
			fmt.Println("could not be registered in Prometheus,value.Name:", value.Name, "error:", err.Error())
		}

		switch value {
		case reqCounter:
			p.reqCnt = metric.(*prometheus.CounterVec)
		case reqDuration:
			p.reqDur = metric.(*prometheus.HistogramVec)
		case resSize:
			p.resSz = metric.(prometheus.Summary)
		case reqSize:
			p.reqSz = metric.(prometheus.Summary)
		}

		value.MetricCollector = metric
	}
}

func (p *Prom) SetListenAddress(address string) {
	p.listenAddress = address
	if p.listenAddress != "" {
		p.router = gin.Default()
	}
}

func GetRpcMetrics(registerName string) []prometheus.Collector {
	for _, v := range options.DefaultRpcGroups {
		switch v {
		case "MessageGateway":
			return []prometheus.Collector{OnlineUserGauge}
		case "Msg":
			return []prometheus.Collector{SingleChatMsgProcessSuccessCounter, SingleChatMsgProcessFailedCounter, GroupChatMsgProcessSuccessCounter, GroupChatMsgProcessFailedCounter}
		case "Transfer":
			return []prometheus.Collector{MsgInsertRedisSuccessCounter, MsgInsertRedisFailedCounter, MsgInsertMongoSuccessCounter, MsgInsertMongoFailedCounter, SeqSetFailedCounter}
		case "Push":
			return []prometheus.Collector{MsgOfflinePushFailedCounter}
		case "Auth":
			return []prometheus.Collector{UserLoginCounter}
		default:
			continue
		}
	}
	return nil
}

func GetGinMetrics() []*Metrics {
	return []*Metrics{APICounter}
}

func NewRegistryOf(metrics []prometheus.Collector) (*prometheus.Registry, *grpc_prometheus.ServerMetrics, error) {
	registry := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewServerMetrics()
	grpcMetrics.EnableHandlingTimeHistogram()
	metrics = append(metrics, grpcMetrics, collectors.NewGoCollector())
	registry.MustRegister(metrics...)
	return registry, grpcMetrics, nil
}
