# A Simple Prepaid Card

This is written in Golang with Gorilla Mux used for routing. A MongoDB database is used to store data. It is hosted on Heroku at https://simple-prepaid-card.herokuapp.com/.

Please note that as this is deployed on Heroku's free tier, the web server will shut down after 30 minutes of inactivity. Therefore the initial request may take slightly longer whilst the server starts up again.

# REST API

The REST API to the prepaid card app is described below.

## Create an account

### Request

`POST /account/create`

    curl -d '{"firstName": "Chris", "surname": "Shepherd"}' -H "Content-Type: application/json" -X POST https://simple-prepaid-card.herokuapp.com/account/create

###### Request Body

    {
      "firstName": "Chris",
      "surname": "Shepherd"
    }

### Response

    HTTP/1.1 201 Created
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 201 Created
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 163

    {
      "id": "5bd764707e84828cfaf21160",
      "cardNumber": "154084260800001",
      "firstName": "Chris",
      "surname": "Shepherd",
      "totalBalance": 0,
      "availableBalance": 0,
      "transactions": null
    }

## Load money

### Request

`PUT /account/load`

    curl -d '{ "cardNumber": "154084260800001", "amount": 150 }' -H "Content-Type: application/json" -X PUT https://simple-prepaid-card.herokuapp.com/account/load

###### Request Body

    {
      "cardNumber": "154084260800001",
      "amount": 150
    }

### Response

    HTTP/1.1 200 OK
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 200 OK
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 20

    {
      "result": "success"
    }

## Get balance

### Request

`GET /account/balance/{cardNumber}`

    curl -i -H 'Accept: application/json' https://simple-prepaid-card.herokuapp.com/account/balance/154084260800001

### Response

    HTTP/1.1 200 OK
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 200 OK
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 43

    {
      "totalBalance": 150,
      "availableBalance": 150
    }

## Get account statement

### Request

`GET /account/statement/{cardNumber}`

'toDate' and 'fromDate' are optional query parameters. If these are not specified then all transactions are returned.

    curl -H "Content-Type: application/json" https://simple-prepaid-card.herokuapp.com/account/statement/154084260800001

    curl -H "Content-Type: application/json" https://simple-prepaid-card.herokuapp.com/account/statement/154084260800001?fromDate=2018-10-29&toDate=2018-10-30

### Response

    HTTP/1.1 200 OK
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 200 OK
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 377

    {
      "accountDetails": {
        "firstName": "Chris",
        "surname": "Shepherd",
        "cardNumber": "154084260800001"
      },
      "fromDate": "2018-10-27T00:00:00Z",
      "toDate": "2018-10-29T00:00:00Z",
      "balance": {
        "totalBalance": 150,
        "availableBalance": 148.25
      },
      "transactions": [
        {
            "id": "5bd76ba19382d300047ccc0c",
            "authorisedAmount": 1.75,
            "capturedAmount": 0,
            "merchantName": "Coffee Shop",
            "timestamp": "2018-10-29T20:20:49.954Z"
        }
      ]
    }

## Create transaction authorisation request

### Request

`POST /transaction/authorisation`

    curl -d '{ "cardNumber": "154084260800001", "amount": 1.75, "merchantName": "Coffee Shop" }' -H "Content-Type: application/json" -X POST https://simple-prepaid-card.herokuapp.com/transaction/authorisation

###### Request Body

    {
      "cardNumber": "154084260800001",
      "amount": 1.75,
      "merchantName": "Coffee Shop"
    }

### Response

    HTTP/1.1 201 Created
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 201 Created
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 149

    {
      "id": "5bd76ba19382d300047ccc0c",
      "authorisedAmount": 175,
      "capturedAmount": 0,
      "merchantName": "Coffee Shop",
      "timestamp": "2018-10-29T20:20:49.954413649Z"
    }

## Capture a transaction

### Request

`PUT /transaction/capture`

    curl -d '{ "transactionId": "5bd7428f7e848277e39a402a", "amount": 1.99 }' -H "Content-Type: application/json" -X PUT https://simple-prepaid-card.herokuapp.com/transaction/capture

###### Request Body

    {
      "transactionId": "5bd7428f7e848277e39a402a",
      "amount": 1.99
    }

### Response

    HTTP/1.1 200 OK
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 200 OK
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 142

    {
      "id": "5bd7428f7e848277e39a402a",
      "authorisedAmount": 299,
      "capturedAmount": 199,
      "merchantName": "Go Fresh",
      "timestamp": "2018-10-29T17:25:35.337Z"
    }

## Reverse authorisation

### Request

`PUT /transaction/reverse`

    curl -d '{ "transactionId": "5bd7428f7e848277e39a402a", "amount": 0.89 }' -H "Content-Type: application/json" -X PUT https://simple-prepaid-card.herokuapp.com/transaction/reverse

###### Request Body

    {
      "transactionId": "5bd7428f7e848277e39a402a",
      "amount": 0.89
    }

### Response

    HTTP/1.1 200 OK
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 200 OK
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 142

    {
      "id": "5bd7428f7e848277e39a402a",
      "authorisedAmount": 210,
      "capturedAmount": 199,
      "merchantName": "Go Fresh",
      "timestamp": "2018-10-29T17:25:35.337Z"
    }

## Refund transaction

### Request

`POST /transaction/refund`

    curl -d '{ "transactionId": "5bd7428f7e848277e39a402a", "amount": 1.94 }' -H "Content-Type: application/json" -X POST https://simple-prepaid-card.herokuapp.com/transaction/refund

###### Request Body

    {
      "transactionId": "5bd7428f7e848277e39a402a",
      "amount": 1.94
    }

### Response

    HTTP/1.1 201 Created
    Date: Mon, 29 Oct 2018 19:50:08 GMT
    Status: 201 Created
    Access-Control-Allow-Origin: *
    Content-Type: application/json
    Content-Length: 162

    {
      "id": "5bd772239382d300047ccc0d",
      "authorisedAmount": -100,
      "capturedAmount": -100,
      "merchantName": "REFUND - Coffee Shop",
      "timestamp": "2018-10-29T20:48:35.601504878Z"
    }
