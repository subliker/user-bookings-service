# backendproj
> RESTful API base on Go, PostgreSQL and Docker

### Docs:
```
https://localhost:8000/docs/index.html
```

### Requirements:

![golang](https://badgen.net/static/go/1.13/green?icon=github) ![postgresql](https://badgen.net/static/postgresql/@latest/) ![docker](https://badgen.net/static/docker/@latest/purple)

### Installing:
In main directory:
```bash
docker-compose build
docker-compose up -d postgresdb
docker-compose up -d app
```

### Entities:
 - **User (example)**:
```
{
  "id": 906,
  "username": "Andrew",
  "password": "$2a$14$kv/sGmTWIlNYocbZqd88GuRsrOtKrs9bBFMM7N7HRNZ.qPxF.b.GG",
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


