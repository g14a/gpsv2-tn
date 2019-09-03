package utils

import (
	"math"
)

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

func Distance(prevlat, prevlong, currlat, currlong float64) (km float64) {

	if prevlat == currlat && prevlong == currlong {
		return 0
	}

	if prevlat == 0 || currlat == 0 || prevlong == 0 || currlong == 0 {
		return 0
	}

	lat1 := degreesToRadians(prevlat)
	lon1 := degreesToRadians(prevlong)
	lat2 := degreesToRadians(currlat)
	lon2 := degreesToRadians(currlong)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	km = c * 6371

	km = math.Round(km*100/100)

	return km
}
