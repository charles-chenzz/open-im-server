package logger

import (
	"github.com/OpenIMSDK/tools/log"
)

type Options struct {
	StorageLocation     string `yaml:"storageLocation"`
	RotationTime        uint   `yaml:"rotationTime"`
	RemainRotationCount uint   `yaml:"remainRotationCount"`
	RemainLogLevel      int    `yaml:"remainLogLevel"`
	IsStdout            bool   `yaml:"isStdout"`
	IsJson              bool   `yaml:"isJson"`
	WithStack           bool   `yaml:"withStack"`
}

// Apply log configuration to zap
func Apply(option *Options) error {
	err := log.InitFromConfig(
		"im-logs",
		"api",
		option.RemainLogLevel,
		option.IsStdout,
		option.IsJson,
		option.StorageLocation,
		option.RemainRotationCount,
		option.RotationTime,
	)
	if err != nil {
		return err
	}
	return nil
}
