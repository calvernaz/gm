package log

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
)

const (
	body = "%s\t \x1b[31;1m%s\x1b[0m"
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
	_, _ = fmt.Fprintln(b.writer, fmt.Sprintf(body, args[0], args[1]))
}

func (b *BufferLog) Print() {
	b.mu.Lock()
	defer b.mu.Unlock()
	_ = b.writer.Flush()
}
