# Kapitel 9: Sicherheit

## Setup
```
cd kapitel09
go get github.com/joho/godotenv
go get github.com/dgrijalva/jwt-go
```

## Generate JWT
```
cd identityprovider
go run main.go

Copy token
```

## Start Service
```
cd ..
go run main.go

Add copied token to authorization header of your request. See samples below.
```

## URIs
### POST /invoice
```
curl -X POST \
  http://localhost:8080/invoice \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm8gQnJ1bm5lciIsImFkbWluIjp0cnVlLCJzdWIiOiJyZXN0dm9pY2Uub3JnIn0.x2l5bn5WONnifPXShuPVNFuyeUqOzsPLCcCNAOLfBew' \
  -H 'Content-Type: application/json' \
  -d '{
    "month": 6,
    "year": 2018,
    "customerId": 1
}'
```

### POST /book/{invoiceId}
```
curl -X POST \
  http://localhost:8080/book/1 \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm8gQnJ1bm5lciIsImFkbWluIjp0cnVlLCJzdWIiOiJyZXN0dm9pY2Uub3JnIn0.x2l5bn5WONnifPXShuPVNFuyeUqOzsPLCcCNAOLfBew' \
  -H 'Content-Type: application/json' \
  -d '{
    "day": 6,
    "hours": 2,
    "projectId": 1,
    "activityId": 1,
    "description": "Refactor JWT signing"
}'```