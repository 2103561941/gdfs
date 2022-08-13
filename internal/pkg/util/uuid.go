package util

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// get uuid may error, I choose to try more times.
// If it still doesn't solve the problem, return an error.
func GetUUID() (id string, err error) {
	for i := 0; i < 5; i++ {
		id, err = gonanoid.New()
		if err == nil {
			return
		}
	}
	// TODO: print relative logs to record this problem

	
	return "", fmt.Errorf("couldn't get uuid: %w", err)
}
