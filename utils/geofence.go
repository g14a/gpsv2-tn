package utils

func GeoFence(baselat, baselong, devicelat, devicelong, radius float64) bool {

	radiusInKM := radius/1000
	// if vehicle is inside return true
	if Distance(baselat, baselong, devicelat, devicelong) < radiusInKM {
		return true
	}

	return false
}