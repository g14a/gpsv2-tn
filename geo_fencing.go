package main

import (
	"fmt"
	"math"
)
func main() {
	//chennai := haversine.Coord{Lat: 51.5007, Lon: 0.1246}  // Oxford, UK
	//kandigai  := haversine.Coord{Lat: 40.6892, Lon: 74.0445}  // Turin, Italy
	//mi, km := haversine.Distance(chennai, kandigai)
	//fmt.Println("Miles:", mi, "Kilometers:", km)

	device := &device {
		radius: 2,
		lat: 3,
		long: 2,
		fexist:false,
		farray: [][]float64{{2, 3}, {4, 5}},
	}

	for index, item := range device.farray {
		fmt.Println(item[0], item[1])
		if device.isInside(item[0], item[1]) {
			device.fexist = true
			device.currentfno = string(index)
		}
	}

	fmt.Println(device)
}

type device struct {
	radius float64
	lat float64
	long float64
	currentfno string
	fexist bool
	farray [][]float64
}

func (d *device) isInside(circlex, circley float64) bool {
	xpart := math.Pow(d.lat-circlex, 2)
	ypart := math.Pow(d.long - circley, 2)
	radiusSquared := math.Pow(d.radius, 2)

	if xpart + ypart <= radiusSquared {
		return true
	}

	return false
}