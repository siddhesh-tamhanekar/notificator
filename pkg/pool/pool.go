package pool

import (
	"errors"
	"log"
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
	log.Println("add job start")
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
	pool.spawnWorker()
	log.Println("add job end")
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
	log.Println("spawnWorker start")
	if pool.runningWorkers > pool.MaxWorkers {
		log.Println("max limit reached cant spawn now")
		return
	}
	if len(pool.jobChannel) <= int(float64(pool.runningWorkers)) && pool.runningWorkers > pool.IdleWorkers {
		log.Printf("no need to spawn worker as jobs= %d  workers = %d \n", len(pool.jobChannel), int(float64(pool.runningWorkers)))
		return
	}
	log.Println("spawning new worker")
	go pool.worker(pool.runningWorkers)

	pool.mutex.Lock()
	pool.runningWorkers++
	pool.mutex.Unlock()
	log.Println("spawn worker end")
}

func (pool *Pool) worker(n int) {
	log.Printf("Worker %d Started\n", n)
	i := 0
	for {
		i++
		log.Printf("Worker %d In loop: %d\n", n, i)
		tick := time.Tick(time.Duration(pool.WorkerIdleTimeSecs) * time.Second)
		// tick = time.Tick(2 * time.Second)

		select {
		case job := <-pool.jobChannel:
			log.Printf("Worker %d Got job \n", n)
			job.startedAt = time.Now()
			job.Run()
			job.completedAt = time.Now()
			log.Printf("Worker %d Completed Job  \n", n)

			break
		case <-tick:
			if pool.runningWorkers > pool.IdleWorkers {
				log.Printf("Worker %d waited for idle time and no job available so dying.  \n", n)
				pool.mutex.Lock()
				pool.runningWorkers--
				pool.mutex.Unlock()
				// debug.FreeOSMemory()
				return
			}
		}
	}
}
