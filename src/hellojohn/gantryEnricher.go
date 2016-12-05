package main

import "math"

////////////////////////////////
// ERP Location and Price Join
///////////////////////////////

func getEnrichedGantry() []Gantry {
	var erpGantries = getErpGantries()
	var joinInfomation = getCoordinateToZoneIdMap()

	for key, _ := range erpGantries {
		var erpGantry = erpGantries[key]
		var shortestEuclideanDistance = 9999.0 // abritarily big number
		var bestZoneId = ""
		for erpPoint, erpZoneId := range joinInfomation {
			var candidate = euclideanSum(erpPoint, erpGantry.Line)
			if candidate < shortestEuclideanDistance {
				shortestEuclideanDistance = candidate
				bestZoneId = erpZoneId // erpGantry.ZoneId (zoneId mismatch issue )
			}
		}
		// logging that bestId - erp assignment is done
		//fmt.Println("assigning", bestZoneId, "to erpGantry.ZoneId")
		erpGantry.ZoneId = bestZoneId
		erpGantries[key] = erpGantry
	}

	return erpGantries
}

func euclideanSum(point Point, line ErpLine) float64 {
	return math.Sqrt(math.Pow(point.X-line.HeadX, 2)+math.Pow(point.Y-line.HeadY, 2)) + math.Sqrt(math.Pow(point.X-line.TailX, 2)+math.Pow(point.Y-line.TailY, 2))
}
