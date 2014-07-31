package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"code.google.com/r/skirodriguez-dbf/godbf"
)

type CachedDbfTable struct {
	DbfTable  *godbf.DbfTable
	updatedAt time.Time
	sync.RWMutex
}

type WeatherRecord struct {
	Datetime         time.Time
	LocalPressure    float64
	AbsolutePressure float64
	Temperature      float64
	DewPoint         float64
	RainSum          float64
	RelativeHumidity float64
}

var cachedDbfTable = new(CachedDbfTable)

func currentWeatherHandler(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(readWeatherRecord())
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
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("angular/app")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func readWeatherRecord() *WeatherRecord {
	var err error
	table := new(godbf.DbfTable)
	cachedDbfTable.RLock()
	needsUpdate := cachedDbfTable.updatedAt.Add(time.Minute*20).Before(time.Now()) || cachedDbfTable.DbfTable == nil
	if !needsUpdate {
		*table = *cachedDbfTable.DbfTable
	}
	cachedDbfTable.RUnlock()
	if needsUpdate {
		cachedDbfTable.Lock()
		dbfPath := "http://googledrive.com/host/0B06ZoNF0o91ncXRPdVRuZjBDaE0"
		cachedDbfTable.DbfTable, err = godbf.NewFromUrl(dbfPath, "UTF8")
		cachedDbfTable.updatedAt = time.Now()
		*table = *cachedDbfTable.DbfTable
		cachedDbfTable.Unlock()
	}
	if err != nil {
		return nil
	}
	n := table.NumberOfRecords()
	record := new(WeatherRecord)
	record.RainSum, err = table.Float64FieldValueByName(n-1, "RAIN_SUM")
	record.LocalPressure, err = table.Float64FieldValueByName(n-1, "PRES_LOC")
	record.AbsolutePressure, err = table.Float64FieldValueByName(n-1, "PRES_ABS")
	record.Temperature, err = table.Float64FieldValueByName(n-1, "CHN1_DEG")
	record.DewPoint, err = table.Float64FieldValueByName(n-1, "CHN1_DEW")
	record.RelativeHumidity, err = table.Float64FieldValueByName(n-1, "CHN1_RF")
	var minutes float64
	minutes, err = table.Float64FieldValueByName(n-1, "DATE_TIME")
	record.Datetime = time.Unix(int64((minutes-25569)*86400), 0).UTC()
	return record
}
