package miroir_test

import (
	"io"
	"strings"
	"testing"

	"github.com/mkmik/miroir"
)

func TestReadAll(t *testing.T) {
	in := "foobarbaz"
	left, right := miroir.NewMiroir(strings.NewReader(in))

	lb, err := io.ReadAll(left)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(lb), "foobarbaz"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	rb, err := io.ReadAll(right)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(rb), "foobarbaz"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}

func TestInterleaved(t *testing.T) {
	in := "foobarbaz"
	left, right := miroir.NewMiroir(strings.NewReader(in))

	buf := make([]byte, 3)

	if _, err := left.Read(buf); err != nil {
		t.Fatal(err)
	}
	if got, want := string(buf), "foo"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	if _, err := left.Read(buf); err != nil {
		t.Fatal(err)
	}
	if got, want := string(buf), "bar"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	if _, err := right.Read(buf); err != nil {
		t.Fatal(err)
	}
	if got, want := string(buf), "foo"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	if _, err := left.Read(buf); err != nil {
		t.Fatal(err)
	}
	if got, want := string(buf), "baz"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	if _, err := right.Read(buf); err != nil {
		t.Fatal(err)
	}
	if got, want := string(buf), "bar"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	if _, err := right.Read(buf); err != nil {
		t.Fatal(err)
	}
	if got, want := string(buf), "bar"; got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}
