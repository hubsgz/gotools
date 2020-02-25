package main

import (
	"fmt"
	"gotools/bulktaskpool"
	"log"
	"strings"
	"time"
)

/**
 * 协程池示例(一次处理多个任务)
 */

func main() {

	starttime := time.Now()

	poolNum := 3  //开启协程数
	bulkSize := 8 //一次处理的任务数
	pool := bulktaskpool.NewBulkPool(poolNum, bulkSize).Run()

	taskcount := 30  //任务数量
	var i int
	for i = 1; i <= taskcount; i++ {
		task := bulktaskpool.NewBulkTask(i, func(datas []interface{}, pindex int) {
			arr := make([]string, len(datas))
			for k,data := range datas {
				taskindex := data.(int)
				arr[k] = fmt.Sprintf("%d", taskindex)
			}
			time.Sleep(1*time.Second)
			log.Printf("pindex=%v, run task=%v", pindex, strings.Join(arr, "|"))
		})
		pool.AddBulkTask(task)
		log.Printf("add task %d", i)
	}

	pool.Close()

	cost := time.Since(starttime)
	log.Printf("总耗时=[%v]", cost)
}