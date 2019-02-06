# Kapitel 6: Design

## Start
```
cd kapitel06
go run main.go
```

## URIs
### POST /customers/{customerId}/invoices
```
curl -X POST \
  http://localhost:8080/customers/1/invoices \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018
}'
```

### POST /customers/{customerId}/invoices/{invoiceId}/bookings
```
curl -X POST \
  http://localhost:8080/customers/1/invoices/2/bookings \
  -H 'Content-Type: application/json' \
  -d '{
    "day": 16,
    "hours": 2.5,
    "projectId": 12,
    "activityId": 2,
    "description": "Bankanbindung"
}'
```

### PUT /customers/{customerId}/invoices/{invoiceId}
```
curl -X PUT \
  http://localhost:8080/customers/1/invoices/2 \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018,
    "status": "ready for aggregation"
}'
```

### GET /customers/{customerId}/invoices/{invoiceId}
```
curl -X GET \
  http://localhost:8080/customers/1/invoices/2
```


