package bulktaskpool

import (
	"time"
)

/**
 * 协程池(批量任务)
 */

type (
	BulkTask struct {
		data interface{}
		f    BulkTaskFun
	}
	BulkPool struct {
		PoolNum  int
		BulkSize int
		BulkTask chan *BulkTask
		Stoped   chan bool
	}
	BulkData struct {
		datas []interface{}
		len   int
		f     BulkTaskFun
	}
	BulkTaskFun func(datas []interface{}, pindex int)
)

func NewBulkData(datas []interface{}, len int) *BulkData {
	return &BulkData{datas: datas, len: len}
}

func NewBulkTask(data interface{}, f BulkTaskFun) *BulkTask {
	return &BulkTask{data: data, f: f}
}

func NewBulkPool(poolNum int, bulkSize int) *BulkPool {
	return &BulkPool{PoolNum: poolNum, BulkSize:bulkSize, BulkTask: make(chan *BulkTask), Stoped: make(chan bool)}
}

func (p *BulkPool) AddBulkTask(task *BulkTask) {
	p.BulkTask <- task
}

func resetBulkData(data *BulkData, size int) *BulkData {
	data.len = 0
	for i := 0; i < size; i++ {
		data.datas[i] = ""
	}
	return data
}

func trimBulkData(bd *BulkData) *BulkData {
	for i := 0; i < bd.len; i++ {
		if bd.datas[i] != "" {

		}
	}
	return bd
}

func (p *BulkPool) Run() *BulkPool {
	for i := 0; i < p.PoolNum; i++ {
		go func(pindex int) {
			size := p.BulkSize
			bd := NewBulkData(make([]interface{}, size), 0)
			bd = resetBulkData(bd, size)
			run := true
			for run {
				select {
				case task,ok := <-p.BulkTask:
					if ok {
						bd.datas[bd.len] = task.data
						bd.f = task.f
						bd.len++
						//log.Printf("add data %d\n", bd.len)
						if bd.len == size {
							//log.Println("reach size call func")
							bd.f(bd.datas, pindex)
							bd = resetBulkData(bd, size)
						}
					} else {
						run = false
						//处理未完成的消息再结束
						if bd.len > 0 {
							//log.Println("clean datas before stop")
							subdatas := bd.datas[:bd.len]
							bd.f(subdatas, pindex)
							bd = resetBulkData(bd, size)
						}
					}
				case <-time.After(5 * time.Second):
					if bd.len > 0 {
						//log.Println("timeout call func")
						subdatas := bd.datas[:bd.len]
						bd.f(subdatas, pindex)
						bd = resetBulkData(bd, size)
					}
				}
			}
			p.Stoped <- true
		}(i)
	}
	return p
}

func (p *BulkPool) Close() {
	close(p.BulkTask)
	//等待停止信号
	for i := 0; i < p.PoolNum; i++ {
		<-p.Stoped
		//log.Printf("receive stoped [%d]\n", i)
	}
	//log.Println("close pool success")
}
