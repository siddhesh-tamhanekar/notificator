package pool

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"
)

var pool *Pool

//GetInstance provides single instance of pool
func GetInstance() *Pool {
	if pool == nil {
		pool = &Pool{}
	}
	return pool
}

// Job
type Job struct {
	Id          int
	Run         func()
	CreatedAt   time.Time
	startedAt   time.Time
	completedAt time.Time
}

// Pool provides worker pool implementation
type Pool struct {
	MaxWorkers       int
	runningWorkers   chan struct{}
	JobQueueCapacity int
	avgRunningTime   float64
	jobProcessed     int
	jobChannel       chan Job
}

// AddJob adds the job  into queue.
func (pool *Pool) AddJob(job Job) error {
	if len(pool.jobChannel) == pool.JobQueueCapacity {
		log.Println("job queue was fool let's stop accepting for 1 second")
		time.Sleep(time.Millisecond * 1000)
		if len(pool.jobChannel) == pool.JobQueueCapacity {
			log.Println("after pausing for second still job queue is full denying to add the job")
			return errors.New("Job queue Full please try again lataer")
		}
	}
	pool.jobChannel <- job
	log.Println("job added")
	return nil
}
func (pool *Pool) canSpawn() bool {
	if len(pool.runningWorkers) >= pool.MaxWorkers {
		return false
	}
	return true
}

// Run method starts the worker pool
func (pool *Pool) Run() {
	pool.jobChannel = make(chan Job, pool.JobQueueCapacity)
	pool.runningWorkers = make(chan struct{}, pool.MaxWorkers)
	for {
		if pool.canSpawn() == false {
			fmt.Println("pool max workres reach not able to do job")
			time.Sleep(time.Second * 1)
			continue
		}

		job := <-pool.jobChannel
		pool.runningWorkers <- struct{}{}
		go func() {
			job.startedAt = time.Now()
			job.Run()
			job.completedAt = time.Now()
			<-pool.runningWorkers
		}()
	}

}

// Stats provides information about pool and memory usage.
func (pool *Pool) Stats() interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"Running Workers":        len(pool.runningWorkers),
		"Max Workers":            pool.MaxWorkers,
		"Jobs in the Queue":      len(pool.jobChannel),
		"current goroutines":     runtime.NumGoroutine(),
		"Allocated memory":       bToM(m.Alloc),
		"Total Allocated memory": bToM(m.TotalAlloc),
		"Sytstem memory":         bToM(m.Sys),
		"no of times GC Runs":    m.NumGC,
	}
}

func bToM(bytes uint64) uint64 {
	return bytes / 1024
}
