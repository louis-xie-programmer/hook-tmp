package hook

import (
	"log"
	"time"
)

func recordCost(name string, start time.Time) {
	log.Printf("[METRIC] hook=%s cost=%v", name, time.Since(start))
}

func recordError(name string, err error) {
	log.Printf("[ERROR] hook=%s err=%v", name, err)
}

func recordPanic(name string, r any) {
	log.Printf("[PANIC] hook=%s panic=%v", name, r)
}
