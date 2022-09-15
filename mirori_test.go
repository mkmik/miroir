package miroir_test

import (
	"fmt"
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
	_ = right
	testCases := []struct {
		r    io.Reader
		want string
		last bool
	}{
		{left, "foo", false},
		{right, "foo", false},
		{left, "bar", false},
		{left, "baz", false},
		{left, "", true},
		{right, "bar", false},
		{right, "baz", false},
		{left, "", true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			buf := make([]byte, len(tc.want))
			n, err := tc.r.Read(buf)
			eof := err == io.EOF
			if err != nil && !eof {
				t.Fatal(err)
			}
			if got, want := eof, tc.last; got != want {
				t.Fatalf("got: %v, want: %v", got, want)
			}
			if got, want := n, len(tc.want); got != want {
				t.Fatalf("got: %d, want: %d", got, want)
			}
			if got, want := string(buf[:n]), tc.want; got != want {
				t.Fatalf("got: %q, want: %q", got, want)
			}
		})
	}
}
