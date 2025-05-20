package pkg

import (
	"runtime"
	"server-template/pkg/log"
)

func Recover(clean func()) {
	if err := recover(); err != nil {
		buf := make([]byte, 1024)
		for {
			n := runtime.Stack(buf, false)
			if n < len(buf) {
				buf = buf[:n]
				break
			}
			buf = make([]byte, 2*len(buf))
		}
		log.Errorf("stacktrace from panic, error: %+v\nstack: %s", err, string(buf))
		if clean != nil {
			clean()
		}
	}
}
