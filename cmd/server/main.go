package main

import (
	"fmt"
	"log"
	"net/http"
	"siddhesh-tamhanekar/notificator"
	"siddhesh-tamhanekar/notificator/internal"
	"siddhesh-tamhanekar/notificator/pkg/pool"
)

func main() {
	// initialize worker pool
	pool := pool.GetInstance()

	pool.JobQueueCapacity = notificator.ConfigInstance().Pool.JobQueueCapacity
	pool.MaxWorkers = notificator.ConfigInstance().Pool.MaxWorkers
	go pool.Run()

	server := http.NewServeMux()
	// server.HandleFunc("/debug/pprof/", pprof.Index)
	// server.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// server.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// server.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// server.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// server.Handle("/debug/pprof/block", pprof.Handler("block"))
	// server.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	// server.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	// server.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))

	// add routes
	internal.SetHandlers(server)

	host := notificator.ConfigInstance().Server.Host

	fmt.Println("Server is started on " + host)
	if err := http.ListenAndServe(host, server); err != nil {
		log.Fatal("web server could not be started " + err.Error())
	}
}
