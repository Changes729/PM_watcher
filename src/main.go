package main

import (
	"main/src/api"
	"main/src/delegate"
	"main/src/manager"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	manager.InitDB()
	manager.YamlInit()

	delegate.InitMeterConnector(manager.YamlIPDevices())

	api.DeviceSubRouter(r.PathPrefix("/api/").Subrouter())

	web := http.FileServer(http.Dir("./web/"))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("./web" + r.URL.Path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./web/index.html")
		} else {
			web.ServeHTTP(w, r)
		}
	})

	http.ListenAndServe(":8080", r)
}
