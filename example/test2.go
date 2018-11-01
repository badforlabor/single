/**
 * Auth :   liubo
 * Date :   2018/11/1 17:16
 * Comment: 模拟更复杂的多线程
 */

package main

import (
	"fmt"
	"math/rand"
	"single"
	"time"
)

func complexFunction(tag string, worker *single.Worker, args map[string]string) {

	callback := func(...interface{}) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		fmt.Println("step-2", tag, args)
	}

	fmt.Println("step-1", tag, args)
	worker.BlockJob(single.MakeCommonAction(callback))
	fmt.Println("step-3", tag, args)
}
func test2() {
	fmt.Println("test2")

	worker := &single.Worker{}
	worker.Run()

	// 模拟100个线程
	action := func(index int, count int) {
		for i:= index;i< index+count; i++ {
			args := make(map[string]string)
			v := fmt.Sprint(rand.Intn(10000))
			args[v]=v
			go complexFunction(fmt.Sprint("tag", i), worker, args)
		}
	}
	go action(0, 30)
	go action(30, 30)
	go action(60, 30)

	time.Sleep(10 * time.Second)
	fmt.Println("done")
}
