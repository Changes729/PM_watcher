package main

import (
	"log/slog"
	"main/src/api"
	"main/src/delegate"
	"main/src/manager"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var programLevel = new(slog.LevelVar)
	programLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)

	r := mux.NewRouter()

	manager.YamlInit()
	manager.InitDB()

	delegate.InitMeterConnector(manager.YamlIPDevices())

	api.DeviceSubRouter(r.PathPrefix("/api/").Subrouter())
	api.PowerSubRouter(r.PathPrefix("/api/").Subrouter())

	web := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("./public" + r.URL.Path); os.IsNotExist(err) {
			http.ServeFile(w, r, "./public/index.html")
		} else {
			web.ServeHTTP(w, r)
		}
	})

	if manager.YamlInfo.CROSSupport {
		http.ListenAndServe(":8080", handlers.CORS()(r))
	} else {
		http.ListenAndServe(":8080", r)
	}
}
