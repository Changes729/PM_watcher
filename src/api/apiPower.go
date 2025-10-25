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
}

func _CombinedEnergy() (data []dUnit) {
	formattedCmd := fmt.Sprintf(`
	|> range(start: 0)
	|> filter(fn: (r) => r._field == "combined_power")
	|> filter(fn: (r) => r.source == "meter")
	|> first()`)
	result, err := manager.QueryBucket(formattedCmd)
	if err != nil {
		slog.Error(fmt.Sprintf("Request data failed: %v", err))
	} else if result.Err() != nil {
		slog.Error(fmt.Sprintf("query parsing error: %s\n", result.Err().Error()))
	} else {
		// Iterate over query response
		newData := dUnit{}
		for result.Next() {
			newData.DeviceID = result.Record().Measurement()
			meterDevice, _ := manager.YamlMeterDevice(newData.DeviceID)
			newData.Name = meterDevice.Name
			newData.Magnification = meterDevice.MultiPower
			newData.LatestUpdate = result.Record().Time()

			value := result.Record().Value()
			switch v := value.(type) {
			case float64:
				newData.RawEnergyRecord = float32(v)
			default:
			}
		}
		data = append(data, newData)
	}

	slog.Debug(fmt.Sprintf("data: %v", data))
	return data
}

func GetCombinedEnergy(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(_CombinedEnergy())
	w.Write(str)
}
