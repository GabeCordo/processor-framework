package http

import (
	"context"
	"github.com/GabeCordo/mango/threads"
	"net/http"
	"net/http/pprof"
	"time"
)

func (thread *Thread) Setup() {

	mux := http.NewServeMux()

	mux.HandleFunc("/module", func(w http.ResponseWriter, r *http.Request) {
		thread.moduleCallback(w, r)
	})

	mux.HandleFunc("/supervisor", func(w http.ResponseWriter, r *http.Request) {
		thread.supervisorCallback(w, r)
	})

	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		thread.debugCallback(w, r)
	})

	// TODO - explore this more, fucking cool
	if thread.Config.Debug {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	thread.mux = mux
}

func (thread *Thread) Start() {
	thread.wg.Add(1)

	go func(thread *Thread) {
		err := http.ListenAndServe(thread.Config.Net, thread.mux)
		if err != nil {
			thread.Interrupt <- threads.Panic
		}
	}(thread)

	go func() {
		for supervisorResponse := range thread.C2 {
			if !thread.accepting {
				break
			}
			thread.ProvisionerResponseTable.Write(supervisorResponse.Nonce, supervisorResponse)
		}
	}()

	thread.wg.Wait()
}

func (thread *Thread) Teardown() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	err := thread.server.Shutdown(ctx)
	if err != nil {
		thread.Interrupt <- threads.Panic
	}
}
