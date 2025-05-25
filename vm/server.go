package vm

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	lua "github.com/yuin/gopher-lua"
)

type ServerModule struct {
	vm       *SolVM
	servers  map[string]*http.Server
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

func NewServerModule(vm *SolVM) *ServerModule {
	return &ServerModule{
		vm:      vm,
		servers: make(map[string]*http.Server),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (sm *ServerModule) Register() {
	sm.vm.RegisterFunction("create_server", sm.createServer)
	sm.vm.RegisterFunction("start_server", sm.startServer)
	sm.vm.RegisterFunction("stop_server", sm.stopServer)
	sm.vm.RegisterFunction("handle_http", sm.handleHTTP)
	sm.vm.RegisterFunction("handle_ws", sm.handleWebSocket)
}

func (sm *ServerModule) createServer(L *lua.LState) int {
	serverID := L.CheckString(1)
	port := L.CheckInt(2)
	isHTTPS := L.OptBool(3, false)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	if isHTTPS {
		certFile := L.CheckString(4)
		keyFile := L.CheckString(5)

		config := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		server.TLSConfig = config

		
		server.TLSConfig.GetCertificate = func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load certificate: %v", err)
			}
			return &cert, nil
		}
	}

	sm.mu.Lock()
	sm.servers[serverID] = server
	sm.mu.Unlock()

	return 0
}

func (sm *ServerModule) startServer(L *lua.LState) int {
	serverID := L.CheckString(1)

	sm.mu.RLock()
	server, exists := sm.servers[serverID]
	sm.mu.RUnlock()

	if !exists {
		L.RaiseError("Server %s does not exist", serverID)
		return 0
	}

	go func() {
		var err error
		if server.TLSConfig != nil {
			err = server.ListenAndServeTLS("", "")
		} else {
			err = server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			sm.vm.monitor.handleError(fmt.Errorf("Server %s error: %v", serverID, err))
		}
	}()

	return 0
}

func (sm *ServerModule) stopServer(L *lua.LState) int {
	serverID := L.CheckString(1)

	sm.mu.Lock()
	server, exists := sm.servers[serverID]
	if exists {
		delete(sm.servers, serverID)
	}
	sm.mu.Unlock()

	if exists {
		if err := server.Close(); err != nil {
			sm.vm.monitor.handleError(fmt.Errorf("Error stopping server %s: %v", serverID, err))
		}
	}

	return 0
}

func (sm *ServerModule) handleHTTP(L *lua.LState) int {
	serverID := L.CheckString(1)
	path := L.CheckString(2)
	handler := L.CheckFunction(3)

	sm.mu.RLock()
	_, exists := sm.servers[serverID]
	sm.mu.RUnlock()

	if !exists {
		L.RaiseError("Server %s does not exist", serverID)
		return 0
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		L2 := lua.NewState()
		defer L2.Close()

		
		req := L2.NewTable()
		req.RawSetString("method", lua.LString(r.Method))
		req.RawSetString("path", lua.LString(r.URL.Path))
		req.RawSetString("query", lua.LString(r.URL.RawQuery))

		
		headers := L2.NewTable()
		for key, values := range r.Header {
			if len(values) > 0 {
				headers.RawSetString(key, lua.LString(values[0]))
			}
		}
		req.RawSetString("headers", headers)

		
		fn := L2.NewFunctionFromProto(handler.Proto)
		L2.Push(fn)
		L2.Push(req)
		if err := L2.PCall(1, 1, nil); err != nil {
			sm.vm.monitor.handleError(fmt.Errorf("HTTP handler error: %v", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		
		resp := L2.Get(-1)
		if resp.Type() != lua.LTTable {
			http.Error(w, "Invalid response from handler", http.StatusInternalServerError)
			return
		}

		
		if status := resp.(*lua.LTable).RawGetString("status"); status.Type() == lua.LTNumber {
			w.WriteHeader(int(status.(lua.LNumber)))
		}

		
		if headers := resp.(*lua.LTable).RawGetString("headers"); headers.Type() == lua.LTTable {
			headers.(*lua.LTable).ForEach(func(key, value lua.LValue) {
				w.Header().Set(key.String(), value.String())
			})
		}

		
		if body := resp.(*lua.LTable).RawGetString("body"); body.Type() == lua.LTString {
			w.Write([]byte(body.String()))
		}
	})

	return 0
}

func (sm *ServerModule) handleWebSocket(L *lua.LState) int {
	serverID := L.CheckString(1)
	path := L.CheckString(2)
	handler := L.CheckFunction(3)

	sm.mu.RLock()
	_, exists := sm.servers[serverID]
	sm.mu.RUnlock()

	if !exists {
		L.RaiseError("Server %s does not exist", serverID)
		return 0
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := sm.upgrader.Upgrade(w, r, nil)
		if err != nil {
			sm.vm.monitor.handleError(fmt.Errorf("WebSocket upgrade error: %v", err))
			return
		}
		defer conn.Close()

		L2 := lua.NewState()
		defer L2.Close()

		
		ws := L2.NewTable()
		ws.RawSetString("send", L2.NewFunction(func(L *lua.LState) int {
			message := L.CheckString(1)
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				sm.vm.monitor.handleError(fmt.Errorf("WebSocket send error: %v", err))
			}
			return 0
		}))

		ws.RawSetString("receive", L2.NewFunction(func(L *lua.LState) int {
			_, message, err := conn.ReadMessage()
			if err != nil {
				sm.vm.monitor.handleError(fmt.Errorf("WebSocket receive error: %v", err))
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LString(string(message)))
			return 1
		}))

		
		fn := L2.NewFunctionFromProto(handler.Proto)
		L2.Push(fn)
		L2.Push(ws)
		if err := L2.PCall(1, 0, nil); err != nil {
			sm.vm.monitor.handleError(fmt.Errorf("WebSocket handler error: %v", err))
		}
	})

	return 0
}
