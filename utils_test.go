package main

import (
	"testing"
)

func TestAddMdCode(t *testing.T) {
	t.Parallel()

	s := addMdCode("test")
	want := "```text\ntest\n```"

	if s != want {
		t.Errorf("got:\n%v\nwant:\n%v", s, want)
	}
}
