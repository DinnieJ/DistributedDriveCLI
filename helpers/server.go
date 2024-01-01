package helpers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

const DEFAULT_PORT = 28899

const SRV_KEY = "server"

type Credential struct {
}

func StoreToken(cred *Credential) error {
	return nil
}
func CallbackWrapper(resp http.ResponseWriter, request *http.Request) {

	var ctx = request.Context()
	var cancel = ctx.Value("cancel").(context.CancelFunc)
	var server *http.Server = ctx.Value(SRV_KEY).(*http.Server)
	select {
	case <-ctx.Done():
		{
			if err := server.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
			return
		}
	default:
	}
	var code = request.URL.Query().Get("code")
	if code != "" {
		fmt.Printf("Find code: %s", code)
		io.WriteString(resp, "Fuck you\n")
		cancel()
	}
}

func StartCallbackServer(port int) error {
	var listenerPort int
	if port != 0 {
		listenerPort = port
	} else {
		listenerPort = DEFAULT_PORT
	}

	httpServerExitDone := &sync.WaitGroup{}

	httpServerExitDone.Add(1)

	var mux = http.NewServeMux()
	mux.HandleFunc("/callback", CallbackWrapper)
	var ctx, cancel = context.WithCancel(context.Background())
	var addr = Spr("127.0.0.1:%d", listenerPort)
	LogInfo.Printf("Start callback server at %s\n", addr)
	var server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	server.BaseContext = func(l net.Listener) context.Context {
		ctx = context.WithValue(ctx, SRV_KEY, server)
		ctx = context.WithValue(ctx, "cancel", cancel)
		return ctx
	}
	go func() {
		defer httpServerExitDone.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	httpServerExitDone.Wait()
	return nil
}
