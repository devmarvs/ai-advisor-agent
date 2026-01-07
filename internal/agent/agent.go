package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Model        string
	SystemPrompt string
	MaxTurns     int
	Tools        Toolset
}

func DefaultSystemPrompt() string {
	return strings.TrimSpace(`You are a helpful AI assistant for a financial advisor.
You can answer questions about the user's clients using email and calendar data.
When you need data, you can call tools by returning JSON: {"tool":"name","args":{...}}.
Tools available: search_context, gmail_send, calendar_find_slots, calendar_create_event.
If no tool is needed, just answer.`)
}

func New(cfg Config) *Agent {
	if cfg.MaxTurns == 0 { cfg.MaxTurns = 4 }
	if cfg.Tools == nil { cfg.Tools = DefaultToolset{} }
	return &Agent{cfg: cfg, llm: NewLLM()}
}

type Agent struct {
	cfg Config
	llm *LLM
}

type toolCall struct {
	Tool string                 `json:"tool"`
	Args map[string]interface{} `json:"args"`
}

func (a *Agent) Handle(ctx context.Context, userID string, message string) (string, string, error) {
	turns := 0
	trace := &strings.Builder{}
	for {
		turns++
		if turns > a.cfg.MaxTurns {
			return "I reached the maximum steps. If you want me to continue, please ask again.", trace.String(), nil
		}
		reply, err := a.llm.Complete(ctx, a.cfg.SystemPrompt, message+`\n\nIf calling a tool, respond with only a JSON object: {"tool":"...","args":{...}}. Otherwise, reply normally.`)
		if err != nil { return "", trace.String(), err }

		var call toolCall
		if json.Unmarshal([]byte(reply), &call) == nil && call.Tool != "" {
			out, err := a.execTool(ctx, userID, call)
			if err != nil { return "", trace.String(), err }
			fmt.Fprintf(trace, "→ tool:%s args:%v\n← %s\n", call.Tool, call.Args, out)
			message = fmt.Sprintf("Tool result:\n%s\n\nPlease continue.", out)
			continue
		}
		return reply, trace.String(), nil
	}
}

func (a *Agent) execTool(ctx context.Context, userID string, call toolCall) (string, error) {
	switch call.Tool {
	case "search_context":
		q, _ := call.Args["query"].(string)
		limit := 6
		if v, ok := call.Args["limit"].(float64); ok && v > 0 { limit = int(v) }
		docs, err := a.cfg.Tools.SearchContext(ctx, userID, q, limit)
		if err != nil { return "", err }
		b, _ := json.MarshalIndent(docs, "", "  ")
		return string(b), nil
	case "gmail_send":
		to, _ := call.Args["to"].(string)
		subject, _ := call.Args["subject"].(string)
		text, _ := call.Args["text"].(string)
		err := a.cfg.Tools.SendEmail(ctx, userID, to, subject, text); return "sent", err
	case "calendar_find_slots":
		now := time.Now()
		slots, err := a.cfg.Tools.FindSlots(ctx, userID, now, now.AddDate(0,0,7), nil)
		if err != nil { return "", err }
		b, _ := json.MarshalIndent(slots, "", "  ")
		return string(b), nil
	case "calendar_create_event":
		title, _ := call.Args["title"].(string)
		whenStr, _ := call.Args["when"].(string)
		desc, _ := call.Args["description"].(string)
		when, _ := time.Parse(time.RFC3339, whenStr)
		id, err := a.cfg.Tools.CreateEvent(ctx, userID, title, when, nil, desc)
		if err != nil { return "", err }
		return id, nil
	}
	return "", fmt.Errorf("unknown tool %q", call.Tool)
}
