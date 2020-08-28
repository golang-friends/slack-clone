# Slack Clone
![AuthService Check](https://github.com/golang-friends/slack-clone/workflows/AuthService%20Check/badge.svg)

This repository is the start of a group effort to build a Slack Clone using microservices and grpc. If you want to join please let us know.


## How to run locally

```shell
docker-compose up # --build
```

And then test

```text
❯ grpcurl --plaintext -d '{ "username": "kkweon", "email": "kkweon@gmail.com", "password": "password12345" }' 0.0.0.0:9000 protos.authservice.AuthService/Register

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6MTU5ODYwMTExMiwidXNlciI6Imtrd2VvbiJ9.T12xX4oVpmLmctiltaXaS44RPH53t91stABOIAnZjgE"
}
```

```text
❯ docker-compose logs auth

Attaching to slack-clone_auth_1
auth_1     | 2020/08/26 04:18:00 Starting Server...
auth_1     | 2020/08/26 04:18:00 Connected to Database
auth_1     | 2020/08/26 04:18:00 gRPC server hosted on 0.0.0.0:9000
```
