package config

import "io"

// sshConfigWriter records the first output error so the formatters can report
// it once and avoid retrying writes after a partial config has been emitted.
type sshConfigWriter struct {
	output io.Writer
	err    error
}

func (w *sshConfigWriter) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err := w.output.Write(p)
	if err == nil && n != len(p) {
		err = io.ErrShortWrite
	}
	w.err = err
	return n, err
}
