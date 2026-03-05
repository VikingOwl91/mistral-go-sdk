package mistral

import (
	"strings"
	"testing"
)

func TestSSEReader_SingleEvent(t *testing.T) {
	input := "data: {\"id\":\"1\"}\n\n"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected event")
	}
	if string(ev.Data) != `{"id":"1"}` {
		t.Errorf("got data %q", ev.Data)
	}
}

func TestSSEReader_MultipleEvents(t *testing.T) {
	input := "data: first\n\ndata: second\n\n"
	r := newSSEReader(strings.NewReader(input))

	ev1, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if string(ev1.Data) != "first" {
		t.Errorf("got %q, want %q", ev1.Data, "first")
	}

	ev2, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if string(ev2.Data) != "second" {
		t.Errorf("got %q, want %q", ev2.Data, "second")
	}

	ev3, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev3 != nil {
		t.Errorf("expected nil at EOF, got %+v", ev3)
	}
}

func TestSSEReader_MultiLineData(t *testing.T) {
	input := "data: line1\ndata: line2\ndata: line3\n\n"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	want := "line1\nline2\nline3"
	if string(ev.Data) != want {
		t.Errorf("got %q, want %q", ev.Data, want)
	}
}

func TestSSEReader_EventField(t *testing.T) {
	input := "event: completion\ndata: {\"id\":\"1\"}\n\n"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev.Event != "completion" {
		t.Errorf("got event %q, want %q", ev.Event, "completion")
	}
	if string(ev.Data) != `{"id":"1"}` {
		t.Errorf("got data %q", ev.Data)
	}
}

func TestSSEReader_SkipsComments(t *testing.T) {
	input := ": this is a comment\ndata: hello\n\n"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if string(ev.Data) != "hello" {
		t.Errorf("got %q, want %q", ev.Data, "hello")
	}
}

func TestSSEReader_Done(t *testing.T) {
	input := "data: {\"id\":\"1\"}\n\ndata: [DONE]\n\n"
	r := newSSEReader(strings.NewReader(input))

	ev1, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev1.isDone() {
		t.Error("first event should not be done")
	}

	ev2, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if !ev2.isDone() {
		t.Error("second event should be done")
	}
}

func TestSSEReader_DoneWithWhitespace(t *testing.T) {
	ev := &sseEvent{Data: []byte(" [DONE] ")}
	if !ev.isDone() {
		t.Error("should detect [DONE] with whitespace")
	}
}

func TestSSEReader_EmptyStream(t *testing.T) {
	r := newSSEReader(strings.NewReader(""))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev != nil {
		t.Errorf("expected nil for empty stream, got %+v", ev)
	}
}

func TestSSEReader_OnlyComments(t *testing.T) {
	input := ": comment1\n: comment2\n\n"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev != nil {
		t.Errorf("expected nil, got %+v", ev)
	}
}

func TestSSEReader_NoTrailingNewline(t *testing.T) {
	input := "data: hello"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if ev == nil {
		t.Fatal("expected event for data without trailing blank line")
	}
	if string(ev.Data) != "hello" {
		t.Errorf("got %q, want %q", ev.Data, "hello")
	}
}

func TestSSEReader_DataNoSpace(t *testing.T) {
	input := "data:{\"compact\":true}\n\n"
	r := newSSEReader(strings.NewReader(input))
	ev, err := r.next()
	if err != nil {
		t.Fatal(err)
	}
	if string(ev.Data) != `{"compact":true}` {
		t.Errorf("got %q", ev.Data)
	}
}

func TestSSEReader_MultipleBlankLines(t *testing.T) {
	input := "data: first\n\n\n\ndata: second\n\n"
	r := newSSEReader(strings.NewReader(input))

	ev1, _ := r.next()
	if string(ev1.Data) != "first" {
		t.Errorf("got %q", ev1.Data)
	}
	ev2, _ := r.next()
	if string(ev2.Data) != "second" {
		t.Errorf("got %q", ev2.Data)
	}
}
