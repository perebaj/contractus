# API

Basic usage examples of the contractus API

# Endpoints

## `/upload`

Request:

```bash
    curl -i -X 'POST' \
        'http://localhost:8080/upload' \
        -H 'accept: application/json' \
        -H 'Content-Type: multipart/form-data' \
        -F 'file=@sales.txt;type=text/plain'
```

## `/balance/affiliate?name=<>`

This route is responsible to retrieve the balance for a given affiliate name


### Parameters

**name:** Name of the seller 

Obs: If the affilate's name is "JONATHAN SILVA," it should be input as `?name=JONATHAN%20SILVA` to ensure proper search functionality.

Request: 

```bash
curl -X -i 'GET' \
  'http://localhost:8080/balance/affiliate?name=<>' \
  -H 'accept: application/json'
```

Response:

```JSON
{
  "balance": 0,
  "seller_name": "string"
}
```

## `/balance/producer?name=<>`

This route is responsible to retrieve the balance for a given producer name

### Parameters

**name:** Name of the seller 

Obs: If the producer's name is "JONATHAN SILVA," it should be input as `?name=JONATHAN%20SILVA` to ensure proper search functionality.


Request: 

```bash
curl -X -i 'GET' \
  'http://localhost:8080/balance/producer?name=<>' \
  -H 'accept: application/json'
```

Response:

```JSON
{
  "balance": 0,
  "seller_name": "string"
}
```

## `/transactions`

List all transactions

Request: 

```bash
curl -i -X 'GET' \
  'http://localhost:8080/transactions' \
  -H 'accept: application/json'
```

Response:

```JSON
{
  "transactions": [
    {
      "type": "string",
      "date": "2023-09-27T17:27:14.189Z",
      "productDescription": "string",
      "productPrice": "string",
      "sellerName": "string",
      "sellerType": "affiliate",
      "action": "venda produtor"
    }
  ],
  "total": 0
}
```
