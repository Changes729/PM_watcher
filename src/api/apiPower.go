package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/src/delegate"
	"main/src/manager"
	"net/http"

	"github.com/gorilla/mux"
)

type dUnit struct {
	DeviceID        string  `json:"ID"`
	Tags            string  `json:"Tags"`
	Multiply        int     `json:"Multiply"`
	RawEnergyRecord float32 `json:"RawEnergyRecord"`
	EnergyRecord    float32 `json:"EnergyRecord"`
}

func PowerSubRouter(r *mux.Router) {
	r.HandleFunc("/power", GetCombinedEnergy).Methods("GET")
	r.HandleFunc("/power/csv", CombinedEnergyCsvFile).Methods("GET")
}

func _CombinedEnergy() (data []dUnit) {
	for _, val := range manager.YamlInfo.IPDevice {
		newData := dUnit{
			DeviceID:        delegate.DeviceID(val.IP),
			Tags:            val.Tags,
			Multiply:        val.MultiPower,
			RawEnergyRecord: 0,
		}

		// formattedCmd := fmt.Sprintf(`
		// |> filter(fn: (r) => r._measurement == "%s")
		// |> filter(fn: (r) => r.source == "meter")
		// |> first()`, newData.DeviceID)
		formattedCmd := fmt.Sprintf(`
    |> range(start: 0)
    |> filter(fn: (r) => r._measurement == "%s")
    |> filter(fn: (r) => r.machine == "A")
		|> first()`, "electronic_power")
		result, err := manager.QueryBucket(formattedCmd)
		if err != nil {
			log.Printf("Request data failed: %v", err)
		} else if result.Err() != nil {
			log.Printf("query parsing error: %s\n", result.Err().Error())
		} else {
			// Iterate over query response
			for result.Next() {
				value := result.Record().Value()
				switch v := value.(type) {
				case float64:
					newData.RawEnergyRecord = float32(v)
					newData.EnergyRecord = newData.RawEnergyRecord * float32(newData.Multiply)
				default:
				}
			}
			data = append(data, newData)
		}
	}

	log.Printf("data: %v", data)
	return data
}

func GetCombinedEnergy(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(_CombinedEnergy())
	w.Write(str)
}

func CombinedEnergyCsvFile(w http.ResponseWriter, r *http.Request) {
	csvFile := "ID,Tags,Multiply,RawEnergyRecord\n"

	for _, node := range _CombinedEnergy() {
		csvFile += fmt.Sprintf("%s,%s,%d,%f\n", node.DeviceID, node.Tags, node.Multiply, node.RawEnergyRecord)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=power.csv")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len([]byte(csvFile))))
	_, err := io.Copy(w, bytes.NewReader([]byte(csvFile)))
	if err != nil {
		http.Error(w, "Failed to write file content", http.StatusInternalServerError)
		return
	}
}
