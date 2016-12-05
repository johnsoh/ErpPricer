package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// going to use zoneIds + list
func getCoordinateToZoneIdMap() map[Point]string {

	// get array of zoneId-SgCoordinates information
	var serviceInfo = obtainServicesInfo()
	var erpInfo = serviceInfo.ErpInfoArray[0]
	var arrayOfJoinInformation = erpInfo["ERPINFO"]

	// prepare to fill in zoneToId
	var pointToZoneIdMap = make(map[Point]string)
	var convertedCoordinates, _ = readLines("resources/convertedCoordinates.csv")
	var pointer = 0

	for _, joinInformation := range arrayOfJoinInformation {
		// check that this map contains zoneId-Location information
		zoneId, isPresent := joinInformation["ZONEID"]
		if !isPresent {
			continue
		}

		var convertedArray = strings.Split(convertedCoordinates[pointer], ",")
		var x, _ = strconv.ParseFloat(convertedArray[0], 64)
		var y, _ = strconv.ParseFloat(convertedArray[1], 64)
		pointer = pointer + 1

		var point = Point{X: x, Y: y}
		pointToZoneIdMap[point] = zoneId
	}
	return pointToZoneIdMap
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func obtainServicesInfo() ServicesInfo {
	// load json and convert via Unmarshal
	var zoneIdCoordinatesJoinObject ServicesInfo
	raw, _ := ioutil.ReadFile("resources/zoneIdCoordinatesJoin.json")
	json.Unmarshal(raw, &zoneIdCoordinatesJoinObject)
	return zoneIdCoordinatesJoinObject
}

type ServicesInfo struct {
	ErpInfoArray []map[string][]map[string]string `json:"SERVICESINFO"`
}

type ErpInfo struct {
	vals []map[string]string
}

func convertSvyToWrs(svyCoordinateX int, svyCoordinateY int) []float64 {
	var arr = [2]float64{1.0, 2.0}
	return arr[0:]
}

type Point struct {
	X float64
	Y float64
}

type Coordinate struct {
	HeadX float64
	HeadY float64
	TailX float64
	TailY float64
}

type ZoneIdCoordinate struct {
	HeadX  float64
	HeadY  float64
	TailX  float64
	TailY  float64
	ZoneId string
}
