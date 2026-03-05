package mistral

import (
	"encoding/json"
	"fmt"
	"io"
)

// Stream is a generic iterator for streaming API responses.
// Use Next() to advance, Current() to read the value, Err() for errors,
// and Close() when done.
type Stream[T any] struct {
	reader  *sseReader
	closer  io.Closer
	current T
	err     error
	done    bool
}

func newStream[T any](body io.ReadCloser) *Stream[T] {
	return &Stream[T]{
		reader: newSSEReader(body),
		closer: body,
	}
}

// Next advances to the next event. Returns false when the stream
// is exhausted or an error occurs.
func (s *Stream[T]) Next() bool {
	if s.done || s.err != nil {
		return false
	}
	for {
		event, err := s.reader.next()
		if err != nil {
			s.err = fmt.Errorf("mistral: read stream: %w", err)
			return false
		}
		if event == nil {
			s.done = true
			return false
		}
		if event.isDone() {
			s.done = true
			return false
		}

		var v T
		if err := json.Unmarshal(event.Data, &v); err != nil {
			s.err = fmt.Errorf("mistral: decode stream event: %w", err)
			return false
		}
		s.current = v
		return true
	}
}

// Current returns the most recently read value.
// Only valid after Next() returns true.
func (s *Stream[T]) Current() T {
	return s.current
}

// Err returns any error encountered during streaming.
func (s *Stream[T]) Err() error {
	return s.err
}

// Close releases the underlying HTTP response body.
func (s *Stream[T]) Close() error {
	return s.closer.Close()
}
