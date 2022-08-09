package handlers

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/n0byk/keeper/api_gateway/lib/distribution"

	"github.com/n0byk/keeper/engine"
)

func Registration(payload []byte) (string, error) {
	m := engine.RegistrationRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return "", errors.New("Not valid JSON")
	}

	distribution.RegistrationRequest(m)

}
