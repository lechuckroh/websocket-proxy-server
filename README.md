# websocket-proxy-server

* [한국어](README_kr.md)

A proxy server that relays websocket connections from client to backend server.

Features:
* Write custom middleware for websocket message handling using JavaScript and TypeScript.
* On-demand TypeScript compilation using esbuild when websocket connection established.
* Can write middlewares to ignore/auto-response to messages from backend server.
* Can write middlewares to return a mock response to client.
* Run middleware in isolated v8 engine per websocket connection.
* No vulnerabilities. Unlike node.js, middleware cannot access local resources.

## Documents

* [Middleware](docs/middleware.md)
* [Sequence Diagram](docs/sequence.md)


## Run

```bash
# build 
$ task install && task build

# run
$ ./ws-proxy -l :8000 -b wss://wss.example.com -f scripts/default.js
```

> Install [Task](https://taskfile.dev) to use `task` command.

### Options

You can override options using environment variables.

### `-b`, `env:BACKEND`

Set backend websocket server URL.

### `-f`, `env:SCRIPT_FILE`

Set script file for middleware.

Supports `*.js` and `*.ts` files.

### `-l`, `env:LISTEN`

Set listening address. default value is `:8000`.

Format: `<ip address>:<port>`

### `-r`, `env:RECORD_DIR`

Directory to store traffic records in text format.

If not specified, no records are saved.

It is saved in a separate directory for each websocket connection.
Data is saved as a `*.txt` or `*.json` file depending on message format.

## Run using docker

```bash
$ PORT=1234 task run-docker -- -b wss://wss.example.com
```

* Set `PORT` environment variable to override default listening port(`8000`).
* Append custom CLI options after `--`.


## Known problems

### MacOS M1 is not supported

v8go does not support yet.

See [v8go: add darwin/arm64 to list of architectures to build](https://github.com/rogchap/v8go/issues/54)

### Cross-compile is not supported

See [v8go: Cross-compile for Linux?](https://github.com/rogchap/v8go/issues/35)
