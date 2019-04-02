package console

import (
	"fmt"
	"strings"
	"time"
)

//ProgressBasicTemplate print only bar and percetage
const ProgressBasicTemplate = "{{bar}} {{percent}}"

//Progress implements Progress
type Progress struct {
	size     int64
	template string
	// action   string
	width    int64
	interval time.Duration

	progress  int64
	lastDraw  time.Time
	maxLength int
}

//NewProgress return a new Progress Progress
func NewProgress() *Progress {
	return (&Progress{template: ProgressBasicTemplate, width: int64(10), interval: time.Second})
}

//Size get the size of the data beeing processed
func (p *Progress) Size() int64 {
	return p.size
}

//SetSize set the size of the data beeing processed
func (p *Progress) SetSize(s int64) {
	p.size = s
}

//Template get the the draw template
func (p *Progress) Template() string {
	return p.template
}

//SetTemplate set the draw template
func (p *Progress) SetTemplate(t string) {
	p.template =
		strings.ReplaceAll(
			strings.ReplaceAll(
				t, "{{bar}}", "[%[1]s%[2]s]",
			),
			"{{percent}}", "%[3]d%%",
		)
}

//Width get the progress bar width
func (p *Progress) Width() int64 {
	return p.width
}

//SetWidth set the progress bar width
func (p *Progress) SetWidth(w int64) {
	p.width = w
}

//Interval get the progress bar drawing call interval
func (p *Progress) Interval() time.Duration {
	return p.interval
}

//SetInterval set the progress bar drawing call interval
func (p *Progress) SetInterval(i time.Duration) {
	p.interval = i
}

//Progress get the current progress
func (p *Progress) Progress() int64 {
	return p.progress
}

func (p *Progress) Write(b []byte) (n int, err error) {
	var n64 int64
	n64, err = p.Tick(len(b))
	n = int(n64)
	return
}

//Tick add n to progress
// returns progress
func (p *Progress) Tick(n int) (int64, error) {
	p.progress += int64(n)
	if p.progress != p.size {
		if err := p.drawProgress(); err != nil {
			return 0, err
		}
	} else {
		if err := p.drawBar(p.progress, p.size); err != nil {
			return 0, err
		}
		p.Clear(true)
		p.lastDraw = time.Time{}
	}
	return p.progress, nil
}

//Clear reset progress
// if newline true, print a new line
func (p *Progress) Clear(newline bool) {
	p.progress = 0
	if newline {
		fmt.Printf("\n")
	}
}

func (p *Progress) drawProgress() error {
	if !p.lastDraw.IsZero() && p.interval != -1 {
		nextDraw := p.lastDraw.Add(p.interval)
		if time.Now().Before(nextDraw) {
			return nil
		}
	}

	if err := p.drawBar(p.progress, p.size); err != nil {
		return err
	}

	p.lastDraw = time.Now()
	return nil
}

func (p *Progress) drawBar(progress, total int64) error {
	current := int64((float64(progress) / float64(total)) * float64(p.width))

	line := fmt.Sprintf(
		p.template,
		strings.Repeat("=", int(current)),
		strings.Repeat(" ", int(p.width-current)),
		(progress*100)/total)

	if len(line) < p.maxLength {
		line = fmt.Sprintf(
			"%s%s",
			line,
			strings.Repeat(" ", p.maxLength-len(line)))
	}
	p.maxLength = len(line)

	_, err := fmt.Print(line + "\r")
	return err
}
