package delegate

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"main/src/manager"
	"math/big"
	"net"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

const (
	_CH_ERROR = iota
	_CH_DATA  = iota
)

const (
	_DB_BUCKET = "power"
	_DB_ORG    = "KaYo"
)

type _meterThreadData struct {
	c int
	d manager.DLT_645_2007
}

type ElectricityMeter struct {
	IP   net.IPAddr
	Port int

	_id   string
	_ch   chan _meterThreadData
	_conn net.Conn
}

var _meterList = []ElectricityMeter{}
var _writeAPI api.WriteAPIBlocking = nil

func InitMeterConnector(IParray []string) {
	_writeAPI = manager.InfluxClient.WriteAPIBlocking(_DB_ORG, _DB_BUCKET)

	for _, ip := range IParray {
		newMeter := ElectricityMeter{
			IP:   net.IPAddr{IP: net.ParseIP(ip)},
			Port: 10899,

			_id: manager.EMPTY_ADDRESS,
			_ch: make(chan _meterThreadData),
		}
		_meterList = append(_meterList, newMeter)

		/** connection will be closed on run method. */
		dialAddress := ip + ":" + strconv.Itoa(newMeter.Port)
		conn, err := net.Dial("tcp", dialAddress)
		if err != nil {
			log.Printf("device %s tcp connect failed: %v", newMeter.IP.String(), err)
			continue
		}

		newMeter._conn = conn
		go newMeter.Run()
	}

	for _, meter := range _meterList {
		select {
		case pack := <-meter._ch:
			processPack(pack)
		default:
			// noting todo now.
		}
	}
}

func (e *ElectricityMeter) Run() {
	pack := _meterThreadData{c: _CH_ERROR}

	/** Get Device ID */
	buffer := make([]byte, 1024)
	e._conn.Write(manager.GenCommand(manager.EMPTY_ADDRESS, manager.CC_GET_IP, []byte{}))
	n, err := e._conn.Read(buffer)
	if err != nil {
		log.Printf("device %s command %v read failed: %v", e.IP.String(), manager.CC_GET_IP, err)
	} else if e._id, err = manager.MarshalDeviceAddress(buffer[:n]); err != nil {
		log.Printf("device %s command %v process failed: %v", e.IP.String(), manager.CC_GET_IP, err)
	}

	for range time.Tick(time.Second * 20) {
		if e._id == manager.EMPTY_ADDRESS {
			break
		}

		e._conn.Write(manager.GenCommand(e._id, manager.CC_GET_DATA, manager.FC_COMBINED_ENERGY))
		if n, err := e._conn.Read(buffer); err != nil {
			log.Printf("device %s command %v read failed: %v", e.IP.String(), manager.CC_GET_DATA, err)
		} else if pack.d, err = manager.MarshalData(buffer[:n]); err != nil {
			log.Printf("device %s command %v process failed: %v", e.IP.String(), manager.CC_GET_DATA, err)
		} else {
			pack.c = _CH_DATA
			e._ch <- pack
		}

		e._conn.Write(manager.GenCommand(e._id, manager.CC_GET_DATA, manager.FC_NEUTRAL_LINE_CURRENT))
		if n, err := e._conn.Read(buffer); err != nil {
			log.Printf("device %s command %v read failed: %v", e.IP.String(), manager.CC_GET_DATA, err)
		} else if pack.d, err = manager.MarshalData(buffer[:n]); err != nil {
			log.Printf("device %s command %v process failed: %v", e.IP.String(), manager.CC_GET_DATA, err)
		} else {
			pack.c = _CH_DATA
			e._ch <- pack
		}
	}
}

func processPack(pack _meterThreadData) {
	switch pack.c {
	case _CH_DATA:
		processDLT(pack.d)
	}
}

func processDLT(data manager.DLT_645_2007) {
	switch data.ControlCode {
	case manager.RC_GET_DATA:
	case manager.RC_GET_DATA_N:
	case manager.RC_GET_DATA_E:
		processDLTGetData(data)
	}
}

func processDLTGetData(data manager.DLT_645_2007) {
	var D []byte
	var N []byte

	switch data.ControlCode {
	case manager.RC_GET_DATA:
	case manager.RC_GET_DATA_N:
		D = data.Data[:4]
		N = data.Data[4:]
	case manager.RC_GET_DATA_E:
		log.Printf("err code: %v", data.Data)
		return
	}

	if bytes.Equal(D, manager.FC_COMBINED_ENERGY) {
		bigHex := new(big.Int).SetBytes(N)
		id := hex.EncodeToString(manager.ReverseArray(data.ID))
		power := float32(bigHex.Int64()) / 100.0
		tags := map[string]string{
			"source": "meter",
		}
		data := map[string]interface{}{
			"combined_power": power,
		}
		p := influxdb2.NewPoint(id, tags, data, time.Now())
		err := _writeAPI.WritePoint(context.Background(), p)
		if err != nil {
			log.Printf("DB write data %v failed: %v", data, err)
		}
	} else if bytes.Equal(D, manager.FC_NEUTRAL_LINE_CURRENT) {

	}
}
