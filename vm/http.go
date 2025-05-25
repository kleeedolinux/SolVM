package vm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type HTTPModule struct {
	vm     *SolVM
	client *http.Client
}

func NewHTTPModule(vm *SolVM) *HTTPModule {
	return &HTTPModule{
		vm: vm,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (hm *HTTPModule) Register() {
	hm.vm.RegisterFunction("http_get", hm.get)
	hm.vm.RegisterFunction("http_post", hm.post)
	hm.vm.RegisterFunction("http_put", hm.put)
	hm.vm.RegisterFunction("http_delete", hm.delete)
	hm.vm.RegisterFunction("http_request", hm.request)
}

func (hm *HTTPModule) get(L *lua.LState) int {
	url := L.CheckString(1)

	resp, err := hm.client.Get(url)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP GET failed: %v", err))
		L.Push(lua.LNil)
		return 1
	}
	defer resp.Body.Close()

	return hm.handleResponse(L, resp)
}

func (hm *HTTPModule) post(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.CheckString(2)

	resp, err := hm.client.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP POST failed: %v", err))
		L.Push(lua.LNil)
		return 1
	}
	defer resp.Body.Close()

	return hm.handleResponse(L, resp)
}

func (hm *HTTPModule) put(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.CheckString(2)

	req, err := http.NewRequest("PUT", url, bytes.NewBufferString(body))
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP PUT failed: %v", err))
		L.Push(lua.LNil)
		return 1
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := hm.client.Do(req)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP PUT failed: %v", err))
		L.Push(lua.LNil)
		return 1
	}
	defer resp.Body.Close()

	return hm.handleResponse(L, resp)
}

func (hm *HTTPModule) delete(L *lua.LState) int {
	url := L.CheckString(1)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP DELETE failed: %v", err))
		L.Push(lua.LNil)
		return 1
	}

	resp, err := hm.client.Do(req)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP DELETE failed: %v", err))
		L.Push(lua.LNil)
		return 1
	}
	defer resp.Body.Close()

	return hm.handleResponse(L, resp)
}

func (hm *HTTPModule) request(L *lua.LState) int {
	method := L.CheckString(1)
	url := L.CheckString(2)

	
	var body io.Reader
	if L.GetTop() > 2 {
		body = bytes.NewBufferString(L.CheckString(3))
	}

	
	headers := make(map[string]string)
	if L.GetTop() > 3 {
		headersTable := L.CheckTable(4)
		headersTable.ForEach(func(key, value lua.LValue) {
			headers[key.String()] = value.String()
		})
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP %s request failed: %v", method, err))
		L.Push(lua.LNil)
		return 1
	}

	
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := hm.client.Do(req)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("HTTP %s request failed: %v", method, err))
		L.Push(lua.LNil)
		return 1
	}
	defer resp.Body.Close()

	return hm.handleResponse(L, resp)
}

func (hm *HTTPModule) handleResponse(L *lua.LState, resp *http.Response) int {
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		hm.vm.monitor.handleError(fmt.Errorf("Failed to read response body: %v", err))
		L.Push(lua.LNil)
		return 1
	}

	
	response := L.NewTable()
	response.RawSetString("status", lua.LNumber(resp.StatusCode))
	response.RawSetString("body", lua.LString(string(body)))

	
	headers := L.NewTable()
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers.RawSetString(key, lua.LString(values[0]))
		}
	}
	response.RawSetString("headers", headers)

	L.Push(response)
	return 1
}
