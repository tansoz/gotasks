package gotasks

import (
	"net/http"
	"testing"
)

func TestNewTasks(t *testing.T) {
	NewTasks(5, func() func(Tasks) {
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		})
		http.ListenAndServe(":7575", nil)
		return func(t Tasks) {

		}
	}(), nil)
}
