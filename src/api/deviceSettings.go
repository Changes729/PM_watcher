package api

import (
	"encoding/json"
	"main/src/manager"
	"net/http"

	"github.com/gorilla/mux"
)

func DeviceSubRouter(r *mux.Router) {
	r.HandleFunc("/device", GetDevice).Methods("GET")
	r.HandleFunc("/device/{id}", GetDeviceInfo).Methods("GET")
	r.HandleFunc("/device/{id}", UpdateDevice).Methods("POST")
	r.HandleFunc("/device/{id}", RemoveDevice).Methods("DELETE")
}

func GetDevice(w http.ResponseWriter, r *http.Request) {
	var count struct {
		Ids []string `json:"ids"`
	}

	for key := range manager.YamlInfo.IPDevice {
		count.Ids = append(count.Ids, key)
	}

	str, _ := json.Marshal(count)
	w.Write(str)
}

func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if device, ok := manager.YamlInfo.IPDevice[id]; ok {
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
		manager.YamlInfo.IPDevice[id] = device
		w.WriteHeader(http.StatusOK)
	}
}

func RemoveDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, ok := manager.YamlInfo.IPDevice[id]; ok {
		delete(manager.YamlInfo.IPDevice, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
