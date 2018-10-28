package request

import (
	"encoding/json"
	"errors"
	"io"
)

type TransactionRequest struct {
	TransactionId string  `bson:"transactionId" json:"transactionId"`
	Amount        float64 `bson:"amount" json:"amount"`
}

func DecodeTransactionRequest(body io.ReadCloser) (*TransactionRequest, error) {
	defer body.Close()
	var transactionRequest TransactionRequest

	if err := json.NewDecoder(body).Decode(&transactionRequest); err != nil {
		return nil, errors.New("Invalid request payload")
	}

	if transactionRequest.TransactionId == "" || transactionRequest.Amount <= 0 {
		return nil, errors.New("Invalid request payload")
	}

	return &transactionRequest, nil
}
