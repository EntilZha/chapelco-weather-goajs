package weather

import (
	"sync"
	"time"

	"code.google.com/r/skirodriguez-dbf/godbf"
)

var cachedDbfTable = new(CachedDbfTable)

const (
	rainSum  = "RAIN_SUM"
	presLoc  = "PRES_LOC"
	presAbs  = "PRES_ABS"
	chn1Deg  = "CHN1_DEG"
	chn1Dew  = "CHN1_DEW"
	chn1Rf   = "CHN1_RF"
	dateTime = "DATE_TIME"
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
		dbfPath := "/Users/pedro/Documents/Code/chapelco-weather/SinusOrg.dbf"
		// dbfPath := "http://googledrive.com/host/0B06ZoNF0o91ncXRPdVRuZjBDaE0"
		cachedDbfTable.DbfTable, err = godbf.NewFromFile(dbfPath, "UTF8")
		cachedDbfTable.updatedAt = time.Now()
		*table = *cachedDbfTable.DbfTable
		cachedDbfTable.Unlock()
	}
	return table, err
}

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

func ReadLastNRainSums(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, rainSum)
}

func ReadLastNPressures(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, presLoc)
}

func ReadLastNAbsPressures(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, presAbs)
}

func ReadLastNTemperatures(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, chn1Deg)
}

func ReadLastNDewPoints(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, chn1Dew)
}

func ReadLastNRelativeHumidities(table *godbf.DbfTable, n int) []float64 {
	return ReadLastNFromFloat64Field(table, n, chn1Rf)
}

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
