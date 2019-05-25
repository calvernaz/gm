package log

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
)

const (
	bodySuccess = "%s\t \x1b[31;1m%s\x1b[0m"
	bodyError = "%s\t \x1b[31;1m%s\x1b[33;1m ( %s )\x1b[0m"
)

type BufferLog struct {
	mu     *sync.Mutex
	writer *tabwriter.Writer
}

func NewBufferLog() *BufferLog {
	w := tabwriter.NewWriter(os.Stdout, 20, 4, 1, '\t', tabwriter.TabIndent)
	return &BufferLog{
		mu:     &sync.Mutex{},
		writer: w,
	}
}

func (b *BufferLog) Buffer(args ...string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(args) == 2 {
		_, _ = fmt.Fprintln(b.writer, fmt.Sprintf(bodySuccess, args[0], args[1]))
	} else {
		_, _ = fmt.Fprintln(b.writer, fmt.Sprintf(bodyError, args[0], args[1], args[2]))
	}
}

func (b *BufferLog) Print() {
	b.mu.Lock()
	defer b.mu.Unlock()
	_ = b.writer.Flush()
}
