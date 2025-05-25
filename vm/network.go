package vm

import (
	"fmt"
	"net"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type NetworkModule struct {
	vm       *SolVM
	tcpConns map[int]net.Conn
	udpConns map[int]*net.UDPConn
	mu       sync.RWMutex
	nextID   int
}

func NewNetworkModule(vm *SolVM) *NetworkModule {
	return &NetworkModule{
		vm:       vm,
		tcpConns: make(map[int]net.Conn),
		udpConns: make(map[int]*net.UDPConn),
		nextID:   1,
	}
}

func (nm *NetworkModule) Register() {
	nm.vm.RegisterFunction("tcp_listen", nm.tcpListen)
	nm.vm.RegisterFunction("tcp_connect", nm.tcpConnect)
	nm.vm.RegisterFunction("udp_sendto", nm.udpSendTo)
	nm.vm.RegisterFunction("udp_recvfrom", nm.udpRecvFrom)
	nm.vm.RegisterFunction("resolve_dns", nm.resolveDNS)
}

func (nm *NetworkModule) tcpListen(L *lua.LState) int {
	port := L.CheckInt(1)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		L.RaiseError("Failed to listen: %v", err)
		return 0
	}

	nm.mu.Lock()
	id := nm.nextID
	nm.nextID++
	nm.mu.Unlock()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				nm.vm.monitor.handleError(fmt.Errorf("Accept error: %v", err))
				return
			}

			nm.mu.Lock()
			connID := nm.nextID
			nm.nextID++
			nm.tcpConns[connID] = conn
			nm.mu.Unlock()

			
			connTable := L.NewTable()
			connTable.RawSetString("id", lua.LNumber(connID))
			connTable.RawSetString("remote_addr", lua.LString(conn.RemoteAddr().String()))

			
			connTable.RawSetString("read", L.NewFunction(func(L *lua.LState) int {
				buffer := make([]byte, 1024)
				n, err := conn.Read(buffer)
				if err != nil {
					L.Push(lua.LNil)
					L.Push(lua.LString(err.Error()))
					return 2
				}
				L.Push(lua.LString(string(buffer[:n])))
				return 1
			}))

			
			connTable.RawSetString("write", L.NewFunction(func(L *lua.LState) int {
				data := L.CheckString(1)
				_, err := conn.Write([]byte(data))
				if err != nil {
					L.Push(lua.LBool(false))
					L.Push(lua.LString(err.Error()))
					return 2
				}
				L.Push(lua.LBool(true))
				return 1
			}))

			
			connTable.RawSetString("close", L.NewFunction(func(L *lua.LState) int {
				conn.Close()
				nm.mu.Lock()
				delete(nm.tcpConns, connID)
				nm.mu.Unlock()
				return 0
			}))

			L.Push(connTable)
		}
	}()

	L.Push(lua.LNumber(id))
	return 1
}

func (nm *NetworkModule) tcpConnect(L *lua.LState) int {
	host := L.CheckString(1)
	port := L.CheckInt(2)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		L.RaiseError("Failed to connect: %v", err)
		return 0
	}

	nm.mu.Lock()
	id := nm.nextID
	nm.nextID++
	nm.tcpConns[id] = conn
	nm.mu.Unlock()

	
	connTable := L.NewTable()
	connTable.RawSetString("id", lua.LNumber(id))
	connTable.RawSetString("remote_addr", lua.LString(conn.RemoteAddr().String()))

	
	connTable.RawSetString("read", L.NewFunction(func(L *lua.LState) int {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LString(string(buffer[:n])))
		return 1
	}))

	
	connTable.RawSetString("write", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		_, err := conn.Write([]byte(data))
		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LBool(true))
		return 1
	}))

	
	connTable.RawSetString("close", L.NewFunction(func(L *lua.LState) int {
		conn.Close()
		nm.mu.Lock()
		delete(nm.tcpConns, id)
		nm.mu.Unlock()
		return 0
	}))

	L.Push(connTable)
	return 1
}

func (nm *NetworkModule) udpSendTo(L *lua.LState) int {
	addr := L.CheckString(1)
	port := L.CheckInt(2)
	message := L.CheckString(3)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		L.RaiseError("Failed to create UDP connection: %v", err)
		return 0
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		L.RaiseError("Failed to send UDP message: %v", err)
		return 0
	}

	return 0
}

func (nm *NetworkModule) udpRecvFrom(L *lua.LState) int {
	port := L.CheckInt(1)

	addr := &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		L.RaiseError("Failed to listen on UDP port: %v", err)
		return 0
	}

	nm.mu.Lock()
	id := nm.nextID
	nm.nextID++
	nm.udpConns[id] = conn
	nm.mu.Unlock()

	
	connTable := L.NewTable()
	connTable.RawSetString("id", lua.LNumber(id))

	
	connTable.RawSetString("receive", L.NewFunction(func(L *lua.LState) int {
		buffer := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		result := L.NewTable()
		result.RawSetString("data", lua.LString(string(buffer[:n])))
		result.RawSetString("addr", lua.LString(remoteAddr.IP.String()))
		result.RawSetString("port", lua.LNumber(remoteAddr.Port))

		L.Push(result)
		return 1
	}))

	
	connTable.RawSetString("close", L.NewFunction(func(L *lua.LState) int {
		conn.Close()
		nm.mu.Lock()
		delete(nm.udpConns, id)
		nm.mu.Unlock()
		return 0
	}))

	L.Push(connTable)
	return 1
}

func (nm *NetworkModule) resolveDNS(L *lua.LState) int {
	hostname := L.CheckString(1)

	ips, err := net.LookupIP(hostname)
	if err != nil {
		L.RaiseError("Failed to resolve DNS: %v", err)
		return 0
	}

	if len(ips) == 0 {
		L.RaiseError("No IP addresses found for hostname")
		return 0
	}

	L.Push(lua.LString(ips[0].String()))
	return 1
}
