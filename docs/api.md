# API

Basic usage examples of the contractus API

# Endpoints

## `/upload`

Request:

```
    curl -i -X 'POST' \
        'http://localhost:8080/upload' \
        -H 'accept: application/json' \
        -H 'Content-Type: multipart/form-data' \
        -F 'file=@sales.txt;type=text/plain'
```
