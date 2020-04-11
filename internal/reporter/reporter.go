package reporter

import "fmt"

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
	level int
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

func (message *Message) print(indent int, withArgs bool, extra string) {
	var tab = ""
	for i := 0; i < indent; i++ {
		tab += "  "
	}
	fmt.Print(tab)
	fmt.Println(message.title + extra)
	if withArgs {
		for name, value := range message.args {
			fmt.Printf("  %s- %s: %s\n", tab, name, value)
		}
	}
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
	reporter.message.print(reporter.level, true, "")

	var result = Task{
		level: 1 + reporter.level,
	}
	return &result
}

// Ok display the fact that the reporter ends successfully
func (reporter *Reporter) Ok() {
	// we log the reporter title
	reporter.message.print(reporter.level, false, ": SUCCESS")
}

// Fail indicates that the reporter failed
func (reporter *Reporter) Fail(err *error) {
	// we log the reporter title
	reporter.message.print(reporter.level, false, ": FAILED")
	fmt.Println(err)
}

// SubReporter create sub reporter
func (task *Task) SubReporter(message *Message) *Reporter {
	var result = Reporter{
		message: message,
		level:   task.level + 1,
	}
	return &result
}

// LogM a message in the task
func (task *Task) LogM(message *Message) *Task {
	// we log the message title
	message.print(task.level, true, "")
	return task
}

// Log a message in the task
func (task *Task) Log(title string) *Task {
	// we log the message title
	task.LogM(NewMessage(title))
	return task
}

// Experiment experiment
func Experiment() {
	t1 := NewReporterM(NewMessage("my first reporter").WithArg("arg1", "42")).Start()
	t1.Log("Let's go")
}
