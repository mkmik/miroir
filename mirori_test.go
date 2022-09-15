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

	testCases := []struct {
		r    io.Reader
		want string
	}{
		{left, "foo"},
		{right, "foo"},
		{left, "bar"},
		{left, "baz"},
		{right, "bar"},
		{right, "baz"},
	}

	buf := make([]byte, 3)
	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			n, err := tc.r.Read(buf)
			if err != nil && err != io.EOF {
				t.Fatal(err)
			}
			if got, want := n, len(tc.want); got != want {
				t.Fatalf("got: %d, want: %d", got, want)
			}
			if got, want := string(buf), tc.want; got != want {
				t.Fatalf("got: %q, want: %q", got, want)
			}
		})
	}
}
