# Kapitel 4: HTTP und JSON

## Start
```
go run main.go
```
## URIs
### GET /contacts
```
curl -X GET \
  http://localhost:8080/contacts
```

### GET /contacts/{id}
```
curl -X GET \
  http://localhost:8080/contacts/1
```

### POST /contacts
```
curl -X POST \
  http://localhost:8080/contacts \
  -H 'Content-Type: application/json' \
  -d '{
    "Firstname": "Kater",
    "Lastname": "Kalle"
}'
```

### DELETE /contacts/{id}
```
curl -X DELETE \
  http://localhost:8080/contacts/1
```

### PUT /contacts/{id}
```
curl -X PUT \
  http://localhost:8080/contacts/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "Firstname": "Jens",
    "Lastname": "Wirdemann"
}'
```