// Package mcp is the optional stdio MCP server (JOURNAL_SPEC §8.1):
// journal_search, journal_get(topic), journal_note(entry). Explicit
// journal_note writes are welcomed but nothing depends on them (I1).
package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"clew/internal/globx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/scrub"
)

type Server struct {
	Journal *journal.Journal
	Surface string
	// AfterWrite is called after journal_note writes (sync hook).
	AfterWrite func()
	clientName string
}

type req struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

func (s *Server) Serve(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)
	enc := json.NewEncoder(out)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var r req
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			continue
		}
		if r.ID == nil { // notification
			continue
		}
		resp := map[string]any{"jsonrpc": "2.0", "id": json.RawMessage(r.ID)}
		result, errObj := s.dispatch(&r)
		if errObj != nil {
			resp["error"] = errObj
		} else {
			resp["result"] = result
		}
		if err := enc.Encode(resp); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Server) dispatch(r *req) (any, map[string]any) {
	switch r.Method {
	case "initialize":
		var p struct {
			ProtocolVersion string `json:"protocolVersion"`
			ClientInfo      struct {
				Name string `json:"name"`
			} `json:"clientInfo"`
		}
		json.Unmarshal(r.Params, &p)
		s.clientName = p.ClientInfo.Name
		pv := p.ProtocolVersion
		if pv == "" {
			pv = "2025-06-18"
		}
		return map[string]any{
			"protocolVersion": pv,
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "clew", "version": "0.1.0"},
		}, nil
	case "ping":
		return map[string]any{}, nil
	case "tools/list":
		return map[string]any{"tools": toolDefs()}, nil
	case "tools/call":
		var p struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(r.Params, &p); err != nil {
			return nil, rpcErr(-32602, "bad params")
		}
		text, err := s.call(p.Name, p.Arguments)
		if err != nil {
			return map[string]any{
				"content": []map[string]any{{"type": "text", "text": "error: " + err.Error()}},
				"isError": true,
			}, nil
		}
		return map[string]any{"content": []map[string]any{{"type": "text", "text": text}}}, nil
	default:
		return nil, rpcErr(-32601, "method not found: "+r.Method)
	}
}

func toolDefs() []map[string]any {
	obj := func(props map[string]any, required ...string) map[string]any {
		return map[string]any{"type": "object", "properties": props, "required": required}
	}
	str := func(desc string) map[string]any { return map[string]any{"type": "string", "description": desc} }
	return []map[string]any{
		{
			"name":        "journal_search",
			"description": "Search the project journal (decisions, findings, questions, intents) by free text.",
			"inputSchema": obj(map[string]any{"query": str("free-text query")}, "query"),
		},
		{
			"name":        "journal_get",
			"description": "Get live journal entries for a topic (tag/path glob or title substring).",
			"inputSchema": obj(map[string]any{"topic": str("tag, path, or title substring")}, "topic"),
		},
		{
			"name":        "journal_note",
			"description": "Write a journal entry explicitly (optional; the watcher extracts automatically).",
			"inputSchema": obj(map[string]any{
				"type":  str("decision | finding | question | intent"),
				"title": str("≤80 chars"),
				"body":  str("≤400 chars, plain language"),
				"quote": str("the verbatim statement this entry records"),
			}, "type", "title", "quote"),
		},
	}
}

func (s *Server) call(name string, args json.RawMessage) (string, error) {
	switch name {
	case "journal_search":
		var a struct {
			Query string `json:"query"`
		}
		json.Unmarshal(args, &a)
		return s.search(a.Query), nil
	case "journal_get":
		var a struct {
			Topic string `json:"topic"`
		}
		json.Unmarshal(args, &a)
		return s.get(a.Topic), nil
	case "journal_note":
		var a struct{ Type, Title, Body, Quote string }
		if err := json.Unmarshal(args, &a); err != nil {
			return "", err
		}
		return s.note(a.Type, a.Title, a.Body, a.Quote)
	default:
		return "", fmt.Errorf("unknown tool %q", name)
	}
}

func (s *Server) render(idList []string) string {
	st := journal.Compute(s.Journal, time.Now())
	var b strings.Builder
	for _, id := range idList {
		e := s.Journal.Entries[id]
		c := st[id]
		fmt.Fprintf(&b, "%s [%s/%s] %s\n  quote: %q\n  source: %s %s\n",
			e.ID, e.Type, c.Status, e.Title, clip(e.Quote, 160), e.Source.Kind, e.Source.Ref)
		if e.Body != "" {
			fmt.Fprintf(&b, "  %s\n", clip(e.Body, 300))
		}
	}
	if b.Len() == 0 {
		return "no matching journal entries"
	}
	return b.String()
}

func (s *Server) search(q string) string {
	q = strings.ToLower(strings.TrimSpace(q))
	st := journal.Compute(s.Journal, time.Now())
	var hits []string
	for id, e := range s.Journal.Entries {
		if c := st[id]; c == nil || !journal.Live(c.Status) {
			continue
		}
		hay := strings.ToLower(e.Title + " " + e.Body + " " + e.Quote + " " + strings.Join(e.Tags, " "))
		if q == "" || strings.Contains(hay, q) {
			hits = append(hits, id)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(hits)))
	if len(hits) > 10 {
		hits = hits[:10]
	}
	return s.render(hits)
}

func (s *Server) get(topic string) string {
	topic = strings.TrimSpace(topic)
	st := journal.Compute(s.Journal, time.Now())
	var hits []string
	for id, e := range s.Journal.Entries {
		if c := st[id]; c == nil || !journal.Live(c.Status) {
			continue
		}
		match := false
		for _, t := range e.Tags {
			if strings.EqualFold(t, topic) || globx.Match(t, topic) || globx.Match(topic, t) {
				match = true
			}
		}
		if strings.Contains(strings.ToLower(e.Title), strings.ToLower(topic)) {
			match = true
		}
		if match {
			hits = append(hits, id)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(hits)))
	if len(hits) > 15 {
		hits = hits[:15]
	}
	return s.render(hits)
}

func (s *Server) note(typ, title, body, quote string) (string, error) {
	var et model.EntryType
	switch typ {
	case "decision":
		et = model.Decision
	case "finding":
		et = model.Finding
	case "question":
		et = model.Question
	case "intent":
		et = model.Intent
	default:
		return "", fmt.Errorf("type must be decision|finding|question|intent")
	}
	quote, _ = scrub.Scrub(quote)
	body, _ = scrub.Scrub(body)
	agent := s.clientName
	if agent == "" {
		agent = "mcp-client"
	}
	now := time.Now().UTC()
	e := &model.Entry{
		ID: ids.NewEntry(now), Type: et, Title: title, Body: body, Quote: quote,
		UtteranceBy: model.ByAssistant,
		Source: model.Source{Kind: model.SrcSession, Ref: "mcp:" + agent,
			Agent: agent, Surface: s.Surface, At: now},
		Confidence: 0.8,
	}
	if et == model.Question {
		e.Asks = "any"
	}
	if err := s.Journal.AddEntry(e); err != nil {
		return "", err
	}
	if s.AfterWrite != nil {
		s.AfterWrite()
	}
	return "recorded " + e.ID, nil
}

func rpcErr(code int, msg string) map[string]any {
	return map[string]any{"code": code, "message": msg}
}

func clip(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ServeStdio runs the server on the process's stdio.
func (s *Server) ServeStdio() error { return s.Serve(os.Stdin, os.Stdout) }
