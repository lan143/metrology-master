package cmdlog

import (
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func Signal(key string, sig os.Signal) zap.Field {
	return zap.String(key, SignalName(sig))
}

func Signals(key string, sigs []os.Signal) zap.Field {
	names := make([]string, 0, len(sigs))

	for i := range sigs {
		names = append(
			names,
			SignalName(sigs[i]),
		)
	}

	return zap.Strings(key, names)
}

func SignalName(sig os.Signal) string {
	var name string
	if x, ok := sig.(syscall.Signal); ok {
		name = unix.SignalName(x)
	}

	if name == "" {
		name = sig.String()
	}

	return name
}
