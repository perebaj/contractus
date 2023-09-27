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

Request: 

```bash
curl -X 'GET' \
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

Request: 

```bash
curl -X 'GET' \
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
