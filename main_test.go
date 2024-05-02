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
		ColorGreen.Apply("6:testdir/main.js - TODO: implementation missing"),
		ColorGreen.Apply("3:testdir/test.php - TODO: finish this class"),
		ColorGreen.Apply("7:testdir/test.php - FIXME: Add variable types"),
	}

	res := captureOutput(func() {
		o := o.CreateOutput()
		o.Terminal = true
		sc := SearchContext{
			output:            &o,
			keywords:          []string{"FIXME:", "TODO:"},
			exclude:           []string{".git", ".svn"},
			include_filetypes: []string{},
		}
		FindTodos("testdir", &sc)
	})
	resSplit := strings.Split(res, "\n")

	// Filter out empty
	filtered := filter(resSplit, func(a string) bool { return a != "" })

	if !reflect.DeepEqual(todos, filtered) {
		t.Errorf("Expected %s got %s", todos, filtered)
	}
}
func filter[T any](arr []T, f func(p T) bool) []T {
	var result []T
	for _, el := range arr {
		if f(el) {
			result = append(result, el)
		}
	}
	return result
}

func TestExcludeDir(t *testing.T) {
	todosSkipped := []string{
		ColorGreen.Apply("5:testdir/another.php - TODO: Parse env arguments before running dirFindTodos"),
		ColorRed.Apply("Skipping ignored dir - internal"),
		ColorGreen.Apply("6:testdir/main.js - TODO: implementation missing"),
		ColorGreen.Apply("3:testdir/test.php - TODO: finish this class"),
		ColorGreen.Apply("7:testdir/test.php - FIXME: Add variable types"),
	}

	res := captureOutput(func() {
		o := o.CreateOutput()
		o.Terminal = true
		sc := SearchContext{
			output:            &o,
			keywords:          []string{"FIXME:", "TODO:"},
			exclude:           []string{".git", ".svn", "internal"},
			include_filetypes: []string{},
		}
		FindTodos("testdir", &sc)
	})
	resSplit := strings.Split(res, "\n")

	// Filter out empty
	filtered := filter(resSplit, func(a string) bool { return a != "" })

	if !reflect.DeepEqual(todosSkipped, filtered) {
		t.Errorf("Expected %s got %s", todosSkipped, filtered)
	}

}

func TestOnlyFiletype(t *testing.T) {
	todosSkipped := []string{
		ColorGreen.Apply("6:testdir/main.js - TODO: implementation missing"),
	}

	res := captureOutput(func() {
		o := o.CreateOutput()
		o.Terminal = true
		sc := SearchContext{
			output:            &o,
			keywords:          []string{"FIXME:", "TODO:"},
			exclude:           []string{".git", ".svn", "internal"},
			include_filetypes: []string{"js"},
		}
		FindTodos("testdir", &sc)
	})
	resSplit := strings.Split(res, "\n")

	// Filter out empty
	filtered := filter(resSplit, func(a string) bool { return a != "" })

	if !reflect.DeepEqual(todosSkipped, filtered) {
		t.Errorf("Expected %s got %s", todosSkipped, filtered)
	}
}
