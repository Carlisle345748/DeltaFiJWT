# DeltaFi Project
A simple email signup, login, logout web backend using Gin, GORM, JWT and SQLite.
- Why Gin? Gin is the one of the most popular web framework
- Why GORM? GORM is a popular and flexible ORM framework
- Why SQLite? Simple and fast local database that suitable for demo project

```
.
├── controller/
│   ├── controller.go       // controler
│   ├── controller_test.go  // controler unit test
│   └── jwt.go              // JWT authentication
├── dao/
│   ├── dao.go              // dao layer
│   └── dao_test.go         // dao unit test
├── types/
│   └── types.go
├── main.go                 // server
└── README.md
```

## Quick Start
```shell
    go get -u
    go run main.go
```

## Test
```shell
   go test ./..
```

## API
### SignUp
```
PUT http://localhost:8080/user
{
    "email": "jack@example.com",
    "firstName": "jack",
    "lastName": "mark",
    "password": "123"
}

Response
{
    "code": 0,
    "message": "success",
    "user": {
        "id": 3,
        "email": "jacker@example.com",
        "firstName": "jack",
        "lastName": "mark"
    }
}
```

### Login
```
POST http://localhost:8080/login
{
    "email": "jack@example.com",
    "password": "123"
}

Response
{
    "code": 0,
    "expire_time": "2023-02-05T22:27:21.400753-08:00",
    "message": "success",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzU2NjQ4NDEsImlkIjoyLCJvcmlnX2lhdCI6MTY3NTY2MTI0MX0.jsYtE55jcErmMw7L6vKWaWYK09dyY3JDY6gf7nyGa0E"
}
```

### Logout
```
GET http://localhost:8080/logout

Response
{
    "code": 0,
    "message": "success"
}
```

### Refresh Token
```
GET http://localhost:8080/refreshToken

Response
{
    "code": 200,
    "expire": "2023-02-05T22:30:40-08:00",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzU2NjUwNDAsImlkIjoyLCJvcmlnX2lhdCI6MTY3NTY2MTQ0MH0.TH8pJcQtvdBNxp6Bkcy91PvSo089e1Q1KlluMYUAB2o"
}
```

### Hello
```
GET http://localhost:8080/hello

Response
{
    "code": 0,
    "greeting": "Hello jack mark",
    "message": "success"
}
```

### Update User
```
POST http://localhost:8080/user
{
    "id": 1,
    "firstName": "amy",
    "lastName": "june"
}

Response
{
    "code": 0,
    "message": "success"
}
```

### Delete User
```
DELETE http://localhost:8080/user
{
    "id": 1
}

Response
{
    "code": 0,
    "message": "success"
}
```