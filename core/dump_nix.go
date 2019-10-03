// +build linux darwin

package core

import (
	"os"
	"os/signal"
	"syscall"
)

func DumpSessionOnSig(sess *Session) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)
	for range sigChan {
		sess.Lock()
		err := sess.SaveToFile("aquatone_session.json")
		if err != nil {
			sess.Out.Error("Failed to write session file")
			sess.Out.Debug("Err: %s", err.Error())
		}
		sess.Unlock()
	}
}