package bot

import "testing"

func TestAddMdCode(t *testing.T) {
	t.Parallel()
	want := "```text\ntest\n```"
	got := addMdCode("test")
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v", got, want)
	}
}
