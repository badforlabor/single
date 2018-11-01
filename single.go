/**
 * Auth :   liubo
 * Date :   2018/11/1 16:18
 * Comment: 单线程干活儿工具。暂时没用ringbuffer，以后优化吧。
 * 	计划优化：
 *		ringbuffer实现无锁形式
 */

package single

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	jobs  []*Action
	mutex sync.Mutex
	stop  bool

	// 测试用的
	rid int
	wid int
}

func (self *Worker) Run() {
	go self.loop()
}

// 一个是阻塞式的
func (self *Worker) BlockJob(callback IAction) {
	job := self.addJob(callback)
	<-job.done
}

// 一个是非阻塞式的
func (self *Worker) NonblockJob(callback IAction) {
	self.addJob(callback)
}
func (self *Worker) addJob(callback IAction) *Action {

	job := &Action{callback: callback, done: make(chan bool)}

	self.mutex.Lock()
	self.jobs = append(self.jobs, job)
	self.rid++
	job.id = self.rid
	self.mutex.Unlock()

	return job
}
func (self *Worker) loop() {
	for !self.stop {
		var job *Action
		wid := 0
		self.mutex.Lock()
		if len(self.jobs) > 0 {
			job = self.jobs[0]
			self.jobs = self.jobs[1:]
		}
		self.wid++
		wid = self.wid
		self.mutex.Unlock()

		if job == nil {
			time.Sleep(time.Millisecond * 30)
			continue
		}

		if wid != self.wid {
			panic(fmt.Sprint("错误！并非单线程的！", wid, self.wid))
		}

		job.safeExec()
		job.done <- true
	}
}
