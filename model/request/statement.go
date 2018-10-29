package request

import (
	"encoding/json"
	"errors"
	"io"
	. "time"
)

type StatementRequest struct {
	CardNumber *string `bson:"cardNumber" json:"cardNumber"`
	FromDate   *Time   `bson:"fromDate" json:"fromDate"`
	ToDate     *Time   `bson:"toDate" json:"toDate"`
}

func DecodeStatementRequest(body io.ReadCloser) (*StatementRequest, error) {
	defer body.Close()
	var statementRequest StatementRequest

	if err := json.NewDecoder(body).Decode(&statementRequest); err != nil {
		return nil, errors.New("Invalid request payload")
	}

	if statementRequest.CardNumber == nil {
		return nil, errors.New("Invalid request payload")
	}

	return &statementRequest, nil
}
