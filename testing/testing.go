/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"net/http"
	"net/http/httptest"
	"time"

	. "gopkg.in/check.v1"
)

type MockServer struct {
	HandlerCalled bool
	Body          string
	ResponseCode  int
	Server        *httptest.Server
	URL           string
	CapturedPath  string
	Headers       map[string]string
}

func NewMockServer() *MockServer {
	return &MockServer{
		HandlerCalled: false,
		ResponseCode:  200,
		Headers:       map[string]string{},
	}
}

func (m *MockServer) Start(c *C) *MockServer {
	m.HandlerCalled = false
	m.Server = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range m.Headers {
				w.Header().Set(key, value)
			}
			w.WriteHeader(m.ResponseCode)
			w.Write([]byte(m.Body))
			m.CapturedPath = r.URL.Path
			m.HandlerCalled = true
		}))

	status := 0
	checkStarted := time.Now()
	for status != m.ResponseCode && time.Now().Before(checkStarted.Add(time.Second*10)) {
		time.Sleep(time.Second / 1000)
		resp, err := http.Get(m.Server.URL)
		if err == nil {
			status = resp.StatusCode
		}
	}
	c.Assert(status, Not(Equals), 0)
	m.HandlerCalled = false
	m.URL = m.Server.URL
	return m
}

func (m *MockServer) WithBody(body string) *MockServer {
	m.Body = body
	return m
}

func (m *MockServer) WithResponseCode(code int) *MockServer {
	m.ResponseCode = code
	return m
}

func (m *MockServer) WithHeader(key, value string) *MockServer {
	m.Headers[key] = value
	return m
}

func (m *MockServer) Stop() {
	m.Server.Close()
}

func (m *MockServer) ExpectCalled(c *C, expectCalled bool, path string) {
	c.Assert(m.HandlerCalled, Equals, expectCalled)
	c.Assert(m.CapturedPath, Equals, path)
}

func (m *MockServer) ExpectError(c *C, err error, path, expectedErrorString string) {
	m.ExpectCalled(c, true, path)
	c.Assert(err.Error(), Equals, expectedErrorString)
}
