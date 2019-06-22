<<<<<<< HEAD:errorcheck/error_check.go
// Package errorcheck just checks if an error occurs and logs it on stdout
=======
>>>>>>> dev:errcheck/error_check.go
package errorcheck

import "log"

// CheckError checks if an error occurs and logs it on stdout.
func CheckError(err error) {
	if err != nil {
		log.Println("Some error: ", err.Error())
	}
}
