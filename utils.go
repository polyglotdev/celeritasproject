package celeritas

import (
	"regexp"
	"runtime"
	"time"
)

// LoadTime logs the time it took to load a page
func (c *Celeritas) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc).Name()
	runTimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runTimeFunc.ReplaceAllString(details, "$1")
	c.InfoLog.Println("This page took", elapsed, "to load", name)
}
