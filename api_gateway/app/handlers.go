package app

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
)

func Registration(payload []byte) (string, error) {
	m := RegistrationRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return "", errors.New("Not valid JSON")
	}
	return uuid.Must(uuid.New(), nil).String(), nil
	// RegistrationRequest(m)

}

func KeepData(payload []byte) (string, error) {
	m := KeepDataRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return "", errors.New("Not valid JSON")
	}
	return m.Message, nil
	// RegistrationRequest(m)

}
