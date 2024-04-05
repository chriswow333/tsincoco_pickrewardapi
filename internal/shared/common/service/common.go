package service

import (
	"strings"

	"github.com/google/uuid"
)

func GenUUID() string {

	id := uuid.New().String()
	return strings.Replace(id, "-", "", -1)
}
