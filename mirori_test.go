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

	var right, left io.Reader

	testCases := [][]struct {
		r    *io.Reader
		ln   int
		want string
		last bool
	}{
		{
			{&left, 3, "foo", false},
			{&right, 3, "foo", false},
			{&left, 3, "bar", false},
			{&left, 3, "baz", false},
			{&left, 3, "", true},
			{&right, 3, "bar", false},
			{&right, 3, "baz", false},
			{&right, 3, "", true},
		},
		{
			{&left, 3, "foo", false},
			{&left, 3, "bar", false},
			{&left, 3, "baz", false},
			{&left, 3, "", true},
			{&right, 3, "foo", false},
			{&right, 3, "bar", false},
			{&right, 4, "baz", true},
			{&right, 3, "", true},
		},
		{
			{&left, 3, "foo", false},
			{&right, 6, "foobar", false},
			{&left, 10, "barbaz", false},
			{&left, 10, "", true},
			{&right, 10, "baz", true},
			{&right, 3, "", true},
			{&left, 3, "", true},
		},
	}

	for i, run := range testCases {
		left, right = miroir.NewMiroir(strings.NewReader(in))
		for j, tc := range run {
			t.Run(fmt.Sprint(i, j), func(t *testing.T) {
				buf := make([]byte, tc.ln)
				n, err := (*tc.r).Read(buf)
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
}
