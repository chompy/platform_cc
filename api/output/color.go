/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package output

import "fmt"

const termColorFormat = "\033[%dm"
const termColorReset = "\033[0m"
const termColorSuccess = 32
const termColorError = 31
const termColorProgress = 34
const termColorWarn = 33
const termColorDebug = 35

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
