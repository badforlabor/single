/**
 * Auth :   liubo
 * Date :   2018/11/1 14:51
 * Comment: 每个work都是一个单线程做事儿的
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"single"
)

func main() {

	fmt.Println("test!")

	var w single.Worker
	w.Run()

	cnt := make(chan int, 100)

	logs := ""
	callback := func(a int) {
		logs += fmt.Sprintf("%d\n", a)
		cnt <- a
	}

	actions := func(n int) {
		for i:=n; i< 100+n; i++ {

			callback := func(args ...interface{}) {
				ii := args[0].(int)
				fmt.Println(ii)
				callback(ii)
			}

			w.BlockJob(single.MakeCommonAction(callback, i))
		}
	}

	go func() {
		actions(1)
		actions(100)
		actions(200)
		actions(300)
	}()

	stop := make(chan bool)
	s2 := make(chan bool)
	go func() {
		stop2 := false
		for !stop2 {
			select {
			case a := <-cnt:
				fmt.Println("channel", a)

			case <-stop:
				stop2 = true

			default:
				time.Sleep(time.Millisecond)
			}
		}
		s2 <- true
	}()

	time.Sleep(13 * time.Second)
	stop <- true
	<-s2

	ioutil.WriteFile("1.log", []byte(logs), os.ModePerm)
	fmt.Println("done.")
}
