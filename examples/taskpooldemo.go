package main

import (
	"gotools/taskpool"
	"log"
	"time"
)

/**
 * 协程池示例(一次处理一个任务)
 */

func main() {

	starttime := time.Now()

	poolNum := 3  //开启协程数
	pool := taskpool.NewPool(poolNum).Run()

	taskcount := 6  //任务数量
	var i int
	for i = 1; i <= taskcount; i++ {
		task := taskpool.NewTask(i, func(data interface{}, pindex int) {
			taskindex := data.(int)
			time.Sleep(1*time.Second)
			log.Printf("pindex=%v, run task=%d", pindex, taskindex)
		})
		pool.AddTask(task)
		log.Printf("add task %d", i)
	}

	pool.Close()

	cost := time.Since(starttime)
	log.Printf("总耗时=[%v]", cost)
}