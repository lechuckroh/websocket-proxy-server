# websocket-proxy-server

* [English](README.md)

백엔드 서버와의 웹소켓 연결을 중계하는 프록시 서버.

제공 기능
* 자바스크립트, 타입스크립트를 사용해 웹소켓 메시지 처리를 위한 미들웨어 설정
* 미들웨어를 사용해 백엔드에서 오는 메시지 무시 및 자동 응답
* 미들웨어를 사용해 클라이언트에서 보내는 메시지에 대한 Mock 응답 반환
* 웹소켓 연결별로 격리된 v8 엔진을 사용해 미들웨어 실행
* 웹소켓 연결시 esbuild를 사용해 타입스크립트 코드를 on-demand 컴파일
* node.js 와 달리 미들 웨어에서는 로컬 자원에 접근할 수 없기 때문에 보안 취약점이 없습니다.

## 문서

* [미들웨어](docs/middleware_kr.md)
* [시퀀스 다이어그램](docs/sequence_kr.md)

## 실행

```bash
# 빌드 
$ task install && task build

# 실행
$ ./ws-proxy -l :8000 -b wss://wss.example.com -f scripts/default.js
```

> `task` 커맨드를 사용하려면 [Task](https://taskfile.dev)를 설치합니다.

### 실행 옵션

실행옵션은 환경변수를 사용해 오버라이드할 수 있습니다.

### `-b`, `env:BACKEND`

웹소켓을 연결할 백엔드 서버 URL을 설정합니다.

### `-f`, `env:SCRIPT_FILE`

웹소켓이 연결된 경우 실행할 스크립트 파일을 지정합니다.

`*.js` 또는 `*.ts` 파일만 지원합니다.

### `-l`, `env:LISTEN`

리스닝할 주소와 포트를 지정합니다. 기본값은 `:8000` 입니다.

형식: `<ip address>:<port>`

### `-r`, `env:RECORD_DIR`

송수신 데이터를 저장할 디렉토리를 지정합니다.

디렉토리를 지정하지 않으면 데이터를 저장하지 않습니다.

각 연결 세션별로 별도의 디렉토리에 저장되며, 데이터 형식에 따라 `*.txt` 또는 `*.json` 파일로 저장됩니다.


## 도커 실행

```bash
$ PORT=1234 task run-docker -- -b wss://wss.example.com
```

* 기본 포트인 `8000` 대신 다른 포트를 사용하려면 `PORT` 환경변수를 사용합니다.
* 실행 옵션은 `--` 뒤에 지정합니다. 


## 알려진 문제

### MacOS M1 미지원

v8go 라이브러리에서 지원 필요.
 
관련 이슈: [v8go: add darwin/arm64 to list of architectures to build](https://github.com/rogchap/v8go/issues/54)


### 크로스 컴파일 미지원

참고문서: [v8go: Cross-compile for Linux?](https://github.com/rogchap/v8go/issues/35)
