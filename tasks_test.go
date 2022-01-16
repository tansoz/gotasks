package gotasks

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
)

func TestNewTasks(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	NewTasks(5, func() func(Tasks) {
		return func(t Tasks) {
			http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				t.AddSync(func() {
					writer.Write([]byte("Hello World!"))
					fmt.Println("跑...")
				})
				fmt.Println("离开...")
			})
			http.ListenAndServe(":7575", nil)
		}
	}(), nil)
	wg.Wait()
}
func TestNewTasks2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	NewTasks(5, func() func(Tasks) {
		listener, err := net.Listen("tcp", ":7575")
		if err != nil {
			panic(err)
		}
		return func(t Tasks) {
			for {
				if client, accept_err := listener.Accept(); accept_err == nil {
					t.Add(func(conn net.Conn) func() {
						return func() {
							fmt.Println(t.Active())
							conn.Write([]byte("HTTP/1.0 200 OK\r\nContent-length: 12\r\n\r\nHello World!"))
						}
					}(client))
				}
			}
		}
	}(), nil)
	wg.Wait()
}
