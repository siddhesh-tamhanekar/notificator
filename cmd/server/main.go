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

	pool.IdleWorkers = notificator.ConfigInstance().Pool.IdleWorkers
	pool.JobQueueCapacity = notificator.ConfigInstance().Pool.JobQueueCapacity
	pool.MaxWorkers = notificator.ConfigInstance().Pool.MaxWorkers
	pool.WorkerIdleTimeSecs = notificator.ConfigInstance().Pool.WorkerIdleTimeSecs
	pool.Run()

	server := http.NewServeMux()
	// add routes
	internal.SetHandlers(server)

	host := notificator.ConfigInstance().Server.Host

	fmt.Println("Server is started on " + host)
	if err := http.ListenAndServe(host, server); err != nil {
		log.Fatal("web server could not be started " + err.Error())
	}
}
