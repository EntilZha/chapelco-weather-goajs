// Copyright 2014 Pedro Rodriguez. All rights reserved.
// Use of this code is governed by the MIT License

// Provides interface to weather data from Chapelco Ski Resort weather station via Google Drive
package weather

import (
	"sync"
	"time"

	"code.google.com/r/skirodriguez-dbf/godbf"
)

// cachedDbfTable holds cached DbfTable. It is only fetched every 20 minutes after it is stale.
var cachedDbfTable = new(CachedDbfTable)

// Constants to access variables from Chapelco weather .dbf file
const (
	rainSum  = "RAIN_SUM"
	presLoc  = "PRES_LOC"
	presAbs  = "PRES_ABS"
	chn1Deg  = "CHN1_DEG"
	chn1Dew  = "CHN1_DEW"
	chn1Rf   = "CHN1_RF"
	dateTime = "DATE_TIME"
)

// CachedDbfTable consists of DbfTable which holds a godbf.DfTable, updatedAt contains the time.Time it was
// last updated, and holds a Read/Write lock to insure that the table is in sync with when it was last updated.
type CachedDbfTable struct {
	DbfTable  *godbf.DbfTable
	updatedAt time.Time
	sync.RWMutex
}

// WeatherRecord represents one weather observation from the weather station
type WeatherRecord struct {
	Datetime         time.Time
	LocalPressure    float64
	AbsolutePressure float64
	Temperature      float64
	DewPoint         float64
	RainSum          float64
	RelativeHumidity float64
}

// getDbf returns a pointer to a dbf table from a constant defined url. On first call it fetches the table, thereafter
// returns the value from the cached table unless it is stale by 20 minutes or more.
func getDbf() (*godbf.DbfTable, error) {
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
		//dbfPath := "/Users/pedro/Downloads/chapelco.dbf"
		dbfPath := "http://googledrive.com/host/0B06ZoNF0o91ncXRPdVRuZjBDaE0"
		cachedDbfTable.DbfTable, err = godbf.NewFromUrl(dbfPath, "UTF8")
		cachedDbfTable.updatedAt = time.Now()
		*table = *cachedDbfTable.DbfTable
		cachedDbfTable.Unlock()
	}
	return table, err
}

// ReadWeatherRecordFromDbf reads a single WeatherRecord from the given Dbf Table.
func ReadWeatherRecordFromDbf(table *godbf.DbfTable, n int) *WeatherRecord {
	var err1, err2, err3, err4, err5, err6, err7 error
	record := new(WeatherRecord)
	record.RainSum, err1 = table.Float64FieldValueByName(n, rainSum)
	record.LocalPressure, err2 = table.Float64FieldValueByName(n, presLoc)
	record.AbsolutePressure, err3 = table.Float64FieldValueByName(n, presAbs)
	record.Temperature, err4 = table.Float64FieldValueByName(n, chn1Deg)
	record.DewPoint, err5 = table.Float64FieldValueByName(n, chn1Dew)
	record.RelativeHumidity, err6 = table.Float64FieldValueByName(n, chn1Rf)
	var minutes float64
	minutes, err7 = table.Float64FieldValueByName(n, dateTime)
	record.Datetime = time.Unix(int64((minutes-25569.0)*86400.0)+60*60*4, 0).UTC()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		return nil
	}
	return record
}

// ReadLastNWeatherRecordsFromDbf reads the last n WeatherRecords from the DbfTable
func ReadLastNWeatherRecordsFromDbf(table *godbf.DbfTable, n int) []WeatherRecord {
	total := table.NumberOfRecords()
	start := total - n
	if start < 0 {
		return nil
	}
	records := make([]WeatherRecord, n)
	for i := 0; i < n; i++ {
		r := ReadWeatherRecordFromDbf(table, i+start)
		if r == nil {
			return nil
		}
		records[i] = *r
	}
	return records
}

// ReadCurrentWeatherRecord reads the most recent (last 1) WeatherRecord from the DbfTable
func ReadCurrentWeatherRecord() *WeatherRecord {
	table, err := getDbf()
	if err != nil {
		return nil
	}
	n := table.NumberOfRecords() - 1
	return ReadWeatherRecordFromDbf(table, n)
}

// ReadLastNWeatherRecords reads the last n records from the cached DbfTable
func ReadLastNWeatherRecords(n int) []WeatherRecord {
	table, err := getDbf()
	if err != nil {
		return nil
	}
	return ReadLastNWeatherRecordsFromDbf(table, n)
}

// ReadLastNWeatherRecordsToMap reads the last n records in separate lists into a map with keys from code.
func ReadLastNWeatherRecordsToMap(n int) map[string]interface{} {
	table, err := getDbf()
	if err != nil {
		return nil
	}
	fields := make(map[string]interface{})
	fields[rainSum] = ReadLastNRainSums(table, n)
	fields[presLoc] = ReadLastNPressures(table, n)
	fields[presAbs] = ReadLastNAbsPressures(table, n)
	fields[chn1Deg] = ReadLastNTemperatures(table, n)
	fields[chn1Dew] = ReadLastNDewPoints(table, n)
	fields[chn1Rf] = ReadLastNRelativeHumidities(table, n)
	fields[dateTime] = ReadLastNDateTimes(table, n)
	return fields
}

// ReadLastNRainSums reads the last n RAIN_SUM records
func ReadLastNRainSums(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, rainSum)
}

// ReadLastNPressures reads the last n PRES_LOC records
func ReadLastNPressures(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, presLoc)
}

// ReadLastNAbsPressures reads the last n PRES_ABS records
func ReadLastNAbsPressures(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, presAbs)
}

// ReadLastNTemperatures reads the last n CHN1_DEG records
func ReadLastNTemperatures(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, chn1Deg)
}

// ReadLastNDewPoints reads the last n CHN1_DEW records
func ReadLastNDewPoints(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, chn1Dew)
}

// ReadLastNRelativeHumidities reads the last n CHN1_RF records
func ReadLastNRelativeHumidities(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, chn1Rf)
}

// ReadLastNFromFloat64Field reads the last n records by field string
func ReadLastNFromFloat64Field(table *godbf.DbfTable, n int, field string) []float64 {
	rows := make([]float64, n)
	var err error
	start := table.NumberOfRecords() - n
	for i := 0; i < n; i++ {
		rows[i], err = table.Float64FieldValueByName(i+start, field)
		if err != nil {
			return nil
		}
	}
	return rows
}

// ReadLastNDateTimes reads the last n DATE_TIME records
func ReadLastNDateTimes(table *godbf.DbfTable, n int) []string {
	rows := make([]string, n)
	start := table.NumberOfRecords() - n
	for i := 0; i < n; i++ {
		rawVal, err := table.Float64FieldValueByName(i+start, dateTime)
		seconds := int64((rawVal - 25569) * 86400)
		rows[i] = time.Unix(seconds, 0).UTC().Format("1/2 15:04")
		if err != nil {
			return nil
		}
	}
	return rows
}
