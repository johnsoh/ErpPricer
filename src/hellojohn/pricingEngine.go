package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/////////////////////////////
// ERP Price Data processing
////////////////////////////

func (pricingData PricingData) getPrice(vehicleType string, zoneId string, dateTime string) float64 {

	var dateTimeArr = strings.Split(dateTime, " ")
	var form = "2006-01-02"
	var t, _ = time.Parse(form, dateTimeArr[0])

	var dayType string
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		dayType = "Weekends"
	} else {
		dayType = "Weekdays"
	}

	var intTime = convertTimeStringToInt(dateTimeArr[1])

	var pricings = pricingData.VehicleToZoneId[vehicleType].ZoneIdToDayType[zoneId].DayTypeToTimePrice[dayType]
	for _, priceTimeInfo := range pricings {
		if priceTimeInfo.StartTime <= intTime && intTime <= priceTimeInfo.EndTime {
			return priceTimeInfo.Price
		}
	}
	return 0.0
}

func getPricingData() PricingData {
	var pricingDataJson = getPricingDataJsonObjectFromJson()
	var pricingData = PricingData{VehicleToZoneId: make(map[string]ZoneIdData)}
	for _, pricingDataEntry := range pricingDataJson.PricingDataEntries {
		chargeAmount, _ := strconv.ParseFloat(pricingDataEntry["ChargeAmount"], 64)
		if chargeAmount == 0 {
			continue
		}

		var vehicleType = pricingDataEntry["VehicleType"]
		var dayType = pricingDataEntry["DayType"]
		var startTime = convertTimeStringToInt(pricingDataEntry["StartTime"])
		var endTime = convertTimeStringToInt(pricingDataEntry["EndTime"])
		var zoneID = pricingDataEntry["ZoneID"]

		zoneIdObject, ok := pricingData.VehicleToZoneId[vehicleType]
		if !ok {
			zoneIdObject = ZoneIdData{}
		}

		dayTypeObject, dayTypeObjectIsPresent := zoneIdObject.ZoneIdToDayType[dayType]
		if !dayTypeObjectIsPresent {
			dayTypeObject = DayType{}
		}

		var timePrice = TimePrice{
			StartTime: startTime,
			EndTime:   endTime,
			Price:     chargeAmount}

		dayTypeObject.DayTypeToTimePrice[dayType] = append(dayTypeObject.DayTypeToTimePrice[dayType], timePrice)
		zoneIdObject.ZoneIdToDayType[zoneID] = dayTypeObject
		pricingData.VehicleToZoneId[vehicleType] = zoneIdObject
	}
	return pricingData
}

func convertTimeStringToInt(timeString string) int64 {
	arr := strings.Split(timeString, ":")
	value, _ := strconv.ParseInt(arr[0]+arr[1], 10, 64)
	return value
}

func getPricingDataJsonObjectFromJson() PricingDataJson {
	// load json and convert via Unmarshal
	var pricingDataJson PricingDataJson
	raw, _ := ioutil.ReadFile("resources/resources/PriceTimingZoneId.json")
	json.Unmarshal(raw, &pricingDataJson)
	return pricingDataJson
}

type PricingDataJson struct {
	PricingDataEntries []map[string]string `json:"features"`
}

type PricingData struct {
	VehicleToZoneId map[string]ZoneIdData
}

type ZoneIdData struct {
	ZoneIdToDayType map[string]DayType
}

type DayType struct {
	DayTypeToTimePrice map[string][]TimePrice
}

type TimePrice struct {
	StartTime int64
	EndTime   int64
	Price     float64
}

////////////////////////////
// ERP Price Data retrieval
///////////////////////////

func getLatestErpAndCache() {
	// Prepare request
	erpUrl := "http://datamall2.mytransport.sg/ltaodataservice/ERPRates"
	req, _ := http.NewRequest("GET", erpUrl, nil)
	req.Header.Add("AccountKey", "p33ITka7RQ+Y4fQQgdDJYQ==")
	req.Header.Add("UniqueUserID", "e724fbb1-6f15-4aa5-96e9-26ba4db92798")
	req.Header.Add("accept", "application/json")

	// Send request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, _ := client.Do(req)

	// Write to file
	defer resp.Body.Close()
	out, _ := os.Create("resources/PriceTimingZoneIdNew.json")
	defer out.Close()
	io.Copy(out, resp.Body)
}
