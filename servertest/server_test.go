package servertest

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/hypnoglow/x/server"
)

func TestServer(t *testing.T) {
	addr := fmt.Sprintf(":%d", getFreePort())
	handler := http.HandlerFunc(testHandler)
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	t.Run("Should execute standard flow", func(t *testing.T) {
		gsrv := server.New(addr, handler, logger)
		go gsrv.Start()

		go func() {
			defer gsrv.Stop()

			body, err := getBody("http://" + addr)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			if body != "Just testing!" {
				t.Fatalf("Unexpected response body: %s", string(body))
			}
		}()

		gsrv.Wait()
		gsrv.Shutdown()
	})
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Just testing!")
}

func getBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getFreePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}
