package main

import (
	"bytes"
	o "cli/todo/output"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIsTerminal(t *testing.T) {
	output := o.CreateOutput()
	stat, _ := os.Stdout.Stat()
	expected := (stat.Mode() & os.ModeCharDevice) != 0

	if output.IsTerminal() != expected {
		t.Errorf("Expected %v got %v", expected, output.IsTerminal())
	}
}

func captureOutput(f func()) string {
	out := os.Stdout
	defer func() { os.Stdout = out }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()

}

func TestFindTodos(t *testing.T) {
	todos := []string{
		ColorGreen.Apply("5:testdir/another.php - TODO: Parse env arguments before running dirFindTodos"),
		ColorGreen.Apply("3:testdir/internal/find.php - TODO: Find todos in a directory"),
		ColorGreen.Apply("3:testdir/test.php - TODO: finish this class"),
	}

	res := captureOutput(func() {
		o := o.CreateOutput()
		o.Terminal = true
		FindTodos("testdir", &o)
	})
	resSplit := strings.Split(res, "\n")

	filtered := []string{}

	for _, str := range resSplit {
		if str != "" {
			filtered = append(filtered, str)
		}
	}
	if !reflect.DeepEqual(todos, filtered) {
		t.Errorf("Expected %s got %s", todos, filtered)
	}
}
