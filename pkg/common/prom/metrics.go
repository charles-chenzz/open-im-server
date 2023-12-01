package prom

import "github.com/prometheus/client_golang/prometheus"

var (
	reqCounter = &Metrics{
		ID:          "reqCnt",
		Name:        "requests_total",
		Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		Type:        "counter_vec",
		Args:        []string{"code", "method", "handler", "host", "url"}}

	reqDuration = &Metrics{
		ID:          "reqDur",
		Name:        "request_duration_seconds",
		Description: "The HTTP request latencies in seconds.",
		Type:        "histogram_vec",
		Args:        []string{"code", "method", "url"},
	}

	resSize = &Metrics{
		ID:          "resSz",
		Name:        "response_size_bytes",
		Description: "The HTTP response sizes in bytes.",
		Type:        "summary"}

	reqSize = &Metrics{
		ID:          "reqSz",
		Name:        "request_size_bytes",
		Description: "The HTTP request sizes in bytes.",
		Type:        "summary"}

	standardMetrics = []*Metrics{
		reqCounter,
		reqDuration,
		resSize,
		reqSize,
	}

	// custom metrics
	APICounter = &Metrics{
		Name:        "total",
		Description: "counter events",
		Type:        "counter_vec",
		Args:        []string{"label_one", "label_two"},
	}
)

var (
	UserLoginCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "user_login_total",
		Help: "The number of user login",
	})

	SingleChatMsgProcessSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_success_total",
		Help: "The number of single chat msg successful processed",
	})

	SingleChatMsgProcessFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "single_chat_msg_process_failed_total",
		Help: "The number of single chat msg failed processed",
	})

	GroupChatMsgProcessSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_success_total",
		Help: "The number of group chat msg successful processed",
	})

	GroupChatMsgProcessFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "group_chat_msg_process_failed_total",
		Help: "The number of group chat msg failed processed",
	})

	MsgOfflinePushFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_offline_push_failed_total",
		Help: "The number of msg failed offline pushed",
	})

	MsgInsertRedisSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_success_total",
		Help: "The number of successful insert msg to redis",
	})

	MsgInsertRedisFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_redis_failed_total",
		Help: "The number of failed insert msg to redis",
	})

	MsgInsertMongoSuccessCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_success_total",
		Help: "The number of successful insert msg to mongo",
	})

	MsgInsertMongoFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "msg_insert_mongo_failed_total",
		Help: "The number of failed insert msg to mongo",
	})

	SeqSetFailedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "seq_set_failed_total",
		Help: "The number of failed set seq",
	})

	OnlineUserGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "online_user_num",
		Help: "The number of online user num",
	})
)

func applyMetricsOf(subSystem string, m *Metrics) (metrics prometheus.Collector) {
	switch m.Type {
	case "counter_vec":
		metrics = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metrics = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metrics = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metrics = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metrics = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "histogram":
		metrics = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metrics = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metrics = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subSystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metrics
}
