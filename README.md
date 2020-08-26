# Slack-Clone
This repository is the start of a group effort to build a Slack Clone using microservices and grpc. If you want to join please let us know.


## How to run locally

```shell
docker-compose up # --build
```

And then test

```text
❯ grpcurl --plaintext -d '{ "username": "kkweon", "email": "kkweon@gmail.com", "password": "password12345" }' 0.0.0.0:9000 protos.authservice.AuthService/Register

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiSURcIjpcIjVmNDVlMjg2NjZkMzAyYWE5YTE1NTU0OVwiLFwiVXNlcm5hbWVcIjpcImtrd2VvblwiLFwiUGFzc3dvcmRcIjpcIiQyYSQxMCRRLlJiTzBkR3NJWjFFZTM0NzFTOS91SXU3NnNpMVRzOWJOZy42bHpCb05JMDNndzJIQ3lPLlwiLFwiRW1haWxcIjpcImtrd2VvbkBnbWFpbC5jb21cIixcIkFkbWluXCI6ZmFsc2V9In0.ODKqn3gpR3IA-FXBS_LlSWXu5t0mfupV7WiQRcgQPFE"
}
```

```text
❯ docker-compose logs auth

Attaching to slack-clone_auth_1
auth_1     | 2020/08/26 04:18:00 Starting Server...
auth_1     | 2020/08/26 04:18:00 Connected to Database
auth_1     | 2020/08/26 04:18:00 gRPC server hosted on 0.0.0.0:9000
```