package util

import gonanoid "github.com/matoous/go-nanoid/v2"

func GetUUID() (id string, err error) {
	id, err = gonanoid.New()
	return
}
