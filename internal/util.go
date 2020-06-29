package internal

import (
	"log"
)

func Report(err error) {
	log.Printf("An error occurred %s\n", err)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
