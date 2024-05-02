package main

import (
	"bufio"
	o "cli/todo/output"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

type Color string

const (
	ColorRed   Color = "\033[31m"
	ColorGreen Color = "\033[32m"
	ColorBlue  Color = "\033[34m"
	ColorReset Color = "\033[0m"
)

// var keywords = []string{
// 	"TODO:",
// 	"FIXME:",
// }

func UpdateProgress(filename string) {
	fmt.Printf("\rProcessing file: %s", filename)
}

func (c Color) Apply(str string) string {
	return string(c) + str + string(ColorReset)
}

func FindTodos(dir string, sctx *SearchContext) {

	files, err := os.ReadDir(dir)

	if err != nil {
		sctx.output.ExitWithError(err.Error(), 2)
	}
	var wg sync.WaitGroup

	for _, entry := range files {

		filename := entry.Name()
		fp := fmt.Sprintf("%s/%s", dir, filename)

		if entry.IsDir() {
			if slices.Contains(sctx.exclude, filename) && sctx.output.IsTerminal() {
				if sctx.output.IsTerminal() {
					sctx.output.Printf(ColorRed.Apply("Skipping ignored dir - %s")+"\n", filename)
				}
				continue
			}
			wg.Add(1) // Increment group counter

			go func(path string) {
				defer wg.Done()
				FindTodos(path, sctx)
			}(fp)
		}

		extension := strings.Trim(filepath.Ext(fp), ".")

		if len(sctx.include_filetypes) > 0 && !slices.Contains(sctx.include_filetypes, extension) {
			continue
		}

		f, err := os.Open(fp)

		if err != nil {
			sctx.output.ExitWithError(err.Error(), 2)
		}
		scanner := bufio.NewScanner(f)

		lineno := 1
		for scanner.Scan() {
			line := scanner.Text()

			for _, keyword := range sctx.keywords {
				if strings.Contains(line, keyword) {
					line = strings.Replace(line, "//", "", -1)
					line = strings.Trim(line, " ")
					rawLine := fmt.Sprintf("%d:%s - %s", lineno, fp, line)

					line = ColorGreen.Apply(rawLine)

					if sctx.output.IsTerminal() {
						sctx.output.Print(line)
					} else {
						sctx.output.Printf("- [ ] %s\n", rawLine)
					}
				}

			}

			lineno++
		}

		f.Close()
		wg.Wait()
	}
}

type SearchContext struct {
	output            *o.Output
	keywords          []string
	include_filetypes []string
	exclude           []string
}

func main() {

	output := o.CreateOutput()
	sc := SearchContext{
		output:            &output,
		keywords:          []string{"FIXME:", "TODO:"},
		exclude:           []string{".git", ".svn"},
		include_filetypes: []string{},
	}

	exclude := flag.String("exclude-dir", "", "Exclude a directory")
	filetype := flag.String("filetype", "", "Search only inside filetype")

	flag.Parse()

	if len(os.Args) <= 1 {
		output.ExitWithError("USAGE:  program <folder_name>", 1)
	}

	dir := flag.Args()[0]

	if *exclude != "" {
		sc.exclude = append(sc.exclude, *exclude)
	}

	if *filetype != "" {
		sc.include_filetypes = append(sc.include_filetypes, *filetype)
	}

	stat, err := os.Stat(dir)

	if err != nil {
		output.ExitWithError(err.Error(), 1)
	}

	if !stat.IsDir() {
		output.ExitWithError(fmt.Sprintf("'%s' is not a directory\n", dir), 1)
	}
	FindTodos(dir, &sc)
}
