package manager

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var InfluxClient influxdb2.Client

const (
	_url   = "http://localhost:8086"
	_token = "6fvvj9Q5pnx9X5opO6t4emyJuwn1BgUDVb6KwDIAOaLlQljRf_AAvQCdS8xtP9M1sku8ZueF8VgnKn2dRIRNRQ=="
)

func InitDB() {
	InfluxClient = influxdb2.NewClient(_url, _token)
}
