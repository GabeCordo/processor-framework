package http

import (
	"context"
	"fmt"
	"github.com/GabeCordo/mango-go/processor/threads/common"
	"github.com/GabeCordo/mango/threads"
	"net/http"
	"net/http/pprof"
	"time"
)

func (httpThread *Thread) Setup() {

	mux := http.NewServeMux()

	mux.HandleFunc("/module", func(w http.ResponseWriter, r *http.Request) {
		httpThread.moduleCallback(w, r)
	})

	mux.HandleFunc("/supervisor", func(w http.ResponseWriter, r *http.Request) {
		httpThread.supervisorCallback(w, r)
	})

	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		httpThread.debugCallback(w, r)
	})

	// TODO - explore this more, fucking cool
	if common.GetConfigInstance().Debug {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	httpThread.mux = mux
}

func (httpThread *Thread) Start() {
	httpThread.wg.Add(1)

	go func(thread *Thread) {
		net := common.GetConfigInstance().Net
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", net.Host, net.Port), httpThread.mux)
		if err != nil {
			thread.Interrupt <- threads.Panic
		}
	}(httpThread)

	go func() {
		for supervisorResponse := range httpThread.C2 {
			if !httpThread.accepting {
				break
			}
			httpThread.ProvisionerResponseTable.Write(supervisorResponse.Nonce, supervisorResponse)
		}
	}()

	httpThread.wg.Wait()
}

func (httpThread *Thread) Teardown() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	err := httpThread.server.Shutdown(ctx)
	if err != nil {
		httpThread.Interrupt <- threads.Panic
	}
}
