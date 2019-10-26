package errorhandler

import (
	"log"
)

// HandleError a general handler for errors in go
func HandleError(location string, err error) {
	if err != nil {
		log.Fatal(location + err.Error())
	}
}
