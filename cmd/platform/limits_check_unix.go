// +build !windows

package main

import (
	"syscall"

	l4g "github.com/alecthomas/log4go"

	"github.com/mattermost/platform/utils"
)

func limitsCheck() {
	var rlim syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlim)
	if err != nil {
		l4g.Error(utils.T("mattermost.rlimit.error"), err.Error())
	} else {
		if rlim.Cur < 50000 {
			l4g.Error(utils.T("mattermost.rlimit_low.error"), rlim.Cur)
		} else {
			l4g.Info(utils.T("mattermost.rlimit.info"), rlim.Cur)
		}
	}
}
