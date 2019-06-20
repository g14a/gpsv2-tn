package errcheck

import "log"

// CheckError checks if an error occurs and logs it on stdout.
func CheckError(err error) {
	if err != nil {
		log.Println("Some error: ", err.Error())
	}
}
