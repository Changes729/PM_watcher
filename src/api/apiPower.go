package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"main/src/manager"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type dUnit struct {
	DeviceID        string    `json:"ID"`
	Name            string    `json:"Name"`
	Magnification   int       `json:"Magnification"`
	RawEnergyRecord float32   `json:"Watt"`
	LatestUpdate    time.Time `json:"LatestUpdate"`
}

func PowerSubRouter(r *mux.Router) {
	r.HandleFunc("/power", GetCombinedEnergy).Methods("GET")
	r.HandleFunc("/device-info", UpdateDeviceInfos).Methods("POST")
}

func _CombinedEnergy() (data []dUnit) {
	formattedCmd := fmt.Sprintf(`
	|> range(start: 0)
	|> filter(fn: (r) => r._field == "combined_power")
	|> filter(fn: (r) => r.source == "meter")
	|> last()`)
	result, err := manager.QueryBucket(formattedCmd)
	if err != nil {
		slog.Error(fmt.Sprintf("Request data failed: %v", err))
	} else if result.Err() != nil {
		slog.Error(fmt.Sprintf("query parsing error: %s\n", result.Err().Error()))
	} else {
		// Iterate over query response
		loc, _ := time.LoadLocation("Asia/Shanghai")
		for result.Next() {
			newData := dUnit{}
			newData.DeviceID = result.Record().Measurement()
			meterDevice, _ := manager.YamlMeterDevice(newData.DeviceID)
			newData.Name = meterDevice.Name
			newData.Magnification = meterDevice.MultiPower
			newData.LatestUpdate = result.Record().Time().In(loc)

			value := result.Record().Value()
			switch v := value.(type) {
			case float64:
				newData.RawEnergyRecord = float32(v)
			default:
			}
			data = append(data, newData)
		}
	}

	slog.Debug(fmt.Sprintf("data: %v", data))
	return data
}

func GetCombinedEnergy(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(_CombinedEnergy())
	w.Write(str)
}

func UpdateDeviceInfos(w http.ResponseWriter, r *http.Request) {
	var newDevices []dUnit
	err := json.NewDecoder(r.Body).Decode(&newDevices)
	if err != nil {
		slog.Error(fmt.Sprintf("Decode body failed: %v", err))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	slog.Debug(fmt.Sprintf("New devices: %v", newDevices))

	config := manager.YamlInfo
	for _, device := range newDevices {
		config.MeterDevice[device.DeviceID] = manager.MeterDevice{
			MultiPower: device.Magnification,
			Name:       device.Name,
		}
	}
	err = manager.SaveConfig(config)
	if err != nil {
		slog.Error(fmt.Sprintf("Save config failed: %v", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	manager.ReadConfig()

	w.WriteHeader(http.StatusOK)
}
