package errorcheck

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

// CheckError checks if an error occurs and logs it on listener.log
func CheckError(err error) {
	if err != nil {
		logger.Println("error: ", err.Error())
	}
}

func init() {
	logPath := "/root/gpsv2-tn/listener.log"

	file, err := os.Create(logPath)

	CheckError(err)

	logger = log.New(file, "", log.LstdFlags | log.Lshortfile)
}
