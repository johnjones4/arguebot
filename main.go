package main

import (
	"context"
	"flag"
	"fmt"
	"main/core"
	"main/impl"
	"main/impl/llm"
	"path"
	"time"

	"go.uber.org/zap"
)

func main() {
	ollamaUrl := flag.String("ollama-url", "http://localhost:11434", "URL root of the desired Ollama instance.")
	ollamaModel := flag.String("ollama-model", "mistral", "Desired Ollama model.")
	logLevel := flag.String("log-level", "info", "Desired log level")
	outputDir := flag.String("output-dir", "./data/", "Directory to output conversations")
	runtime := flag.Duration("duration", time.Hour, "Time limit for conversation")
	agent1 := flag.String("agent-1", "", "Instruction prompt for first agent")
	agent2 := flag.String("agent-2", "", "Instruction prompt for second agent")

	flag.Parse()

	config := zap.NewDevelopmentConfig()
	level, err := zap.ParseAtomicLevel(*logLevel)
	if err != nil {
		flag.Usage()
		panic(err)
	}
	config.Level = level
	l, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer l.Sync()
	log := l.Sugar()

	rt := impl.Runtime{
		Log: log,
		LLM: &llm.Ollama{
			Model: *ollamaModel,
			URL:   *ollamaUrl,
			Log:   log,
		},
		File: path.Join(*outputDir, fmt.Sprintf("conv_%s.json", time.Now().String())),
		RuntimeData: impl.RuntimeData{
			Agents: [2]core.Agent{
				{
					Intention: *agent1,
				},
				{
					Intention: *agent2,
				},
			},
		},
	}

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(*runtime))
	err = rt.Start(ctx)
	if err != nil {
		panic(err)
	}
}
