package mango

import (
	"net/http"
	"testing"
)

func routingTestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Hello World!")
}

func routingATestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Server A")
}

func routingBTestServer(env Env) (Status, Headers, Body) {
	return 200, Headers{}, Body("Server B")
}

func TestRoutingSuccess(t *testing.T) {
	// Compile the stack
	routingStack := new(Stack)
	routes := make(map[string]App)
	routes["/a"] = routingATestServer
	routes["/b"] = routingBTestServer
	routingStack.Middleware(Routing(routes))
	routingApp := routingStack.Compile(routingTestServer)

	// Request against A
	request, err := http.NewRequest("GET", "http://localhost:3000/a", nil)
	status, _, body := routingApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "Server A"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}

	// Request against B
	request, err = http.NewRequest("GET", "http://localhost:3000/b", nil)
	status, _, body = routingApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected = "Server B"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func TestRoutingFailure(t *testing.T) {
	// Compile the stack
	routingStack := new(Stack)
	routes := make(map[string]App)
	routes["/a"] = routingATestServer
	routes["/b"] = routingBTestServer
	routingStack.Middleware(Routing(routes))
	routingApp := routingStack.Compile(routingTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/c", nil)
	status, _, body := routingApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "Hello World!"
	if string(body) != expected {
		t.Error("Expected body:", string(body), "to equal:", expected)
	}
}

func BenchmarkRouting(b *testing.B) {
	b.StopTimer()

	routingStack := new(Stack)
	routes := make(map[string]App)
	routes["/a"] = routingATestServer
	routes["/b"] = routingBTestServer
	routingStack.Middleware(Routing(routes))
	routingApp := routingStack.Compile(routingTestServer)

	request, _ := http.NewRequest("GET", "http://localhost:3000/a", nil)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		routingApp(Env{"mango.request": &Request{request}})
	}
	b.StopTimer()
}
