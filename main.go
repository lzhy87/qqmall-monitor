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

	wg.Wait()

}
