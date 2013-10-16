package task

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type taskParser struct {
	Dir string
}

func (l *taskParser) Parse() (taskSet *TaskSet, err error) {
	dir, err := expandDir(l.Dir)
	if err != nil {
		return
	}

	p, e := build.ImportDir(dir, 0)
	taskFiles := append(p.GoFiles, p.IgnoredGoFiles...)
	taskFiles = append(taskFiles, p.CgoFiles...)
	if e != nil {
		// task files may be ignored for build
		if _, ok := e.(*build.NoGoError); !ok || len(taskFiles) == 0 {
			err = e
			return
		}
	}

	tasks, err := loadTasks(dir, taskFiles)
	if err != nil {
		return
	}

	taskSet = &TaskSet{Name: p.Name, Dir: p.Dir, ImportPath: p.ImportPath, Tasks: tasks}

	return
}

func expandDir(dir string) (expanded string, err error) {
	expanded, err = filepath.Abs(dir)
	if err != nil {
		return
	}

	if !isFileExist(dir) {
		err = fmt.Errorf("Directory %s does not exist", dir)
		return
	}

	return
}

func loadTasks(dir string, files []string) (tasks []Task, err error) {
	taskFiles := filterTaskFiles(files)
	for _, taskFile := range taskFiles {
		ts, e := parseTasks(filepath.Join(dir, taskFile))
		if e != nil {
			err = e
			return
		}

		tasks = append(tasks, ts...)
	}

	return
}

func filterTaskFiles(files []string) (taskFiles []string) {
	for _, f := range files {
		if isTaskFile(f, "_task.go") {
			taskFiles = append(taskFiles, f)
		}
	}

	return
}

func parseTasks(filename string) (tasks []Task, err error) {
	taskFileSet := token.NewFileSet()
	f, err := parser.ParseFile(taskFileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return
	}

	for _, d := range f.Decls {
		n, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if n.Recv != nil {
			continue
		}

		actionName := n.Name.String()
		if isTask(actionName, "Task") {
			usage, desc, e := parseUsageAndDesc(n.Doc.Text())
			if e != nil {
				continue
			}

			name := convertActionNameToTaskName(actionName)
			t := Task{Name: name, ActionName: actionName, Usage: usage, Description: desc}
			tasks = append(tasks, t)
		}
	}

	return
}

func isTaskFile(name, suffix string) bool {
	if strings.HasSuffix(name, suffix) {
		return true
	}

	return false
}

func isTask(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Task" is ok
		return true
	}

	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}

func convertActionNameToTaskName(s string) string {
	n := strings.TrimPrefix(s, "Task")
	return dasherize(n)
}

func parseUsageAndDesc(doc string) (usage, desc string, err error) {
	reader := bufio.NewReader(bytes.NewReader([]byte(doc)))
	r := regexp.MustCompile("\\S")
	var usageParts, descParts []string

	line, err := readLine(reader)
	for err == nil {
		if len(descParts) == 0 && r.MatchString(line) {
			usageParts = append(usageParts, line)
		} else {
			descParts = append(descParts, line)
		}

		line, err = readLine(reader)
	}

	if err == io.EOF {
		err = nil
	}

	usage = strings.Join(usageParts, " ")
	usage = strings.TrimSpace(usage)

	desc = strings.Join(descParts, "\n")
	desc = strings.TrimSpace(desc)

	return
}

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}
