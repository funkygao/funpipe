package main

import (
	"github.com/funkygao/dpipe/engine"
)

var (
	globals *engine.GlobalConfigStruct

	BuildID = "unknown" // git version id, passed in from shell

	options struct {
		verbose            bool
		veryVerbose        bool
		configfile         string
		showversion        bool
		logfile            string
		debug              bool
		tick               int
		dryrun             bool
		cpuprof            string
		memprof            string
		lockfile           string
		diagnosticInterval int
	}
)

const (
	USAGE = `dpiped - Distributed Data Pipeline

Flags:
`
)
