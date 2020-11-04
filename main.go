package main

import (
	"sync"

	"v1/utils"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	// go utils.MonitorPort()
	go utils.CheckService()
	// utils.DingToInfo("服务出现问题，请及时处理")

	wg.Wait()

}
