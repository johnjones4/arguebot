package core

import "context"

type LLM interface {
	Completion(ctx context.Context, role string, prompt string) (string, error)
	Embedding(ctx context.Context, text string) ([]float32, error)
}

type Log interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
}
