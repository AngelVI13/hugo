package utils

import (
	"bufio"
	"os"
	"strings"
)

// InputWaiting check if input is waiting
func InputWaiting() {

}

// ReadInput reads cmd line input
func ReadInput(info *SearchInfo) {

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// if something is piped to this .exe

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		info.stopped = true

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.Contains(text, "quit") {
			info.Quit = true
		}
	} else {
		return
	}
}
