package request

import (
	"encoding/json"
	"errors"
	"io"
)

type LoadMoney struct {
	CardNumber *string  `bson:"cardNumber" json:"cardNumber"`
	Amount     *float64 `bson:"amount" json:"amount"`
}

func DecodeLoadMoney(body io.ReadCloser) (*LoadMoney, error) {
	defer body.Close()

	var loadMoney LoadMoney

	if err := json.NewDecoder(body).Decode(&loadMoney); err != nil {
		return nil, errors.New("Invalid request payload")
	}

	if loadMoney.CardNumber == nil || loadMoney.Amount == nil {
		return nil, errors.New("Invalid request payload")
	}

	if *loadMoney.Amount < 0 {
		return nil, errors.New("Cannot load a negative amount onto card")
	}

	return &loadMoney, nil
}