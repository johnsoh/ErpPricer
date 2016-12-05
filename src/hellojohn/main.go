package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CONSTANTS
const INPUT_FILE_PATH = "resources/erp16nov/a2e77e0b-ab82-11e6-8129-aef751f686cf.csv" // "resources/test_vehicle.csv"
const INPUT_DELIMITER = "|"                                                           // "," for initial dataset
const VEHICLE_TYPE = "Very Heavy Goods Vehicles/Big Buses"

func main() {

	testAll("resources/erp16nov")
	printOutAllErpData()

	// init and obtain erp positions, pricing data
	//calculatePrice(getPricingData(), getEnrichedGantry(), VEHICLE_TYPE, INPUT_FILE_PATH, INPUT_DELIMITER)
}

func calculatePrice(pricingData PricingData, erpGantries []Gantry, vehicleType string, inputFilePath string, inputDelimiter string) {
	// prepare vehicle position scanner
	file, _ := os.Open(inputFilePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fmt.Println("Loaded vehicle data with headers:", scanner.Text())

	// algorithm start
	var vehicleCost = 0.0
	scanner.Scan()
	_, prevX, prevY := getDataFromCsvLine(scanner.Text(), inputDelimiter)

	for scanner.Scan() {
		// get coordinates, do calc, reassign
		dateTimeString, currX, currY := getDataFromCsvLine(scanner.Text(), inputDelimiter)
		var res, gantry = doesRouteIntersectWithGantry(prevX, prevY, currX, currY, erpGantries)
		if res {
			fmt.Println(dateTimeString, ": intersection detected alogn path ", prevX, ",", prevY, "to", currX, ",", currY, "(", gantry.ZoneId, ")")
			vehicleCost = vehicleCost + pricingData.getPrice(vehicleType, gantry.ZoneId, dateTimeString)
		}
		prevX, prevY = currX, currY
	}
	fmt.Println("Total cost:", vehicleCost)
}

func getDataFromCsvLine(line string, delimiter string) (string, float64, float64) {
	data := strings.Split(line, delimiter)
	x, _ := strconv.ParseFloat(data[1], 64)
	y, _ := strconv.ParseFloat(data[2], 64)

	return data[0], x, y
}

func doesRouteIntersectWithGantry(prevX float64, prevY float64, currX float64, currY float64, points []Gantry) (bool, Gantry) {

	for _, gantry := range points {
		if doesIntersect(prevX, prevY, currX, currY, gantry.HeadX, gantry.HeadY, gantry.TailX, gantry.TailY) {
			return true, gantry
		}
	}
	return false, Gantry{}
}

func doesIntersect(line1HeadX float64, line1HeadY float64, line1TailX float64, line1TailY float64, line2HeadX float64, line2HeadY float64, line2TailX float64, line2TailY float64) bool {
	return (line1HeadX <= line2HeadX && line1HeadY <= line2HeadY) && (line1TailX >= line2TailX && line1TailY >= line2TailY)
}
