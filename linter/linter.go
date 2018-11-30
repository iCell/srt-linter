package linter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const lineBreak = '\n'

var ErrNoContent = errors.New("No Content")
var ErrExtraSpace = errors.New("Extra Space")
var ErrInvalidCount = errors.New("Invalid Count")
var ErrTimeFormat = errors.New("Invalid Time Format")
var ErrStartEnd = errors.New("end time should later than start time")
var ErrTimeThanLast = errors.New("start time should later or equal than last end time")
var ErrStartNum = errors.New("Subtitle index should start from 1")
var ErrInvalidNum = errors.New("Subtitle index should increase 1")

// Subtitle struct contain a numeric counter identifying each sequential subtitle,
// start and end time that subtitle should appear on the screen
// subtitle text itself
type Subtitle struct {
	Num   int
	Start time.Time
	End   time.Time
	Text  string

	newLine bool
}

// LintResult describe the error detail by linter, and the error line
type LintResult struct {
	LineNum int
	Error   error
}

type process int

const (
	pNum process = iota
	pTime
	pText
	pNewline
)

type Linter struct {
	FilePath string

	pre     *Subtitle
	current *Subtitle
	proc    process
}

// NewLinter construct a linter struct with initial values
func NewLinter(f string) *Linter {
	return &Linter{
		FilePath: f,
		proc:     pNum,
		pre:      nil,
		current:  &Subtitle{},
	}
}

// Lint srt content
// ncluding timeline check, count check
// when the check all passed, it will return nil
func (l *Linter) Lint() []LintResult {
	f, err := os.Open(l.FilePath)
	if err != nil {
		panic(fmt.Sprintf("open file %s error: %s", l.FilePath, err.Error()))
	}
	defer f.Close()

	lineNum := 1
	r := bufio.NewReader(f)
	var result []LintResult

	for {
		b, err := r.ReadBytes(lineBreak)
		if err == io.EOF {
			return result
		}
		if err != nil {
			panic(err)
		}
		err = l.procLine(b)
		if err != nil {
			result = append(
				result,
				LintResult{
					Error:   err,
					LineNum: lineNum,
				},
			)
		}
		lineNum++
	}
}

func (l *Linter) procLine(b []byte) error {
	var err error
	line := strings.TrimSpace(string(b))

	switch l.proc {
	case pNum:
		err = l.procNum(line)
	case pTime:
		err = l.procTime(line)
	case pText:
		err = l.procText(line)
	}

	if l.current.newLine {
		err = l.validatePreCur()
		l.reset()
	}

	return err
}

func (l *Linter) procNum(line string) error {
	if line == "" {
		return ErrExtraSpace
	}

	c, err := strconv.Atoi(line)
	if err != nil {
		l.proc++
		return ErrInvalidCount
	}
	l.proc++
	l.current.Num = c
	return nil
}

func (l *Linter) procTime(line string) error {
	if line == "" {
		return ErrExtraSpace
	}

	l.proc++

	times := strings.Split(line, " --> ")
	if len(times) != 2 {
		return ErrTimeFormat
	}

	if strings.Contains(line, ".") {
		return ErrTimeFormat
	}

	layout := "15:04:05.000"
	start, errS := time.Parse(layout, strings.Replace(times[0], ",", ".", 1))
	end, errE := time.Parse(layout, strings.Replace(times[1], ",", ".", 1))
	if errS != nil || errE != nil {
		return ErrTimeFormat
	}

	if end.Unix() < start.Unix() {
		return ErrStartEnd
	}

	l.current.Start = start
	l.current.End = end
	return nil
}

func (l *Linter) procText(line string) error {
	if line == "" {
		if l.proc == pText {
			l.current.newLine = true
			return nil
		}
		return ErrExtraSpace
	}

	l.current.Text += line
	return nil
}

func (l *Linter) validatePreCur() error {
	if l.pre == nil && l.current.Num != 1 {
		return ErrStartNum
	}
	if l.current.Text == "" {
		return ErrNoContent
	}
	if l.pre != nil && l.current.Start.Unix() < l.pre.End.Unix() {
		return ErrTimeThanLast
	}
	if l.pre != nil && l.current.Num-l.pre.Num != 1 {
		return ErrInvalidNum
	}
	return nil
}

func (l *Linter) reset() {
	l.proc = pNum
	l.pre = l.current
	l.current = &Subtitle{}
}
