package utils

import (
	"math"
)

func haversine(theta float64) float64 {
	return .5 * (1 - math.Cos(theta))
}

type pos struct {
	lat  float64 // latitude, radians
	long float64 // longitude, radians
}

func degPos(lat, lon float64) pos {
	return pos{lat * math.Pi / 180, lon * math.Pi / 180}
}

const rEarth = 6372.8 // km

func hsDist(p1, p2 pos) float64 {
	return 2 * rEarth * math.Asin(math.Sqrt(haversine(p2.lat-p1.lat)+
		math.Cos(p1.lat)*math.Cos(p2.lat)*haversine(p2.long-p1.long)))
}
