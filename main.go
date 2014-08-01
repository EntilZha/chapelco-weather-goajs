package main

import (
	"encoding/json"
	"net/http"
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

func pastWeatherHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/weather/current", currentWeatherHandler)
	router.HandleFunc("/api/weather/past/{n}", pastWeatherHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("angular/app")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
