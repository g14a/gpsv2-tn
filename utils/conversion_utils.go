// Package utils consists of conversion utilities as of now
package utils

import (
	"strconv"
	"time"
)

// ConvertTimeGTPL converts raw data and time coming from a GTPL device
// to a standard Golang time.Time object. rawData and rawTime contains 6 characters each.
func ConvertTimeGTPL(rawDate, rawTime string) time.Time {
	var istTime time.Time

	if len(rawDate) == 6 && len(rawTime) == 6 {
		hr24, _ := strconv.Atoi(string(rawTime[0]) + string(rawTime[1]))

		month := monthsMap[string(rawDate[2:4])]

		h2 := string(rawDate[0:2]) + "/" + month + "/2019" + ":" + strconv.Itoa(hr24) + ":" + string(rawTime[2:4]) + ":" + string(rawTime[4:6]) + " +0530"

		t, _ := time.Parse("02/Jan/2006:15:04:05 -0700", h2)

		istTime = t.Add(time.Hour*5 + time.Minute*30)
	}

	return istTime
}

// ConvertTimeBSTPL converts raw data and time coming from a BSTPL device
// to a standard Golang time.Time object. rawData and rawTime contains 6 characters each.
func ConvertTimeBSTPL(rawDate, rawTime string) time.Time {
	var istTime time.Time

	if len(rawDate) == 6 && len(rawTime) == 6 {
		hr24, _ := strconv.Atoi(string(rawTime[0]) + string(rawTime[1]))

		month := monthsMap[string(rawDate[2:4])]

		h2 := string(rawDate[0:2]) + "/" + month + "/2019" + ":" + strconv.Itoa(hr24) + ":" + string(rawTime[2:4]) + ":" + string(rawTime[4:6]) + " +0530"

		t, _ := time.Parse("02/Jan/2006:15:04:05 -0700", h2)

		istTime = t.Add(time.Hour*5 + time.Minute*30)
	}

	return istTime
}

// ConvertTimeAIS140 converts raw data and time coming from a AIS140 device
// to a standard Golang time.Time object. rawData and rawTime contains 8 characters and 6 characters respectively.
func ConvertTimeAIS140(rawDate, rawTime string) time.Time {
	var istTime time.Time

	if len(rawDate) == 8 && len(rawTime) == 6 {
		hr24, _ := strconv.Atoi(string(rawTime[0]) + string(rawTime[1]))

		month := monthsMap[string(rawDate[2:4])]

		h2 := string(rawDate[0:2]) + "/" + month + "/2019" + ":" + strconv.Itoa(hr24) + ":" + string(rawTime[2:4]) + ":" + string(rawTime[4:6]) + " +0530"

		t, _ := time.Parse("02/Jan/2006:15:04:05 -0700", h2)

		istTime = t.Add(time.Hour*5 + time.Minute*30)
	}
	return istTime
}

// Add5Hrs adds 5 hrs and 30 minutes to deviceTime and returns the added time.
func Add5Hrs(deviceTime time.Time) time.Time {
	istTime := deviceTime.Add(time.Hour*5 + time.Minute*30)

	return istTime
}


// GTPLCheckLiveHistory checks if the GTPL device time is a live packet
// by having a 60 second boundary.
func GTPLCheckLiveHistory(istTime time.Time) bool {
	difference := time.Now().Sub(istTime)

	if difference.Seconds() > 60 {
		return false
	}

	return true
}

// monthsMap is a map which maps indexes of the months of a year
// to its English representation
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
