package delegate

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"main/src/manager"
	"net"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

const (
	_CH_ERROR = iota
	_CH_DATA  = iota
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

func InitMeterConnector(IParray []string) {
	for _, ip := range IParray {
		slog.Info(fmt.Sprintf("load ip device: %v", ip))
		newMeter := ElectricityMeter{
			IP:   net.IPAddr{IP: net.ParseIP(ip)},
			Port: 10899,

			_id: manager.EMPTY_ADDRESS,
			_ch: make(chan _meterThreadData),
		}
		_meterList = append(_meterList, newMeter)

		/** connection will be closed on run method. */
		dialAddress := ip + ":" + strconv.Itoa(newMeter.Port)
		dial := net.Dialer{Timeout: time.Second * 1}
		conn, err := dial.Dial("tcp", dialAddress)
		if err != nil {
			// FIXME: use log manager. like log.debug etc.
			slog.Error(fmt.Sprintf("device %s tcp connect failed: %v", newMeter.IP.String(), err))
			continue
		}

		newMeter._conn = conn
		go newMeter.Run()
	}

	go func() {
		for {
			for _, meter := range _meterList {
				select {
				case pack := <-meter._ch:
					slog.Debug(fmt.Sprintf("Receive pack: %v", pack))
					processPack(pack)
				default:
					// noting todo now.
				}
			}
		}
	}()
}

func DevicesID() (ids []string) {
	for _, meter := range _meterList {
		ids = append(ids, meter._id)
	}

	return
}

func DeviceID(ip string) (id string) {
	for _, meter := range _meterList {
		if ip == meter.IP.String() {
			id = meter._id
		}
	}

	return
}

func (e *ElectricityMeter) Run() {
	pack := _meterThreadData{c: _CH_ERROR}
	buffer := make([]byte, 0, 1024)
	tmp := make([]byte, 256)
	var err error

	/** Get Device ID */
	slog.Debug(fmt.Sprintf("Get Device %v ID", e.IP.String()))
	e._conn.Write(manager.GenCommand(manager.EMPTY_ADDRESS, manager.CC_GET_IP, []byte{}))
	for isIntegrity, _ := manager.CheckPackIntegrity(buffer); isIntegrity == false; {
		n, err := e._conn.Read(tmp)
		if err != nil {
			slog.Error(fmt.Sprintf("device %s command %v read failed: %v", e.IP.String(), manager.CC_GET_IP, err))
			break
		}

		buffer = append(buffer, tmp[:n]...)
		isIntegrity, _ = manager.CheckPackIntegrity(buffer)
	}

	if e._id, err = manager.MarshalDeviceAddress(buffer); err != nil {
		slog.Error(fmt.Sprintf("device %s command %v process failed: %v", e.IP.String(), manager.CC_GET_IP, err))
	} else {
		slog.Debug(fmt.Sprintf("device address: %v", e._id))
	}

	for range time.Tick(time.Second * time.Duration(manager.YamlInfo.Frequency.IntervalPower)) {
		if e._id == manager.EMPTY_ADDRESS {
			break
		}

		e._conn.Write(manager.GenCommand(e._id, manager.CC_GET_DATA, manager.FC_COMBINED_ENERGY))
		slog.Debug(fmt.Sprintf("Get Device %v energy", e.IP.String()))

		buffer = nil
		for isIntegrity, _ := manager.CheckPackIntegrity(buffer); isIntegrity == false; {
			n, err := e._conn.Read(tmp)
			if err != nil {
				break
			}

			buffer = append(buffer, tmp[:n]...)
			isIntegrity, _ = manager.CheckPackIntegrity(buffer)
		}

		if err != nil {
			slog.Error(fmt.Sprintf("device %s command %v read failed: %v", e.IP.String(), manager.CC_GET_DATA, err))
		} else if pack.d, err = manager.MarshalData(buffer); err != nil {
			slog.Error(fmt.Sprintf("device %s command %v process failed: %v", e.IP.String(), manager.CC_GET_DATA, err))
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
	case manager.RC_GET_DATA,
		manager.RC_GET_DATA_N,
		manager.RC_GET_DATA_E:
		processDLTGetData(data)
	}
}

func processDLTGetData(data manager.DLT_645_2007) {
	var D []byte
	var N []byte

	switch data.ControlCode {
	case manager.RC_GET_DATA,
		manager.RC_GET_DATA_N:
		D = data.Data[:4]
		N = data.Data[4:]
		slog.Debug(fmt.Sprintf("Delegate: %X, %X", D, N))
	case manager.RC_GET_DATA_E:
		slog.Error(fmt.Sprintf("err code: %v", data.Data))
		return
	}

	if bytes.Equal(D, manager.FC_COMBINED_ENERGY) {
		bigHex := hex.EncodeToString(manager.ReverseArray(N))
		slog.Debug(fmt.Sprintf("power raw: %v", bigHex))
		id := hex.EncodeToString(manager.ReverseArray(data.ID))
		power_raw, _ := strconv.Atoi(bigHex)
		power := float32(power_raw) / 100.0
		tags := map[string]string{
			"source": "meter",
		}
		data := map[string]interface{}{
			"combined_power": power,
		}
		p := influxdb2.NewPoint(id, tags, data, time.Now())
		err := manager.WriteAPI.WritePoint(context.Background(), p)
		slog.Debug(fmt.Sprintf("DB write data %f", power))
		if err != nil {
			slog.Error(fmt.Sprintf("DB write data %v failed: %v", data, err))
		}
	} else if bytes.Equal(D, manager.FC_NEUTRAL_LINE_CURRENT) {
	}
}
