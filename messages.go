package is

import (
	"flag"
	"fmt"
	"os"
)

var noColorFlag bool

var messages = struct {
	start, run, passSuite, failSuite, passTest, failTest, err1, err2 *message
}{
	start:     &message{"****** Running %s ******\n", "\x1b[1;36m****** Running %s ******\x1b[0m\n"},
	passSuite: &message{"****** %s Passed ******\n", "\x1b[1;36m****** %s Passed ******\x1b[0m\n"},
	failSuite: &message{"****** %s Failed ******\n", "\x1b[1;36m****** %s Failed ******\x1b[0m\n"},
	run:       &message{"Running %s", "\x1b[36mRunning %s\x1b[0m"},
	passTest:  &message{" -- PASS \xf0\x9f\x97\xb8 (%.2fs)\n", "\x1b[36m -- \x1b[32mPASS \xf0\x9f\x97\xb8 (%.2fs)\x1b[0m\n"},
	failTest:  &message{" -- FAIL \xc3\x97 (%.2fs)\n", "\x1b[36m -- \x1b[31mFAIL \xc3\x97 (%.2fs)\x1b[0m\n"},
	err1:      &message{"--- Error: %s\n", "\x1b[31m--- Error: %s\x1b[0m\n"},
	err2:      &message{"--- Fail: %s\n", "\x1b[31m--- Fail: %s\x1b[0m\n"},
}

type message struct {
	normal, color string
}

func printf(m *message, color bool, v ...interface{}) {
	if color && !noColorFlag {
		fmt.Printf(m.color, v...)
	} else {
		fmt.Printf(m.normal, v...)
	}
}

//NoColor disables color
func NoColor() {
	noColorFlag = true
}

func init() {
	envNoColor := os.Getenv("COLOR_TEST") == "true"
	flag.BoolVar(&noColorFlag, "nocolor", !envNoColor, "turns off colors")
}
