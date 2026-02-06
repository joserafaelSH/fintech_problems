package utils

import "github.com/google/uuid"

func GenerateID() string {
	v7, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return v7.String()
}
