package utils

import (
	"strconv"
)

// For GTPL 170319,183004
func ConverttoUnixTS(rawDate, rawTime string) string {

	var layoutDate string

	if rawDate[4:] == "2019" {
		layoutDate = rawDate[4:] + "-" + string(rawDate[2:4]) + "-" + string(rawDate[0:2])
	} else {
		layoutDate = "20" + string(rawDate[4:]) + "-" + string(rawDate[2:4]) + "-" + string(rawDate[0:2])
	}

	hr24, _ := strconv.Atoi(string(rawTime[0]) + string(rawTime[1]))

	if hr24 > 12 {
		hr24 = hr24 - 12
	}

	layoutTime := strconv.Itoa(hr24) + ":" + string(rawTime[2]) + string(rawTime[3]) +
		":" + string(rawTime[4]) + string(rawTime[5]) + " "

	humanReadableDateTime := layoutDate + " " + layoutTime

	return humanReadableDateTime
}

/*
ts
speed
gsm strength
ign stat
device id
external battery
internal 11
satellite count
fuel
sos alert
tamper alert
reg no
company id , proj id
device time
lat long
live history

*/
