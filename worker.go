package main

import (
	"github.com/ActiveState/tail"
	parser "github.com/funkygao/alsparser"
	"sync"
)

// Each single log file is a worker
// Workers share some singleton parsers
func run_worker(logfile string, conf jsonItem, wg *sync.WaitGroup, chLines chan int) {
	defer wg.Done()

	var tailConfig tail.Config
	if options.tailmode {
		tailConfig = tail.Config{
			Follow: true,
			ReOpen: true,
		}
	}

	t, err := tail.TailFile(logfile, tailConfig)
	if err != nil {
		panic(err)
	}

	defer t.Stop()

	for line := range t.Lines {
		// a valid line scanned
		chLines <- 1

		for _, p := range conf.Parsers {
			parser.Dispatch(p, line.Text)
		}
	}

	if options.verbose {
		logger.Println(logfile, "finished")
	}
}
