# Kapitel 6: Design

## Start
```
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

### POST /customers/1/invoices/2/bookings
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

### PUT /customers/1/invoices/2
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


