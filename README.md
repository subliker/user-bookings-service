# backendproj
> RESTful API based on Go, PostgreSQL and Docker

### Docs (Sweagger:
```
https://localhost:8000/docs/index.html
```

### Requirements:
#### With Docker:
 ![docker](https://badgen.net/static/docker/@latest/purple)<br/>
 You can install Docker <a href="https://docs.docker.com/engine/install/">there</a>

#### Without Docker:
 ![golang](https://badgen.net/static/go/1.13/green?icon=github) ![postgresql](https://badgen.net/static/postgresql/@latest/)<br/>
 You can install Golang <a href="https://go.dev/doc/install">there</a><br/>
 You can install PostgreSQL <a href="https://www.postgresql.org/download/">there</a>

### Installing:
1. Clone repository 
2. In main directory:<br/>
   With Docker:
    for Windows users:
      ```bash
      docker-compose build
      docker-compose up -d postgresdb
      docker-compose up -d app
      ```
    for Linux users:
      ```bash
      sudo docker compose build
      sudo docker compose up -d postgresdb
      sudo docker compose up -d app
      ```
   Without Docker:
    ```
    go run main.go
    ```
    and <br/>
    Set PostgreSQL in pgAdmin (see in env file (.env))

### Entities:
 - **User (example)**:
```
{
  "id": 906,
  "username": "Andrew",
  "password": "$2a$14$kv/sGmTWIlNYocbZqd88GuRsrOtKrs9bBFMM7N7HRNZ.qPxF.b.GG", //bcrypt hash
  "created_at": "2023-09-24T17:13:42Z",
  "updated_at": "2023-09-27T11:10:23Z"
}
```
 - **Booking (example)**:
```
{
  "id": 1021,
  "user_id": 906,
  "end_time": "2023-10-01T14:30:00Z",
  "start_time": "2023-10-01T12:00:00Z",
  "comment": "I may be a little late"
}
```

### Requests
- /user/{id} [get]
  <br/>Get User by id
- /user [post]
  <br/>Create User from postForm: username, password
- /user/{user_id} [delete]
  <br/>Delete User and its bookings by user_id
- /user/{id} [put]
  <br/>Update User data (optional: username, password) by id (set new timestamp in update_at)

- /booking [get]
  <br/>Get all bookings ordered by id
- /booking/{id} [get]
  <br/>Get Booking by id (optional: set limit, page(required limit), offset(required limit) in params)
- /booking [post]
  <br/>Create User from postForm: user_id, start_time, end_time, comment(optional)
- /booking/{id} [delete]
  <br/>Delete Booking by id
- /booking/{id} [put]
  <br/>Update Booking data (optional: start_time, end_time, comments) by id
