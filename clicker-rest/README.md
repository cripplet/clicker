clicker-rest
====
REST API for Cookie Clicker

```bash
go run server.go \
    --port=8080 \
    --credentials=/path/to/credentials.json \
    --project=firebase-project-id \
    --environment=(dev|prod) &

wget --post-data='' localhost:8080/game
wget --post-data='' localhost:8080/game/generated-game-id/cookie/mine
```
