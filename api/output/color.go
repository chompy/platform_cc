package output

import "fmt"

const termColorFormat = "\033[%dm"
const termColorReset = "\033[0m"
const termColorSuccess = 32
const termColorError = 31
const termColorProgress = 34
const termColorWarn = 33

func color(msg string, color int) string {
	if !isTTY() {
		return msg
	}
	return fmt.Sprintf(termColorFormat, color) + msg + termColorReset
}

func colorSuccess(msg string) string {
	return color(msg, termColorSuccess)
}

func colorError(msg string) string {
	return color(msg, termColorError)
}

func colorWarn(msg string) string {
	return color(msg, termColorWarn)
}

func colorProgress(msg string) string {
	return color(msg, termColorProgress)
}
