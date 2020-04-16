package reporter

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/google/goterm/term"
)

// Message Holds a message with arguments
type Message struct {
	title string
	args  map[string]string
}

// Reporter support reporting functions
type Reporter struct {
	message *Message
	level   int
}

// Task unit of work
type Task struct {
	reporter *Reporter
	level    int
	lastLog  string
}

// colorIx
// 0: initial message
// 1: default color
// 2: section
// 8: success message
// 9: error message
func printMessageAndArgs(indent int, title string, args map[string]string, colorIx int) {
	var tab = ""
	for i := 0; i < indent; i++ {
		tab += "  "
	}
	fmt.Print(tab)
	switch colorIx {
	case 0:
		fmt.Println(term.Yellow(title))
	case 1:
		fmt.Println(term.Yellow(title))
	case 2:
		fmt.Println(term.BBlue(title))
	case 8:
		fmt.Println(term.Green(title))
	case 9:
		fmt.Println(term.Red(title))
	default:
		fmt.Println(term.Yellow(title))
	}
	if args != nil {
		for name, value := range args {
			fmt.Printf(term.Cyanf("  %s- %s: %s\n", tab, name, value))
		}
	}
}

func printError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		fmt.Println(term.Red("AWS Error is:"))
		fmt.Println(term.Red(" - code:" + aerr.Code()))
		fmt.Println(term.Red(" - error:" + aerr.Error()))
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(term.Red("Error is:"))
		fmt.Println(term.Red(err.Error()))
	}

}

// NewMessage create a message
func NewMessage(title string) *Message {
	var result = Message{
		title: title,
		args:  make(map[string]string),
	}
	return &result
}

// WithArg add arg to task description
func (message *Message) WithArg(name string, value string) *Message {
	message.args[name] = value
	return message
}

// WithArgs add arg to task description
func (message *Message) WithArgs(args map[string]*string) *Message {
	for key, element := range args {
		message.args[key] = *element
	}
	return message
}

func (message *Message) print(indent int, withArgs bool, extra string, colorIx int) {
	var args map[string]string = nil
	if withArgs {
		args = message.args
	}
	printMessageAndArgs(indent, message.title+extra, args, colorIx)
}

// NewReporterM constructor
func NewReporterM(message *Message) *Reporter {
	var result = Reporter{
		message: message,
		level:   0,
	}
	return &result
}

// NewReporter simple ctor
func NewReporter(title string) *Reporter {
	return NewReporterM(NewMessage(title))
}

// Start start the task being reported
func (reporter *Reporter) Start() *Task {
	// we log the reporter title
	if reporter.level == 0 {
		fmt.Println("")
	}

	reporter.message.print(reporter.level, true, "", 0)

	var result = Task{
		level:    1 + reporter.level,
		reporter: reporter,
	}
	return &result
}

// Ok display the fact that the reporter ends successfully
func (reporter *Reporter) Ok() {
	// we log the reporter title
	printMessageAndArgs(reporter.level, reporter.message.title+": SUCCESS", nil, 8)
}

// Okr success with result to display
func (reporter *Reporter) Okr(result map[string]string) {
	// we log the reporter title
	printMessageAndArgs(reporter.level, reporter.message.title+": SUCCESS", result, 8)
}

// Fail indicates that the reporter failed
func (reporter *Reporter) Fail(err error) {
	// we log the reporter title
	reporter.message.print(reporter.level, false, ": FAILED", 9)
	printError(err)
}

// SubM create sub reporter
func (task *Task) SubM(message *Message) *Task {
	var result = Reporter{
		message: message,
		level:   task.level + 1,
	}
	return result.Start()
}

// Sub create sub reporter
func (task *Task) Sub(title string) *Task {
	return task.SubM(NewMessage(title))
}

// LogM a message in the task
func (task *Task) LogM(message *Message) *Task {
	// we log the message title
	task.lastLog = message.title
	message.print(task.level, true, "", 1)
	return task
}

// SectionM log a section title
func (task *Task) SectionM(message *Message) *Task {
	// we log the message title
	task.lastLog = message.title
	message.print(task.level, true, "", 2)
	return task
}

// Log a message in the task
func (task *Task) Log(title string) *Task {
	// we log the message title
	task.LogM(NewMessage(title))
	return task
}

// Section create section title
func (task *Task) Section(title string) *Task {
	// we log the message title
	task.SectionM(NewMessage(title))
	return task
}

// Okr indicates success and print some values
func (task *Task) Okr(result map[string]string) {
	task.reporter.Okr(result)
	// printMessageAndArgs(task.level, task.lastLog+": SUCCESS", result)
}

// Ok indicates success
func (task *Task) Ok() {
	task.reporter.Ok()
	// if task.lastLog != "" {
	// 	printMessageAndArgs(task.level, task.lastLog+": SUCCESS", nil)
	// }
}

// Fail indicates failure
func (task *Task) Fail(err error) {
	task.reporter.Fail(err)
	// printMessageAndArgs(task.level, task.lastLog+": FAILURE", nil)
	// printError(err)
}

// Experiment experiment
func Experiment() {
	r := NewReporterM(NewMessage("my first reporter").WithArg("arg1", "42"))
	t0 := r.Start()
	t0.Log("Let's go")
	t1 := t0.Sub("Sub task...")
	t1.Log("step 1")
	t1.Okr(map[string]string{"a": "42"})
	t1.Log("step 2")
	t1.Fail(errors.New("There is an error"))
	r.Ok()
}
