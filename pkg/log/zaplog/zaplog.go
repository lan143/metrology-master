package zaplog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ParseLevel(level string) (zapcore.LevelEnabler, error) {
	enabler := zap.NewAtomicLevel()

	switch level {
	case "debug":
		enabler.SetLevel(zapcore.DebugLevel)
	case "info":
		enabler.SetLevel(zapcore.InfoLevel)
	case "warn":
		enabler.SetLevel(zapcore.WarnLevel)
	case "error":
		enabler.SetLevel(zapcore.ErrorLevel)
	case "fatal":
		enabler.SetLevel(zapcore.FatalLevel)
	default:
		return nil, fmt.Errorf("unknown log level %q", level)
	}

	return enabler, nil
}

func ParseOutput(output string) (zapcore.WriteSyncer, error) {
	switch output {
	case "stdout":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		return nil, fmt.Errorf("unknown log output %q", output)
	}
}
