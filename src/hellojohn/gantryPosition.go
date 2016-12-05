package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

////////////////////////////////
// ERP Location json processing
///////////////////////////////

func getErpGantries() []Gantry {
	// initialize required data structures
	var gantryData = obtainGantryPositionData()
	var gantryPositions = []Gantry{}

	for _, feature := range gantryData.Features {
		// gantry data is long-lat not lat-long therefore: x: key[pos][1] y: key[pos][0]
		var key = feature.Geometry.LineCoordinates

		var newGantry = Gantry{
			HeadX: key[0][1],
			HeadY: key[0][0],
			TailX: key[1][1],
			TailY: key[1][0]}

		newGantry.Line = ErpLine{
			HeadX: key[0][1],
			HeadY: key[0][0],
			TailX: key[1][1],
			TailY: key[1][0]}

		gantryPositions = append(gantryPositions, newGantry)
	}
	fmt.Println("GantryPosition: Loaded ", len(gantryPositions[0:]), "ERP gantries")
	return gantryPositions[0:] // returns a slice
}

func obtainGantryPositionData() GantryData {
	// load json and convert via Unmarshal
	var gantryData GantryData
	raw, _ := ioutil.ReadFile("resources/GantryPosition.json")
	json.Unmarshal(raw, &gantryData)
	return gantryData
}

func (gantry Gantry) toString() string {
	return fmt.Sprint("ZoneId: ", gantry.ZoneId, ", HeadX: ", float64(gantry.HeadX), ", HeadY: ", float64(gantry.HeadY))
	//return "ZoneId: " + gantry.ZoneId + ", HeadX: " + gantry.
	//return "ZoneId: " + gantry.ZoneId + ", HeadX: " + gantry.HeadX + ", HeadY: " + gantry.HeadY + ", TailX: " + gantry.TailX + ", TailY: " + gantry.TailY
}

type Gantry struct {
	HeadX  float64
	HeadY  float64
	TailX  float64
	TailY  float64
	Line   ErpLine
	ZoneId string
}

type ErpLine struct {
	HeadX float64
	HeadY float64
	TailX float64
	TailY float64
}

type GantryData struct {
	Type     string    `json:"type"`
	Crs      Crs       `json:"crs"`
	Features []Feature `json:"features"`
}

type Crs struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"`
}

type Feature struct {
	Type       string            `json:"type"`
	Properties map[string]string `json:"properties"` // use map for these fields
	Geometry   Geometry          `json:"geometry"`
}

type Geometry struct {
	Type            string        `json:"type"`
	LineCoordinates [2][2]float64 `json:"coordinates"`
}

type Line struct {
	Points [][]float64
}
