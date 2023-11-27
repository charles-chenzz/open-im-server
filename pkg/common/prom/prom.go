package prom

type Monitor struct {
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
