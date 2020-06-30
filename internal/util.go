package internal

import (
	"log"
)

// Report reports errors
func Report(err error) {
	log.Printf("An error occurred %s\n", err)
}

// StringInSlice checks if a string is present in a slice of strings
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
