package cli

import (
	"fmt"
	"io"
	"os"
)

// progressWriter wraps an io.Writer and prints byte progress to stderr when verbose.
type progressWriter struct {
	w       io.Writer
	label   string
	written int64
	lastPct int
	total   int64 // 0 = unknown
}

func newProgressWriter(w io.Writer, label string, total int64) *progressWriter {
	return &progressWriter{w: w, label: label, total: total, lastPct: -1}
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n, err := p.w.Write(b)
	p.written += int64(n)
	if flagVerbose && !flagQuiet && !flagJSON {
		if p.total > 0 {
			pct := int(p.written * 100 / p.total)
			if pct != p.lastPct && (pct%5 == 0 || pct == 100) {
				fmt.Fprintf(os.Stderr, "%s: %d%% (%d / %d bytes)\n", p.label, pct, p.written, p.total)
				p.lastPct = pct
			}
		} else if p.written == int64(n) || p.written%(1<<20) < int64(n) {
			fmt.Fprintf(os.Stderr, "%s: %d bytes written\n", p.label, p.written)
		}
	}
	return n, err
}
