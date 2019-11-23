# Akerun Slack Notifier

## Env

| Name                   | Decription |
| ---------------------- | ---------- |
| `AKERUN_CLIENT_ID`     |            |
| `AKERUN_CLIENT_SECRET` |            |
| `AKERUN_REDIRECT_URL`  |            |
| `S3_BUCKET`            |            |
| `SLACK_WEBHOOK_URL`    |            |

## Dev

```sh
npm i -g serverless
go get -u github.com/golang/dep/cmd/dep
make deploy
```
