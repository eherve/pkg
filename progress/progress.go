package progress

import (
	"fmt"
	"strings"
	"time"
)

//BasicTemplate print only bar and percetage
const BasicTemplate = "{{bar}} {{percent}}"

//Progress progress interface
type Progress interface {
	Size() int64
	SetSize(s int64) *Progress
	Template() string
	SetTemplate(t string) *Progress
	Width() int64
	SetWidth(w int64) *Progress
	Interval() time.Duration
	SetInterval(i time.Duration) *Progress
	Progress() int64
	Write(p []byte) (n int, err error)
	Tick(n int) (int64, error)
	Clear(newline bool)
}

//Bar implements Progress
type Bar struct {
	size     int64
	template string
	// action   string
	width    int64
	interval time.Duration

	progress  int64
	lastDraw  time.Time
	maxLength int
}

//Size get the size of the data beeing processed
func (p *Bar) Size() int64 {
	return p.size
}

//SetSize set the size of the data beeing processed
func (p *Bar) SetSize(s int64) *Bar {
	p.size = s
	return p
}

//Template get the the draw template
func (p *Bar) Template() string {
	return p.template
}

//SetTemplate set the draw template
func (p *Bar) SetTemplate(t string) *Bar {
	p.template =
		strings.ReplaceAll(
			strings.ReplaceAll(
				t, "{{bar}}", "[%[1]s%[2]s]",
			),
			"{{percent}}", "%[3]d%%",
		)
	return p
}

//Width get the progress bar width
func (p *Bar) Width() int64 {
	return p.width
}

//SetWidth set the progress bar width
func (p *Bar) SetWidth(w int64) *Bar {
	p.width = w
	return p
}

//Interval get the progress bar drawing call interval
func (p *Bar) Interval() time.Duration {
	return p.interval
}

//SetInterval set the progress bar drawing call interval
func (p *Bar) SetInterval(i time.Duration) *Bar {
	p.interval = i
	return p
}

//Progress get the current progress
func (p *Bar) Progress() int64 {
	return p.progress
}

func (p *Bar) Write(b []byte) (n int, err error) {
	var n64 int64
	n64, err = p.Tick(len(b))
	n = int(n64)
	return
}

//Tick add n to progress
// returns progress
func (p *Bar) Tick(n int) (int64, error) {
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
func (p *Bar) Clear(newline bool) {
	p.progress = 0
	if newline {
		fmt.Printf("\n")
	}
}

func (p *Bar) drawProgress() error {
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

func (p *Bar) drawBar(progress, total int64) error {
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

//NewBar return a new Progress Bar
func NewBar() *Bar {
	return (&Bar{}).
		SetTemplate(BasicTemplate).
		SetWidth(int64(10)).
		SetInterval(time.Second)
}
