package utils

import (
	"strconv"
	"time"
)

// For GTPL 170319,183004
func ConvertToUnixTS(rawDate, rawTime string) time.Time {

	amPm := "am"

	hr24, _ := strconv.Atoi(string(rawTime[0]) + string(rawTime[1]))

	if hr24 > 12 {
		hr24 = hr24 - 12
		amPm = "pm"
	}

	month := monthsMap[string(rawDate[2:4])]

	humanReadableDateTime := month + " " + string(rawDate[0:2]) + ", 2019 at " +
		strconv.Itoa(hr24) + ":" + string(rawTime[2:4]) + amPm + " (IST)"

	const longForm = "Jan 02, 2006 at 3:04pm (MST)"

	t, _ := time.Parse(longForm, humanReadableDateTime)

	return t

}

var monthsMap = map[string]string{
	"01": "Jan",
	"02": "Feb",
	"03": "Mar",
	"04": "Apr",
	"05": "May",
	"06": "Jun",
	"07": "Jul",
	"08": "Aug",
	"09": "Sep",
	"10": "Oct",
	"11": "Nov",
	"12": "Dec",
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
