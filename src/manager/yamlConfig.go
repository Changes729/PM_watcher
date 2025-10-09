package manager

import (
	"os"
	"log"

	"github.com/goccy/go-yaml"
)

const _FILE_NAME = "./pm-conf.yaml"

type IPDevice struct {
	IP         string `yaml:"ip"`
	Tags       string `yaml:"tags"`
	MultiPower int    `yaml:"multi-power"`
}

type YamlConfig struct {
	IPDevice map[string]IPDevice `yaml:"ip-device"`
}

var YamlInfo YamlConfig

func YamlInit() {
	info, err := ReadConfig()
	if err != nil {
		log.Printf("Read config failed: %v", err)
	} else {
		YamlInfo = info
	}
}

func YamlIPDevices() (devices []string) {
	for _, deviceInfo := range YamlInfo.IPDevice {
		devices = append(devices, deviceInfo.IP)
	}

	log.Printf("Devices: %v", devices)

	return
}

/** File Operations */
func ReadConfig() (YamlConfig, error) {
	var config YamlConfig
	var err error
	_string, err := os.ReadFile(_FILE_NAME)
	if err == nil {
		err = yaml.Unmarshal(_string, &config)
	} else {
		log.Printf("Load file %s err: %v", _FILE_NAME, err)
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
