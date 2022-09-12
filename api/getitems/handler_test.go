package getitems

import "testing"

func TestHandler(t *testing.T) {
  got := "foo"
  want := "bar"

  if got != want {
    t.Errorf("got %q want %q", got, want)
  }
}
