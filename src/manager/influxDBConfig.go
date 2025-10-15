package manager

import (
	"context"
	"fmt"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var InfluxClient influxdb2.Client
var WriteAPI api.WriteAPIBlocking = nil
var queryAPI api.QueryAPI = nil

func InitDB() {
	InfluxClient = influxdb2.NewClient(
		YamlInfo.InfluxSetting.Url, YamlInfo.InfluxSetting.Token)

	WriteAPI = InfluxClient.WriteAPIBlocking(
		YamlInfo.InfluxSetting.Org, YamlInfo.InfluxSetting.Bucket)

	queryAPI = InfluxClient.QueryAPI(YamlInfo.InfluxSetting.Org)
}

func QueryBucket(cmd string) (*api.QueryTableResult, error) {
	fullCmd := fmt.Sprintf(`
	from(bucket: "%s") %s
	`, YamlInfo.InfluxSetting.Bucket, cmd)

	log.Printf("bucket: %s db cmd: %s", YamlInfo.InfluxSetting.Bucket, fullCmd)
	return queryAPI.Query(context.Background(), fullCmd)
}
