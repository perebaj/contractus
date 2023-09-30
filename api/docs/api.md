# API

Basic usage examples of the contractus API

**@Obs**: To run locally, the right host to put in the <> is the host that will retrieve after running `make ip`

# Endpoints

## `/upload`

Request:

```bash
    curl -i -X 'POST' \
        'http://<>:8080/upload' \
        -H 'accept: application/json' \
        -H 'Content-Type: multipart/form-data' \
        -H 'Cookie: jwt=<>' \
        -F 'file=@sales.txt;type=text/plain' 
```

Response:

```JSON
{
  "msg": "file uploaded successfully"
}
```

## `/balance/affiliate?name=<>`

This route is responsible to retrieve the balance for a given affiliate name


### Parameters

**name:** Name of the seller 

Obs: If the affilate's name is "JONATHAN SILVA," it should be input as `?name=JONATHAN%20SILVA` to ensure proper search functionality.

Request: 

```bash
curl -i -X 'GET' \
  'http://<>:8080/balance/affiliate?name=<>' \
  -H 'accept: application/json' \
  -H 'Cookie: jwt=<>'
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
curl -i -X 'GET' \
  'http://<>:8080/balance/producer?name=<>' \
  -H 'accept: application/json' \
  -H 'Cookie: jwt=<>'
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
  'http://<>:8080/transactions' \
  -H 'accept: application/json' \
  -H 'Cookie: jwt=<>'
```

Response:

```JSON
{
  "transactions": [
    {
      "email": "jj@example.com",
      "type": "string",
      "date": "2023-09-27T17:34:41.455Z",
      "product_description": "string",
      "product_price": "string",
      "seller_name": "string",
      "seller_type": "affiliate",
      "action": "venda produtor"
    }
  ],
  "total": 0
}
```
