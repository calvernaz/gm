package log

import (
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
)

func TestBufferLog(t *testing.T) {
	w := tabwriter.NewWriter(os.Stdout, 20, 4, 1, '\t', tabwriter.TabIndent)
	_, _ = fmt.Fprintln(w, fmt.Sprintf(body, "repository", "updated"))
	_, _ = fmt.Fprintln(w, fmt.Sprintf(body, "repositoryTest", "not updated"))
	_ = w.Flush()
}
