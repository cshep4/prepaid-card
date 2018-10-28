package request

import (
	"encoding/json"
	"errors"
	"io"
)

type AuthorisationRequest struct {
	CardNumber   *string  `bson:"cardNumber" json:"cardNumber"`
	MerchantName *string  `bson:"merchantName" json:"merchantName"`
	Amount       *float64 `bson:"amount" json:"amount"`
}

func DecodeAuthorisationRequest(body io.ReadCloser) (*AuthorisationRequest, error) {
	defer body.Close()
	var authorisationRequest AuthorisationRequest

	if err := json.NewDecoder(body).Decode(&authorisationRequest); err != nil {
		return nil, errors.New("Invalid request payload")
	}

	if authorisationRequest.CardNumber == nil || authorisationRequest.MerchantName == nil || authorisationRequest.Amount == nil {
		return nil, errors.New("Invalid request payload")
	}

	return &authorisationRequest, nil
}