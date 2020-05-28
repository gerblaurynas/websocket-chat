# websocket-chat
Websocket chat server in go

All connected clients receive text messages in json format:

```json
{
    "text":"hello world",
    "timestamp":"2020-05-28T15:32:17.510029905Z",
    "username":"foobar"
}
```

## Setup

1. (Optional) Adjust server port if necessary in .env `cp .env.dist .env` and edit SERVER_PORT. `8080` is the default value.
1. Run server: `docker-compose up`.
1. Connect client `ws://localhost:8080/websocket?username=foobar`.
1. Send json messages like `{"text":"hello world"}`
