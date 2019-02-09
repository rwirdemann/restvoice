# Kapitel 6: Design

## Start
```
cd kapitel08
go run main.go
```

## URIs
### POST /invoice
```
curl -X POST \
  http://localhost:8080/invoice \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018
}'
```

### GET /invoice/{invoiceId}
```
curl -X GET \
  http://localhost:8080/invoice/1 \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache'
```