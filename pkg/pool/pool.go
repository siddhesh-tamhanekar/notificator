package pool

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

var pool *Pool

//GetInstance provides single instance of pool
func GetInstance() *Pool {
	if pool == nil {
		pool = &Pool{
			mutex: &sync.Mutex{},
		}
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

type Pool struct {
	MaxWorkers         int
	IdleWorkers        int
	runningWorkers     int
	JobQueueCapacity   int
	WorkerIdleTimeSecs int
	jobChannel         chan Job
	avgRunningTime     int
	mutex              *sync.Mutex
}

// AddJob adds the job  into queue.
func (pool *Pool) AddJob(job Job) error {
	if len(pool.jobChannel) == pool.JobQueueCapacity {
		time.Sleep(time.Millisecond * 100)
		if len(pool.jobChannel) == pool.JobQueueCapacity {
			return errors.New("Job queue Full please try again lataer")
		}
	}
	pool.jobChannel <- job
	pool.spawnWorker()
	return nil
}

// Run method starts the worker pool
func (pool *Pool) Run() {
	pool.jobChannel = make(chan Job, pool.JobQueueCapacity)
	// lets start ideal workers.
	for i := 0; i < pool.IdleWorkers; i++ {
		pool.spawnWorker()
	}
}

func (pool *Pool) Stats() interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"running Workers":    pool.runningWorkers,
		"max":                pool.MaxWorkers,
		"jobs":               len(pool.jobChannel),
		"current goroutines": runtime.NumGoroutine(),
		"alloc memory":       bToM(m.Alloc),
		"total alloc memory": bToM(m.TotalAlloc),
		"sys memory":         bToM(m.Sys),
		"num gc":             m.NumGC,
	}
}
func bToM(b uint64) uint64 {
	return b / 1024
}

func (pool *Pool) spawnWorker() {
	if pool.runningWorkers > pool.MaxWorkers {
		return
	}
	if len(pool.jobChannel) <= int(float64(pool.runningWorkers)) {
		return
	}

	go worker(pool)

	pool.mutex.Lock()
	pool.runningWorkers++
	pool.mutex.Unlock()
}

func worker(pool *Pool) {
	for {
		tick := time.Tick(time.Duration(pool.WorkerIdleTimeSecs) * time.Second)
		// tick = time.Tick(2 * time.Second)

		select {
		case job := <-pool.jobChannel:
			job.startedAt = time.Now()
			job.Run()
			job.completedAt = time.Now()

			break
		case <-tick:
			// if pool.runningWorkers > pool.IdleWorkers {
			pool.mutex.Lock()
			pool.runningWorkers--
			pool.mutex.Unlock()
			return
			// }
		}
	}
}
