package errcheck

import "log"

func CheckError(err error) {
	if err != nil {
		log.Println("Some error: ", err.Error())
	}
}
