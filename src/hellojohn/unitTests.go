package main

import (
	"fmt"
	"io/ioutil"
)

func testAll(testFilesPath string) {
	var erpGantries = getEnrichedGantry()
	var pricingData = getPricingData()

	files, _ := ioutil.ReadDir(testFilesPath)
	for _, f := range files {
		fmt.Println(testFilesPath + "/" + f.Name())
		var filePath = testFilesPath + "/" + f.Name()
		calculatePrice(pricingData, erpGantries, VEHICLE_TYPE, filePath, INPUT_DELIMITER)
	}
}

func printOutAllErpData() {
	fmt.Println("== printing out all erp data in program ===")
	var erpGantries = getEnrichedGantry()
	for index, gantry := range erpGantries {
		fmt.Println(index, gantry.toString())
	}
}
