/*
Package ndbc is a pure Go client library for parsing data from NOAA NDBC (National Data Buoy Center)
*/
package ndbc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// BuoyData standard NDBC Data Structure
type BuoyData struct {
	Year               *int
	Month              *int
	Day                *int
	Hour               *int
	Minute             *int
	WindDirection      *int
	WindSpeed          *float64
	WindGust           *float64
	WaveHeight         *float64
	WaveDominantPeriod *float64
	WaveAveragePeriod  *float64
	WaveMeanDirection  *int
	Pressure           *float64
	AirTemp            *float64
	WaterTemp          *float64
	DewPoint           *float64
	Visibility         *float64
	PTDY               *float64
	Tide               *float64
}

// GetBuoyData5Day pulls from NOAA NDBC 5day data structure
func GetBuoyData5Day(buoyID string, maxRecords int) []byte {
	// sample 5 day data
	// http://www.ndbc.noaa.gov/data/5day2/44030_5day.txt

	url := fmt.Sprintf("http://www.ndbc.noaa.gov/data/5day2/%v_5day.txt", buoyID)

	return getBuoyTabularData(url, maxRecords)
}

// GetBuoyData45Day pulls from NOAA NDBC 45day data structure
func GetBuoyData45Day(buoyID string, maxRecords int) []byte {
	// sample 45day data
	// http://www.ndbc.noaa.gov/data/realtime2/44030.txt

	url := fmt.Sprintf("http://www.ndbc.noaa.gov/data/realtime2/%v.txt", buoyID)

	return getBuoyTabularData(url, maxRecords)
}

func getBuoyTabularData(url string, maxRecords int) []byte {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	rawBuoyData, err := ioutil.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	// turn the raw byte array into string
	buoyDataString := string(rawBuoyData[:len(rawBuoyData)])

	// trim the last newline character
	buoyDataString = strings.TrimRight(buoyDataString, "\n")

	// remove first 2 lines because they are just headers for the tabular data
	buoyDataLines := strings.Split(buoyDataString, "\n")[2:]

	// by default get all records
	nRecords := len(buoyDataLines)

	// get only the newest n records if specified
	if maxRecords != 0 && maxRecords < nRecords {
		nRecords = maxRecords
	}

	// create a slice containing the BuoyData structs we want to populate
	buoyDataSet := make([]BuoyData, nRecords)

	for i := range buoyDataSet {
		buoyDataSet[i] = CreateBuoyHourRecord(buoyDataLines[i])
	}

	data, err := json.Marshal(buoyDataSet)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func CreateBuoyHourRecord(row string) BuoyData {
	col := strings.Split(row, " ")
	return BuoyData{
		Year:               getInt(col[0]),
		Month:              getInt(col[1]),
		Day:                getInt(col[2]),
		Hour:               getInt(col[3]),
		Minute:             getInt(col[4]),
		WindDirection:      getInt(col[5]),
		WindSpeed:          getFloat(col[6]),
		WindGust:           getFloat(col[7]),
		WaveHeight:         getFloat(col[8]),
		WaveDominantPeriod: getFloat(col[9]),
		WaveAveragePeriod:  getFloat(col[10]),
		WaveMeanDirection:  getInt(col[11]),
		Pressure:           getFloat(col[12]),
		AirTemp:            getFloat(col[13]),
		WaterTemp:          getFloat(col[14]),
		DewPoint:           getFloat(col[15]),
		Visibility:         getFloat(col[16]),
		PTDY:               getFloat(col[17]),
		Tide:               getFloat(col[18]),
	}
}

func getInt(input string) *int {
	i, err := strconv.Atoi(input)

	if err != nil {
		return nil
	}

	return &i
}

func getFloat(input string) *float64 {
	f, err := strconv.ParseFloat(input, 64)

	if err != nil {
		return nil
	}

	return &f
}
