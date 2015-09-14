package gostatic

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// Preprocessor is a processor which will be executed during initialization
	// stage
	Pre = 1 << iota
	Hidden
	Post
)

type Processor interface {
	Process(page *Page, args []string) error
	Description() string
	Mode() int
}

// it is necessary to wrap assignment in a function since go compiler is strict
// enough to throw error when there may be an assignment loop, but not smart
// enough to determine if there is one (Processors definition loops with
// ProcessTags)
func (s *Site) InitProcessors(m map[string]Processor) {
	s.Processors = m
}

func (s *Site) ProcessorSummary() {
	keys := make([]string, 0, len(s.Processors))
	for k := range s.Processors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		p := s.Processors[k]
		if p.Mode()&Hidden != 0 {
			continue
		}
		pre := ""
		if p.Mode()&Pre != 0 {
			pre = "(preprocessor)"
		}
		fmt.Printf("- %s %s\n\t%s\n", k, pre, p.Description())
	}
}

func (s *Site) ProcessCommand(page *Page, cmd *Command, pre bool) error {
	c := string(*cmd)
	if strings.HasPrefix(c, ":") {
		c = "external " + c[1:]
	}
	bits := strings.Split(c, " ")

	processor := s.Processors[bits[0]]
	if processor == nil {
		return fmt.Errorf("Processor is nil")
	}
	if (processor.Mode()&Pre != 0) != pre {
		return nil
	}
	if processor == nil {
		return fmt.Errorf("processor '%s' not found", bits[0])
	}
	return processor.Process(page, bits[1:])
}
