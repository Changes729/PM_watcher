package manager

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var InfluxClient influxdb2.Client
var WriteAPI api.WriteAPIBlocking = nil

func InitDB() {
	InfluxClient = influxdb2.NewClient(
		YamlInfo.influxSetting.url, YamlInfo.influxSetting.token)

	WriteAPI = InfluxClient.WriteAPIBlocking(
		YamlInfo.influxSetting.org, YamlInfo.influxSetting.bucket)
}
