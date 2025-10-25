package manager

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

const _FILE_NAME = "./pm-conf.yaml"

type IPDevice struct {
	LatestID string `yaml:"latest-id"`
}

type InfluxSettings struct {
	Url    string `yaml:"url"`
	Token  string `yaml:"token"`
	Bucket string `yaml:"bucket"`
	Org    string `yaml:"org"`
}

type Frequency struct {
	IntervalPower int `yaml:"power"`
}

type MeterDevice struct {
	MultiPower int    `yaml:"multi-power"`
	Name       string `yaml:"name"`
}

type YamlConfig struct {
	Frequency     Frequency              `yaml:"frequency"`
	InfluxSetting InfluxSettings         `yaml:"influxDB"`
	IPDevice      map[string]IPDevice    `yaml:"ip-device"`
	MeterDevice   map[string]MeterDevice `yaml:"meter-device"`
}

var YamlInfo YamlConfig

func YamlInit() {
	info, err := ReadConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("Read config failed: %v", err))
	} else {
		YamlInfo = info
		slog.Debug(fmt.Sprintf("Read config: %v", info))
	}
}

func YamlIPDevices() (devices []string) {
	for k := range YamlInfo.IPDevice {
		devices = append(devices, k)
	}

	slog.Debug(fmt.Sprintf("Devices: %v", devices))

	return
}

func YamlMeterDevice(id string) (device MeterDevice, ok bool) {
	device, ok = YamlInfo.MeterDevice[id]
	return
}

/** File Operations */
func ReadConfig() (YamlConfig, error) {
	var config YamlConfig
	var err error
	_string, err := os.ReadFile(_FILE_NAME)
	if err == nil {
		err = yaml.Unmarshal(_string, &config)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("Load file %s err: %v", _FILE_NAME, err))
	} else {
		if config.Frequency.IntervalPower == 0 {
			config.Frequency.IntervalPower = 20
		}
	}

	return config, err
}

func SaveConfig(config YamlConfig) (err error) {
	_string, err := yaml.Marshal(config)
	if err == nil {
		err = os.WriteFile(_FILE_NAME, _string, 0644)
	}

	return
}
