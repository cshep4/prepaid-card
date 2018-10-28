package request

import (
	"encoding/json"
	"errors"
	"io"
)

type NewAccount struct {
	FirstName string `bson:"firstName" json:"firstName"`
	Surname   string `bson:"surname" json:"surname"`
}

func DecodeNewAccount(body io.ReadCloser) (*NewAccount, error) {
	defer body.Close()

	var newAccount NewAccount

	if err := json.NewDecoder(body).Decode(&newAccount); err != nil {
		return nil, errors.New("Invalid request payload")
	}

	if newAccount.FirstName == "" || newAccount.Surname == "" {
		return nil, errors.New("Invalid request payload")
	}

	return &newAccount, nil
}