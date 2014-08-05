package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/EntilZha/chapelco-weather-goajs/weather"

	"github.com/gorilla/mux"
)

func currentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(weather.ReadCurrentWeatherRecord())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func pastWeatherRecordsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	n, err := strconv.Atoi(params["n"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(weather.ReadLastNWeatherRecords(n))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func pastWeatherListsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	n, err := strconv.Atoi(params["n"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(weather.ReadLastNWeatherRecordsToMap(n))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/weather/current", currentWeatherHandler)
	router.HandleFunc("/api/weather/past-record-list/{n}", pastWeatherRecordsHandler)
	router.HandleFunc("/api/weather/past-field-lists/{n}", pastWeatherListsHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("angular/app")))
	http.Handle("/", router)
	http.ListenAndServe(getPort(), nil)
}
