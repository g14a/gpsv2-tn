package server

import (
	"strconv"
)

func ConvertToUnixTS(rawDate, rawTime string) string {

	layoutDate := "20" + string(rawDate[4]) + string(rawDate[5]) +
				  "-" + string(rawDate[2]) + string(rawDate[3]) +
				  "-" +	string(rawDate[0]) + string(rawDate[1])

	hr24, _ := strconv.Atoi(string(rawTime[0]) + string(rawTime[1]))


	if hr24 > 12 {
		hr24 = hr24 - 12
	}

	layoutTime := strconv.Itoa(hr24) + ":" + string(rawTime[2]) + string(rawTime[3]) +
		           	":" + string(rawTime[4]) + string(rawTime[5]) + " "


	humanReadableDateTime := layoutDate + " " + layoutTime

	return humanReadableDateTime
}
