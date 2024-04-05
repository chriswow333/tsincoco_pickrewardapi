// Package goroutine provides utilities regarding to a goroutine
package goroutine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"
	"time"

	"pickrewardapi/base/ctx"

	"github.com/sirupsen/logrus"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	// An empty inform function placeholder for testing purpose.
	logger = logrus.StandardLogger()
)

var (
	// ErrScheduleTimeout returned by Pool to indicate that there is no free
	// goroutines during some period of time.
	ErrScheduleTimeout = fmt.Errorf("schedule timed out")
	// ErrSchedulePause indicated that pool is closed
	ErrSchedulePause = fmt.Errorf("schedule pause")
)

const (
	defaultTimeoutSec = 20
)

type (
	// Strategy is a part of options to lead scheduling
	Strategy   uint
	poolOption struct {
		timeoutSec    int
		strategy      Strategy
		enableMonitor bool
	}

	// PoolOption is an alias for functional argument in Pool
	PoolOption func(*poolOption)
)

const (
	// Default means puting tasks in queue or initializing new goroutine randomly
	Default Strategy = iota
	// QueueFirst means puting tasks in queue at first time or doing as default
	QueueFirst
)

// Pool contains logic of goroutine reuse.
type Pool struct {
	name       string
	maxSize    int
	workerChan chan struct{}
	workers    []*worker

	queueChan chan func()
	pauseChan chan struct{}
	pauseOnce sync.Once
	workerMut sync.RWMutex
	counter   *counter

	// options
	strategy Strategy
}

// WithStrategy specifies strategy for scheduling
func WithStrategy(s Strategy) PoolOption {
	return func(o *poolOption) {
		o.strategy = s
	}
}

// DisableMonitor disable goroutine which monitors pool
func DisableMonitor() PoolOption {
	return func(o *poolOption) {
		o.enableMonitor = false
	}
}

func initPoolOption() *poolOption {
	return &poolOption{
		enableMonitor: true,
	}
}

// NewPool creates new asynchronously goroutine pool with given size. It also creates a work
// queue of given size. Finally, it spawns given amount of goroutines immediately.
func NewPool(name string, maxWorkerSize, queueSize, initWorkerSize int, options ...PoolOption) *Pool {
	if maxWorkerSize < 0 || queueSize < 0 || initWorkerSize < 0 {
		panic("negative numbers")
	}
	if initWorkerSize <= 0 && queueSize > 0 {
		panic("dead queue configuration detected")
	}
	if initWorkerSize > maxWorkerSize {
		panic("spawn > workers")
	}

	// load options
	o := initPoolOption()
	for _, opt := range options {
		opt(o)
	}

	p := &Pool{
		name:       name,
		maxSize:    maxWorkerSize,
		workerChan: make(chan struct{}, maxWorkerSize),
		queueChan:  make(chan func(), queueSize),
		pauseChan:  make(chan struct{}),
		counter:    &counter{},
		strategy:   o.strategy,
	}

	wg := &sync.WaitGroup{}
	p.workerMut.Lock()
	for i := 0; i < initWorkerSize; i++ {
		p.workerChan <- struct{}{}
		wg.Add(1)
		p.workers = append(p.workers, newWorker(p.counter, p.queueChan, func() {
			defer wg.Done()
		}))
	}
	p.workerMut.Unlock()

	// make sure all workers needing to be initialized ready
	wg.Wait()

	// start monitor goroutine
	if o.enableMonitor {
		Go(func() {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					// fmt.Println("pool.ActiveWorker", float64(p.counter.ActiveWorker()), "name", name)
					// fmt.Println("pool.Worker", float64(p.GetSize()), "name", name)
				case <-p.pauseChan:
					return
				}
			}
		})
	}

	return p
}

// Schedule schedules task to be executed over pool's workers.
func (p *Pool) Schedule(task func()) error {
	return p.schedule(task, nil)
}

// ScheduleTimeout schedules task to be executed over pool's workers.
// It returns ErrScheduleTimeout when no free workers met during given timeout.
func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	return p.schedule(task, timer.C)
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) (err error) {
	// defer fmt.Println("schedule.time", "name", p.name)

	select {
	case <-p.pauseChan:
		return ErrSchedulePause
	default:
	}

	if p.strategy == QueueFirst {
		// Lazy loading for goroutine: put tasks into queue if possible
		select {
		case <-timeout:
			return ErrScheduleTimeout
		case <-p.pauseChan:
			return ErrSchedulePause
		case p.queueChan <- task:
			return nil
		default:
		}
	}

	select {
	case <-timeout:
		return ErrScheduleTimeout
	case <-p.pauseChan:
		return ErrSchedulePause
	case p.queueChan <- task:
		return nil
	case p.workerChan <- struct{}{}:
		p.workerMut.Lock()
		defer p.workerMut.Unlock()

		select {
		case <-p.pauseChan:
			return ErrSchedulePause
		default:
			p.workers = append(p.workers, newWorker(p.counter, p.queueChan, task))
			return nil
		}
	}
}

// GetSize returns the current size of the pool.
func (p *Pool) GetSize() int {
	p.workerMut.RLock()
	defer p.workerMut.RUnlock()

	return len(p.workers)
}

// Close will terminate all workers and close the job channel of this Pool.
func (p *Pool) Close() {
	// only allow calling Close() or GracefulClose() once
	p.pauseOnce.Do(func() {
		p.workerMut.Lock()
		defer p.workerMut.Unlock()
		// Stop opening new workers
		close(p.pauseChan)

		runningWorkers := len(p.workers)

		for i := 0; i < runningWorkers; i++ {
			p.workers[i].stop()
		}

		for i := 0; i < runningWorkers; i++ {
			p.workers[i].join()
		}

		p.workers = p.workers[:0]
	})
}

// WithTimeoutSec specifies the timeout seconds used in GracefulClose()
// default timeout is 20 seconds
func WithTimeoutSec(sec int) PoolOption {
	return func(o *poolOption) {
		o.timeoutSec = sec
	}
}

// GracefulClose will terminate all workers only if all tasks in queue are done
func (p *Pool) GracefulClose(options ...PoolOption) {
	o := &poolOption{}
	for _, opt := range options {
		opt(o)
	}

	if o.timeoutSec == 0 {
		o.timeoutSec = defaultTimeoutSec
	}
	// only allow calling Close() or GracefulClose() once
	p.pauseOnce.Do(func() {
		p.workerMut.Lock()
		defer p.workerMut.Unlock()
		// Stop opening new workers
		close(p.pauseChan)

		runningWorkers := len(p.workers)

		var wg sync.WaitGroup
		wg.Add(1)
		timeout := time.After(time.Duration(o.timeoutSec) * time.Second)
		p.workers = append(p.workers, newWorker(p.counter, p.queueChan, func() {
			defer wg.Done()
			time.Sleep(time.Second)
			// monitoring the size of queueChan to be zero
			select {
			case <-timeout:
				return
			case <-time.After(250 * time.Millisecond):
				if len(p.queueChan) == 0 {
					return
				}
			}
		}))

		for i := 0; i < runningWorkers; i++ {
			p.workers[i].gracefulStop()
		}

		// wait until the monitor returned successfully
		wg.Wait()

		// force all workers including the monitor to die
		runningWorkers = len(p.workers)
		for i := 0; i < runningWorkers; i++ {
			p.workers[i].stop()
		}

		for i := 0; i < runningWorkers; i++ {
			p.workers[i].join()
		}

		p.workers = p.workers[:0]
	})
}

type counter struct {
	active int
	mu     sync.RWMutex
}

func (a *counter) Add(num int) {
	a.mu.Lock()
	a.active += num
	a.mu.Unlock()
}

func (a *counter) ActiveWorker() (active int) {
	a.mu.RLock()
	active = a.active
	a.mu.RUnlock()
	return
}

type worker struct {
	counter   *counter
	queueChan <-chan func()
	// closeChan is used to shut down the worker
	// When it's true, close the worker immediately.
	// Or close the worker only if all tasks in queue are done
	closeChan  chan bool
	closedChan chan struct{}
}

func newWorker(counter *counter, queueChan <-chan func(), task func()) *worker {
	w := worker{
		counter:    counter,
		queueChan:  queueChan,
		closeChan:  make(chan bool),
		closedChan: make(chan struct{}),
	}
	Go(func() {
		w.run(task)
	})
	return &w
}

func (w *worker) cleanQueue() {
	for {
		select {
		case <-w.closeChan:
			return
		case task := <-w.queueChan:
			w.executeTask(task)
		}
	}
}

func (w *worker) run(task func()) {
	defer close(w.closedChan)

	w.executeTask(task)
	for {
		select {
		case forced := <-w.closeChan:
			if forced {
				return
			}
			w.cleanQueue()
			return
		default:
		}

		select {
		case task, ok := <-w.queueChan:
			if !ok {
				return
			}
			w.executeTask(task)
		case forced := <-w.closeChan:
			if forced {
				return
			}
			w.cleanQueue()
			return
		}
	}
}

func (w *worker) executeTask(task func()) {
	defer func() {
		w.counter.Add(-1)
		if r := recover(); r != nil {
			stack := Stack(4)
			logger.WithFields(logrus.Fields{
				"err":   r,
				"stack": string(stack),
			}).Error("panic")
		}
	}()
	w.counter.Add(1)
	task()
}

func (w *worker) stop() {
	// force to close the worker immediately
	w.closeChan <- true
}

func (w *worker) gracefulStop() {
	// close the worker gracefully
	w.closeChan <- false
}

func (w *worker) join() {
	<-w.closedChan
}

// PanicEvent contains the panic and stack trace info
type PanicEvent struct {
	Panic interface{}
	Stack []byte
}

// Go is a wrapper of Golang's "go" syntax that forks a goroutine and recovers potential
// panic function and logs.
// If panic occurs, it will sent to `panicCh`.
func Go(f func()) (panicCh chan *PanicEvent) {
	panicCh = make(chan *PanicEvent, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				stack := Stack(3)
				logger.WithFields(logrus.Fields{
					"err":   p,
					"stack": string(stack),
				}).Error("panic")
				panicCh <- &PanicEvent{Panic: p, Stack: stack}
			} else {
				close(panicCh)
			}
		}()
		f()
	}()
	return panicCh
}

// GoWithParameters is a ugly generic solution.
// But it allow developer to have pass-by-value method and no dependency to the outer layer local variable.
// Thus it will avoid performance issue and some race condition
// If panic occurs, it will sent to `panicCh`.
func GoWithParameters(f func(...interface{}), parameters ...interface{}) (panicCh chan *PanicEvent) {
	panicCh = make(chan *PanicEvent, 1)

	// context is not thread safe, so clone it conveniently.
	if len(parameters) > 0 {
		if context, ok := parameters[0].(ctx.CTX); ok {
			// key "clone" is useless, copy by WithValue() as it contains a memory pool.
			parameters[0] = ctx.WithValue(context, "clone", 1)
		}
	}

	go func(params ...interface{}) {
		defer func() {
			if p := recover(); p != nil {
				stack := Stack(3)
				logger.WithFields(logrus.Fields{
					"err":   p,
					"stack": string(stack),
				}).Error("panic")
				panicCh <- &PanicEvent{Panic: p, Stack: stack}
			} else {
				close(panicCh)
			}
		}()
		f(params...)
	}(parameters...)
	return panicCh
}

// Stack returns stack trace of current pc with skip levels
func Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
