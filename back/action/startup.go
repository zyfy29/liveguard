package action

import (
	"bearguard/rest"
	"runtime"
)

func Startup() {
	if runtime.GOOS == "linux" {
		go WatchTaskAndDownloadLive()
		go WatchTaskAndCallLLM()
		go WatchAndSubmitToMedium()
		go WatchAndRetrySubmit()
	}

	go rest.Startup()
}
