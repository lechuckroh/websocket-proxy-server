package session

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lechuckroh/websocket-proxy-server/polyfill/proxy"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"net/url"
	"rogchap.com/v8go"
	"time"
)

var (
	DefaultDialer   = websocket.DefaultDialer
	DefaultUpgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Session interface {
	Start()
}

type Opts struct {
	TargetURL  *url.URL
	Script     string
	RecordDir  string
	RespWriter http.ResponseWriter
	Request    *http.Request
}

func NewSession(opts *Opts) (Session, error) {
	if opts.TargetURL == nil {
		return nil, errors.New("targetURL is not set")
	}

	req := opts.Request
	backendURL := *opts.TargetURL
	backendURL.Fragment = req.URL.Fragment
	backendURL.Path = req.URL.Path
	backendURL.RawQuery = req.URL.RawQuery

	return &sessionImpl{
		dialer:     DefaultDialer,
		upgrader:   DefaultUpgrader,
		recordDir:  opts.RecordDir,
		respWriter: opts.RespWriter,
		request:    req,
		backendURL: &backendURL,
		script:     opts.Script,
		sessionID:  uuid.NewV4().String(),
	}, nil
}

type sessionImpl struct {
	dialer     *websocket.Dialer
	upgrader   *websocket.Upgrader
	backendURL *url.URL
	script     string
	sessionID  string
	recordDir  string
	respWriter http.ResponseWriter
	request    *http.Request
	iso        *v8go.Isolate
	v8Context  *v8go.Context
	proxy      proxy.Proxy
}

func (s *sessionImpl) Start() {
	connBackend, connClient, err := s.connectBackend()
	if err != nil {
		s.logErrorf("failed to connect backend: %v", err)
		return
	}

	defer func() {
		_ = connBackend.Close()
		_ = connClient.Close()
		s.log("session closed")

		// onDestroy
		if err := s.proxy.ExecuteOnDestroy(); err != nil {
			s.logErrorf("failed to run onDestroy() function: %v", err)
		}
	}()

	// initialize v8go
	if ctx, err := initV8(); err != nil {
		s.logErrorf("failed to initialize v8: %v", err)
		return
	} else {
		// inject proxy object
		prx, err := proxy.InjectTo(ctx)
		if err != nil {
			s.logErrorf("failed to inject proxy object: %v", err)
			return
		}
		s.v8Context = ctx
		s.iso, _ = ctx.Isolate()
		s.proxy = prx
	}

	// run script
	if err := s.runScript(); err != nil {
		s.logErrorf("failed to run script: %v", err)
		return
	}

	// onInit
	if err := s.proxy.ExecuteOnInit(); err != nil {
		s.logErrorf("failed to run onInit() function: %v", err)
	}

	// messageWriter
	receiveMsgWriter := NewMessageWriter(s.recordDir, s.sessionID, s.getFilenameGenerator("recv"))
	sendMsgWriter := NewMessageWriter(s.recordDir, s.sessionID, s.getFilenameGenerator("sent"))
	if err := receiveMsgWriter.Init(); err != nil {
		s.logErrorf("failed to create receivedMsgWriter: %v", err)
		return
	}
	if err := sendMsgWriter.Init(); err != nil {
		s.logErrorf("failed to create sentMsgWriter: %v", err)
		return
	}

	// forward received message
	errClientCh := make(chan error, 1)
	go s.forwardMessage(connBackend, connClient, errClientCh,
		s.proxy.ExecuteReceivedMessageMiddlewares, s.proxy.ExecuteResponseToBackendMessageMiddlewares,
		receiveMsgWriter, sendMsgWriter)
	// forward sent message
	errBackendCh := make(chan error, 1)
	go s.forwardMessage(connClient, connBackend, errBackendCh,
		s.proxy.ExecuteSentMessageMiddlewares, s.proxy.ExecuteResponseToClientMessageMiddlewares,
		sendMsgWriter, receiveMsgWriter)

	// wait for errors
	var errMsg string
	select {
	case err = <-errClientCh:
		if s.isCloseError(err) {
			errMsg = "terminated from server"
		} else {
			errMsg = err.Error()
		}
	case err = <-errBackendCh:
		if s.isCloseError(err) {
			errMsg = "terminated from client"
		} else {
			errMsg = err.Error()
		}
	}

	s.logErrorf("connection closed: %v", errMsg)
}

func (s *sessionImpl) getFilenameGenerator(typeName string) FilenameGenerator {
	idx := 0
	return func(ext string) string {
		idx++
		return fmt.Sprintf("%d_%s_%05d.%s", time.Now().Unix(), typeName, idx, ext)
	}
}

func (s *sessionImpl) isCloseError(err error) bool {
	var e *websocket.CloseError
	return errors.As(err, &e)
}

func (s *sessionImpl) runScript() error {
	_, err := s.v8Context.RunScript(s.script, "")
	return err
}

// connectBackend connect to backend server
func (s *sessionImpl) connectBackend() (*websocket.Conn, *websocket.Conn, error) {
	s.logInfof("connecting to backend: %s", s.backendURL)

	dialer := s.dialer
	if dialer == nil {
		dialer = DefaultDialer
	}

	header := getForwardingHeader(s.request, nil)

	// Connect to the backend URL
	connBackend, resp, err := dialer.Dial(s.backendURL.String(), header)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial backend URL: %s", s.backendURL)
	}

	s.logInfof("connected to backend: %s", s.backendURL)

	// upgrade connection with client
	upgrader := s.upgrader
	if upgrader == nil {
		upgrader = DefaultUpgrader
	}

	respHeader := http.Header{}
	copyHeaders(resp.Header, &respHeader, []string{"Set-Cookie", "Sec-Websocket-Protocol"})
	connClient, err := upgrader.Upgrade(s.respWriter, s.request, respHeader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to upgrade connection: %v", err)
	}
	s.log("connection upgraded to websocket")

	return connBackend, connClient, nil
}

// forwardMessage forwards messages.
// Send error to errCh on connection closed.
func (s *sessionImpl) forwardMessage(
	fromConn, toConn *websocket.Conn,
	errCh chan error,
	executeMiddlewaresFn proxy.ExecuteMiddlewaresFn,
	executeResponseMiddlewaresFn proxy.ExecuteMiddlewaresFn,
	receiveMessageWriter MessageWriter,
	responseMessageWriter MessageWriter,
) {
	// create responseFunc
	var resFn *v8go.Value
	if resFnTpl, err := v8go.NewFunctionTemplate(s.iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) != 1 {
			s.logErrorf("responseFunc requires 1 argument")
			return nil
		}
		arg := args[0]

		// call middleware
		result, err := executeResponseMiddlewaresFn(arg)
		if err != nil {
			s.logErrorf("failed to execute response middleware: %v", err)
			return nil
		}

		// do not forward if middleware result is null or undefined
		if result == nil || result.IsNullOrUndefined() {
			return nil
		}

		if err := responseMessageWriter.WriteValue(result); err != nil {
			s.logErrorf("failed to record message: %v", err)
			errCh <- err
			return nil
		}

		if err := s.sendMessage(fromConn, result); err != nil {
			s.logErrorf("responseFunc failed: %v", err)
			return nil
		}

		return nil
	}); err != nil {
		s.logErrorf("failed to register responseFunc: %v", err)
		return
	} else {
		resFn = resFnTpl.GetFunction(s.v8Context).Value
	}

	for {
		// read message from source
		msgType, readMsgBytes, err := fromConn.ReadMessage()
		if err != nil {
			errCh <- err
			return
		}

		if err := receiveMessageWriter.Write(readMsgBytes); err != nil {
			s.logErrorf("failed to record message: %v", err)
			errCh <- err
			return
		}

		switch msgType {
		case websocket.TextMessage:
			if ok := s.forwardTextMessage(toConn, readMsgBytes, executeMiddlewaresFn, resFn); !ok {
				return
			}
		default:
			s.logErrorf("unhandled messageType: %d", msgType)
		}
	}
}

// forwardTextMessage forwards message to 'toConn'
// Returns true if OK, false otherwise.
func (s *sessionImpl) forwardTextMessage(
	toConn *websocket.Conn,
	data []byte,
	executeMiddlewaresFn proxy.ExecuteMiddlewaresFn,
	resFn *v8go.Value,
) bool {
	dataValue, _ := v8go.NewValue(s.iso, string(data))

	// call middleware
	// arguments: (message, responseFunc)
	result, err := executeMiddlewaresFn(dataValue, resFn)
	if err != nil {
		s.logErrorf("failed to execute middleware: %v", err)
		return false
	}

	// do not forward if middleware result is null or undefined
	if result == nil || result.IsNullOrUndefined() {
		return true
	}

	// forward message to toConn
	err = s.sendMessage(toConn, result)
	if err != nil {
		s.logErrorf("failed to forward message: %v", err)
		return false
	}

	return true
}

func (s *sessionImpl) log(message string) {
	log.Printf("%s >> %s", s.sessionID, message)
}

func (s *sessionImpl) logInfof(format string, args ...interface{}) {
	s.log(fmt.Sprintf(format, args...))
}

func (s *sessionImpl) logErrorf(format string, args ...interface{}) {
	s.log(fmt.Sprintf(format, args...))
}

func (s *sessionImpl) sendMessage(conn *websocket.Conn, message *v8go.Value) error {
	if message.IsString() {
		return conn.WriteMessage(websocket.TextMessage, []byte(message.String()))
	} else {
		return conn.WriteJSON(message)
	}
}
