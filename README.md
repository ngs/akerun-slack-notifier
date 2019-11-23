# Akerun Slack Notifier

## Env

| Name                   | Decription |
| ---------------------- | ---------- |
| `AKERUN_CLIENT_ID`     |            |
| `AKERUN_CLIENT_SECRET` |            |
| `AKERUN_CALLBACK_URL`  |            |
| `S3_BUCKET`            |            |

## Dev

```sh
npm i -g serverless
go get -u github.com/golang/dep/cmd/dep
dep ensure
make
```
