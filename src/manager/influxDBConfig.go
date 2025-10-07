package manager

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var InfluxClient influxdb2.Client

const (
	_url   = "http://localhost:8086"
	_token = "token"
)

func InitDB() {
	InfluxClient = influxdb2.NewClient(_url, _token)
}
