package taskpool

/**
 * 协程池
 */

type (
	Task struct {
		f TaskFun
		data interface{}
	}
	Pool struct {
		PoolNum int
		Task chan *Task
		Done chan bool
	}
	TaskFun func(data interface{}, pindex int)
)

func NewTask(data interface{}, f TaskFun) *Task {
	return &Task{f: f, data:data}
}

func NewPool(poolNum int) *Pool {
	return &Pool{PoolNum: poolNum, Task: make(chan *Task), Done: make(chan bool, poolNum)}
}

func (p *Pool) AddTask(task *Task) {
	p.Task <- task
}

func (p *Pool) Run() *Pool {
	for i := 0; i < p.PoolNum; i++ {
		go func(pindex int) {
			for {
				task, ok := <-p.Task
				if ok {
					task.f(task.data, pindex)
				} else {
					//log.Printf("pindex=[%d] closed\n", pindex)
					break
				}
			}
			p.Done <- true
		}(i)
	}
	return p
}

func (p *Pool) Close() {
	close(p.Task)
	for i := 0; i < p.PoolNum; i++ {
		<-p.Done
	}
}
