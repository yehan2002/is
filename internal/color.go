package internal

import (
	"flag"
	"fmt"
	"os"

	"github.com/enescakir/emoji"
	"github.com/ttacon/chalk"
)

//NoColorFlag if colors are disabled
var NoColorFlag bool

var (
	check = emoji.CheckMark.String()
	cross = emoji.CrossMark.String()
)

//Messages messages
var Messages = struct {
	Start, Run, PassSuite, FailSuite, PassTest, FailTest, Err1, Err2 *Message
}{
	Start:     newMessage("****** Running %s ******", chalk.Cyan),
	PassSuite: newMessage("****** %s Passed ******", chalk.Green),
	FailSuite: newMessage("****** %s Failed ******", chalk.Red),
	Run:       newMessage("Running %s", chalk.Cyan),
	PassTest:  newMessage(" -- PASS %s "+check+" (%.2fs)", chalk.Green),
	FailTest:  newMessage(" -- FAIL %s "+cross+" (%.2fs)", chalk.Red),
	Err1:      newMessage("--- Error: %s", chalk.Red),
	Err2:      newMessage("--- Fail: %s", chalk.Red),
}

//Message a message
type Message struct {
	normal, color string
}

func newMessage(s string, color chalk.Color) *Message {
	return &Message{s + "\n", color.Color(s) + "\n"}
}

//Print print the message to stdout
func (m *Message) Print(color bool, v ...interface{}) {
	if !flag.Parsed() {
		flag.Parse()
	}
	if color && !NoColorFlag {
		fmt.Printf(m.color, v...)
	} else {
		fmt.Printf(m.normal, v...)

	}
}

func init() {
	envNoColor := os.Getenv("COLOR_TEST") == "true"
	flag.BoolVar(&NoColorFlag, "nocolor", !envNoColor, "turns off colors")
}

//NoColor disables color
func NoColor() {
	NoColorFlag = true
}
