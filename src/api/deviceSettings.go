package api

import (
	"encoding/json"
	"main/src/manager"
	"net/http"

	"github.com/gorilla/mux"
)

var Config *manager.IPDeviceConfig = nil

func DeviceSubRouter(r *mux.Router) {
	Config, _ = manager.ReadConfig()

	r.HandleFunc("/device", GetDevice).Methods("GET")
	r.HandleFunc("/device/{id}", GetDeviceInfo).Methods("GET")
	r.HandleFunc("/device/{id}", UpdateDevice).Methods("POST")
	r.HandleFunc("/device/{id}", RemoveDevice).Methods("DELETE")
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	var count struct {
		Ids []string `json:"ids"`
	}

	for key := range Config.IPDevice {
		count.Ids = append(count.Ids, key)
	}

	str, _ := json.Marshal(count)
	w.Write(str)
}

func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if device, ok := Config.IPDevice[id]; ok {
		str, _ := json.Marshal(device)
		w.Write(str)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

// TODO: write to file
func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var device manager.IPDevice
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		Config.IPDevice[id] = device
		w.WriteHeader(http.StatusOK)
	}
}

func RemoveDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, ok := Config.IPDevice[id]; ok {
		delete(Config.IPDevice, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
