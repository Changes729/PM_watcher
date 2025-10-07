package manager

import (
	"os"

	"github.com/goccy/go-yaml"
)

var FILE_NAME = "pm-conf.yaml"

type IPDevice struct {
	IP         string `yaml:"ip"`
	Tags       string `yaml:"tags"`
	MultiPower int    `yaml:"multi-power"`
}

type IPDeviceConfig struct {
	IPDevice map[string]IPDevice `yaml:"ip-device"`
}

/** File Operations */
func ReadConfig() (*IPDeviceConfig, error) {
	var config IPDeviceConfig

	_string, err := os.ReadFile(FILE_NAME)
	if err == nil {
		err = yaml.Unmarshal(_string, &config)
	}

	return &config, nil
}

func SaveConfig(config *IPDeviceConfig) (err error) {
	_string, err := yaml.Marshal(config)
	if err == nil {
		err = os.WriteFile(FILE_NAME, _string, 0644)
	}

	return
}
