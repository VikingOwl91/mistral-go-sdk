package mistral

import (
	"bufio"
	"bytes"
	"io"
)

// sseEvent represents a single Server-Sent Event.
type sseEvent struct {
	Event string
	Data  []byte
}

// isDone returns true if this event signals end-of-stream.
func (e *sseEvent) isDone() bool {
	return string(bytes.TrimSpace(e.Data)) == "[DONE]"
}

// sseReader reads Server-Sent Events from an io.Reader.
type sseReader struct {
	scanner *bufio.Scanner
}

func newSSEReader(r io.Reader) *sseReader {
	return &sseReader{scanner: bufio.NewScanner(r)}
}

// next reads the next SSE event. Returns nil, nil at EOF.
func (r *sseReader) next() (*sseEvent, error) {
	var event sseEvent
	var hasData bool

	for r.scanner.Scan() {
		line := r.scanner.Bytes()

		// Blank line = end of event
		if len(line) == 0 {
			if hasData {
				return &event, nil
			}
			continue
		}

		// Skip comments
		if line[0] == ':' {
			continue
		}

		field, value, _ := bytes.Cut(line, []byte(":"))
		// Strip single leading space from value per SSE spec
		value = bytes.TrimPrefix(value, []byte(" "))

		switch string(field) {
		case "event":
			event.Event = string(value)
		case "data":
			if hasData {
				event.Data = append(event.Data, '\n')
			}
			event.Data = append(event.Data, value...)
			hasData = true
		}
	}

	if err := r.scanner.Err(); err != nil {
		return nil, err
	}

	// Final event without trailing blank line
	if hasData {
		return &event, nil
	}
	return nil, nil
}
