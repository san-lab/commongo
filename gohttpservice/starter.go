package gohttpservice

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

func Startserver(ha http.Handler) {
	StartserverWithAuth(ha, nil)
}

var DefPort = "8080"
var WithAuth = false

//Starts a server with the handler ha and the user/pass validator v
func StartserverWithAuth(ha http.Handler, v func(string, string) bool) {
	httpPort := flag.String("httpPort", DefPort, "http port")
	httpsPort := flag.Int("httpsPort", 0, "https port. tls not started if not provided. requires server.crt & server.key")
	withBasicAuth := flag.Bool("withAuth", WithAuth, "should Basic Authentication be enabled")
	flag.Parse()

	//Allow spanning go routines and wait for them to complete
	interruptChan := make(chan os.Signal)
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	ctx = context.WithValue(ctx, "WaitGroup", wg)
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/static/", http.StripPrefix("/static", fs).ServeHTTP)
	signal.Notify(interruptChan, os.Interrupt)

	//Optionally wrap the handler provided in Basic Auth
	if *withBasicAuth {
		//Provide your user/password validating function here
		if v == nil {
			v = func(s1 string, s2 string) bool { return true }
		}
		http.Handle("/", BasicAuthHandler(ha, v))
	} else {
		http.Handle("/", ha)
	}

	srv := http.Server{Addr: ":" + *httpPort}
	var tlsSrv http.Server
	if *httpsPort > 0 {
		tlsSrv = http.Server{Addr: ":" + strconv.Itoa(*httpsPort)}
	}

	go func() {
		select {
		case <-interruptChan:
			cancel()
			srv.Shutdown(context.TODO())
			if *httpsPort > 0 {
				tlsSrv.Shutdown(context.TODO())
			}
			return
		}
	}()

	if *httpsPort > 0 {
		go func() { log.Println(tlsSrv.ListenAndServeTLS("server.crt", "server.key")) }()
	}

	log.Println(srv.ListenAndServe())
	wg.Wait()
}

func BasicAuthHandler(h http.Handler, validate func(string, string) bool) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			username, password, authOK := r.BasicAuth()
			if authOK && validate(username, password) {
				h.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Not authorized", 401)
		})

}
