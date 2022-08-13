package app

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
)

var db = NewDatabase()

func Registration(payload []byte) (string, error) {

	var m TRegistrationRequest = TRegistrationRequest{}

	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return "", errors.New("Not valid JSON")
	}
	// return uuid.Must(uuid.New(), nil).String(), nil

	err = db.AddUser(m)
	if err != nil {

		return "", err
	}

	return "ok", nil

}

func Authenticate(payload []byte) (string, error) {
	m := TAuthenticateRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return "", errors.New("Not valid JSON")
	}
	valid, err := db.FindUser(m.Login, m.Password)
	if err != nil {
		log.Println(err)
		return "", errors.New("No user found")
	}
	if valid {

		token := uuid.Must(uuid.New(), nil).String()
		db.SetToken(m.Login, token)

		return token, nil
	}
	return "", nil
}

func CreateLobby(payload []byte) (string, error) {
	m := TCreateLobbyRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		log.Println("Something is wrong with the message!")
		return "", errors.New("Not valid JSON")
	}

	_, err = db.ValidateToken(m.Token)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid TOKEN")
	}
	log.Println(m)
	err = db.AddLobby(m)
	if err != nil {
		log.Println(err)
		return m.LobbyName, errors.New("Already created")
	}

	return m.LobbyName, nil
}

func AddToLobby(payload []byte) (string, error) {
	m := TAddToLobbyRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid JSON")
	}
	_, err = db.ValidateToken(m.Token)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid TOKEN")
	}
	err = db.SetUserToLobby(m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Cant't insert set data")
	}
	return "ok", nil
}

func PublicTokenGet(payload []byte) (string, error) {
	m := TGetPublicTokenRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid JSON")
	}
	token, err := db.findPublicToken(m.Login)
	if err != nil {
		return "", errors.New("Not TOKEN")
	}
	return token, nil
}

func SetData(payload []byte) (string, error) {
	m := TSetDataRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid JSON")
	}
	message, err := db.AddMessage(m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not lobby found")
	}
	return message, nil
}

func ValidatePermissions(payload []byte) (string, error) {
	m := TValidatePermissions{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid JSON")
	}
	lobby, err := db.ValidateLobbyPermission(m)
	if err != nil {
		log.Println(err)
		return "", errors.New("No lobby found")
	}
	return lobby, nil
}

func PublicGetLobbyId(payload []byte) (string, error) {
	m := TLobbyIdRequest{}
	err := json.Unmarshal(payload, &m)
	if err != nil {
		log.Println(err)
		return "", errors.New("Not valid JSON")
	}
	id, err := db.findLobbyId(m.LobbyName)
	if err != nil {
		return "", errors.New("Not Lobby")
	}
	return id, nil
}
