package main

import (
	"sync"

	"github.com/lzhy87/qqcmall-monitor/utils"
)
var wg sync.WaitGroup
func main() {

	wg.Add(1)
	go utils.MonitorPort()
	go utils.CheckService()

	wg.Wait()

}
