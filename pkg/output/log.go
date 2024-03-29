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

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

// LogPath is the path to the log file.
const LogPath = ".platform_cc.log"
const logMaxSize = 2097152 // 2MB
const logOverageMax = logMaxSize * 2

func writeLogFile(msg string) {
	if !Logging {
		return
	}
	f, err := os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Warn("Could not open log file, " + err.Error())
		return
	}
	defer f.Close()
	ts := time.Now().Format(time.RFC1123)
	_, err = f.Write([]byte("[" + ts + "] " + msg + "\n"))
	if err != nil {
		Warn("Could not write log file, " + err.Error())
		return
	}
	f.Sync()
}

// LogDebug writes debug message to log file.
func LogDebug(msg string, data interface{}) {
	logMsg := "[DEBUG] " + msg
	if data != nil {
		dataJSON, _ := json.Marshal(data)
		logMsg += " " + string(dataJSON)
	}
	writeLogFile(logMsg)
	if Verbose {
		os.Stderr.Write([]byte(Color(">> "+logMsg+"\n", termColorDebug)))
	}
}

// LogInfo writes info to log file.
func LogInfo(msg string) {
	writeLogFile("[INFO] " + msg)
}

// LogWarn writes warning to log file.
func LogWarn(msg string) {
	writeLogFile("[WARN] " + msg)
}

// LogError writes error to log file.
func LogError(err error) {
	if err == nil {
		return
	}
	out := fmt.Sprintf("[ERROR] %s\n\n%s", err.Error(), errorStackTrace(err))
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			out += fmt.Sprintf("%+s:%d\n", f, f)
		}
	}
	writeLogFile(out)
}

// LogRotate trims the log file.
func LogRotate() error {
	// stat log file
	s, err := os.Stat(LogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.WithStack(err)
	}
	// calculate how many bytes over size limit log file is
	overage := s.Size() - logMaxSize
	// not reached max size
	if overage < 0 {
		return nil
	}
	// overage is too large, just delete file
	if overage > logOverageMax {
		return errors.WithStack(os.Remove(LogPath))
	}

	// open file
	f, err := os.Open(LogPath)
	if err != nil {
		return errors.WithStack(err)
	}

	// scan line by line for overage content
	scanner := bufio.NewScanner(f)
	bytesScanned := int64(0)
	for scanner.Scan() {
		bytesScanned += int64(len(scanner.Bytes()))
		if bytesScanned >= overage {
			break
		}
	}
	// scan line by line for content that should be remain
	trimmedContents := make([]byte, 0)
	for scanner.Scan() {
		trimmedContents = append(trimmedContents, scanner.Bytes()...)
		trimmedContents = append(trimmedContents, '\n')
	}
	f.Close()
	// delete old file
	if err := os.Remove(LogPath); err != nil {
		return errors.WithStack(err)
	}
	// create new
	f, err = os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	if _, err := f.Write(trimmedContents); err != nil {
		return errors.WithStack(err)
	}
	if err := f.Sync(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
