# Kapitel 5: Restvoice

## Start
```
go run main.go
```
## URIs
### GET /customers
```
curl -X GET \
  http://localhost:8080/customers
```

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

### GET /customers/{customerId}/projects
```
curl -X GET \
  http://localhost:8080/customers/1/projects
```

### GET /activities
```
curl -X GET \
  http://localhost:8080/activities
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
    "activityId": 2
}'
```

### DELETE /customers/{customerId}/invoices/{invoiceId}/bookings/{bookingId}
```
curl -X DELETE \
  http://localhost:8080/customers/1/invoices/2/bookings/1
```

### PUT /customers/1/invoices/1
```
curl -X PUT \
  http://localhost:8080/customers/1/invoices/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018
    "status": "ready for aggregation"
}'
```

### GET 
/customers/1/invoices/1234

```
curl -X GET \
  http://localhost:8080/customers/1/invoices/2 \
  -H 'Accept: application/pdf'
```


