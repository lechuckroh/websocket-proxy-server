# websocket-proxy-server

## Run

```bash
# build 
$ task install && task build

# run
$ ./ws-proxy -l :80 -b wss://wss.example.com -f scripts/default.js
```

> Install [Task](https://taskfile.dev) to use `task` command.

### Options

### `-b`, `env:BACKEND`

Set backend websocket server URL.

### `-f`, `env:SCRIPT_FILE`

Set script file for middleware.

### `-l`, `env:LISTEN`

Set listening address. default value is `:80`.

Format: `<ip address>:<port>`

### `-r`, `env:RECORD_DIR`

Directory to store traffic records in JSON format.
feat