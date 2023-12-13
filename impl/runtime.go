package impl

import (
	"context"
	"encoding/json"
	"main/core"
	"os"
	"strings"
)

type RuntimeData struct {
	Agents     [2]core.Agent
	Transcript []core.TranscriptEntry
}

type Runtime struct {
	RuntimeData
	Log        core.Log
	LLM        core.LLM
	File       string
	Iterations int
}

func (r *Runtime) prepareTranscript(agentIndex int) string {
	bldr := new(strings.Builder)
	for _, t := range r.Transcript {
		if t.Agent == agentIndex {
			bldr.WriteString("You: ")
		} else {
			bldr.WriteString("Agent: ")
		}
		bldr.WriteString(t.Text)
		bldr.WriteString("\n\n")
	}
	return bldr.String()
}

func (r *Runtime) Start(ctx context.Context) error {
	r.Transcript = make([]core.TranscriptEntry, 0)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			idx := r.Iterations % len(r.Agents)
			r.Log.Infof("Conversation step %d (agent %d)", r.Iterations, idx)
			currentAgent := r.Agents[idx]
			prompt := new(strings.Builder)

			if len(r.Transcript) > 0 {
				prompt.WriteString("Given the chat transcript below and given it is your turn to speak in the conversation, write the next line of the conversation that drives it towards your intended goal. Respond only with text that should belong in the transcript. Do not add a speaker label to the text. If the conversation has resolved itself to a satisfactory conclusion that satisfies your goal, then respond with \"DONE\".\n\n")
				prompt.WriteString(r.prepareTranscript(idx))
			} else {
				prompt.WriteString("Write the first line of the conversation that drives it towards your intended goal. Respond only with text that should belong in the transcript. Do not add a speaker label to the text.")
			}

			prmpt := prompt.String()

			res, err := r.LLM.Completion(ctx, currentAgent.Intention, prmpt)
			if err != nil {
				return err
			}
			col := strings.Index(res, ":")
			if col > 0 {
				res = res[col+1:]
			}
			r.Transcript = append(r.Transcript, core.TranscriptEntry{
				Agent: idx,
				Text:  res,
			})
			tsx, err := json.Marshal(r.RuntimeData)
			if err != nil {
				return err
			}
			err = os.WriteFile(r.File, tsx, 0777)
			if err != nil {
				return err
			}
			r.Iterations++
		}
	}
}
