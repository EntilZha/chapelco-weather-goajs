package weather

import (
	"sync"
	"time"

	"code.google.com/r/skirodriguez-dbf/godbf"
)

var cachedDbfTable = new(CachedDbfTable)

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
		dbfPath := "http://googledrive.com/host/0B06ZoNF0o91ncXRPdVRuZjBDaE0"
		cachedDbfTable.DbfTable, err = godbf.NewFromUrl(dbfPath, "UTF8")
		cachedDbfTable.updatedAt = time.Now()
		*table = *cachedDbfTable.DbfTable
		cachedDbfTable.Unlock()
	}
	return table, err
}

func ReadWeatherRecordFromDbf(table *godbf.DbfTable, n int) *WeatherRecord {
	var err1, err2, err3, err4, err5, err6, err7 error
	record := new(WeatherRecord)
	record.RainSum, err1 = table.Float64FieldValueByName(n, "RAIN_SUM")
	record.LocalPressure, err2 = table.Float64FieldValueByName(n, "PRES_LOC")
	record.AbsolutePressure, err3 = table.Float64FieldValueByName(n, "PRES_ABS")
	record.Temperature, err4 = table.Float64FieldValueByName(n, "CHN1_DEG")
	record.DewPoint, err5 = table.Float64FieldValueByName(n, "CHN1_DEW")
	record.RelativeHumidity, err6 = table.Float64FieldValueByName(n, "CHN1_RF")
	var minutes float64
	minutes, err7 = table.Float64FieldValueByName(n, "DATE_TIME")
	record.Datetime = time.Unix(int64((minutes-25569)*86400), 0).UTC()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		return nil
	}
	return record
}

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

func ReadCurrentWeatherRecord() *WeatherRecord {
	table, err := getDbf()
	if err != nil {
		return nil
	}
	n := table.NumberOfRecords() - 1
	return ReadWeatherRecordFromDbf(table, n)
}

func ReadLastNWeatherRecords(n int) []WeatherRecord {
	table, err := getDbf()
	if err != nil {
		return nil
	}
	return ReadLastNWeatherRecordsFromDbf(table, n)
}
