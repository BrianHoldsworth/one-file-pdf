// -----------------------------------------------------------------------------
// (c) balarabe@protonmail.com                                      License: MIT
// :v: 2018-03-11 00:38:55 90E51C                              [test/support.go]
// -----------------------------------------------------------------------------

// Contains slightly-altered testing functions from Zircon-Go lib:
// github.com/balacode/zr

package main

import "fmt"     // standard
import "os"      // standard
import "reflect" // standard
import "runtime" // standard
import "strings" // standard
import "testing" // standard
import "time"    // standard

const showSourceFileNames = false

// does nothing, just allows this package to build properly
// this package is only used for testing
func main() {
}

// CallerList returns a human-friendly list of strings showing the
// call stack with each calling method or function's name and line number.
//
// The most immediate callers are listed first, followed by their callers,
// and so on. For brevity, 'runtime.*' and 'syscall.*'
// and other top-level callers are not included.
func CallerList() []string {
	var ret []string
	var i = 0
	for {
		i++
		var programCounter, filename, lineNo, _ = runtime.Caller(i)
		var funcName = runtime.FuncForPC(programCounter).Name()
		//
		// end loop on reaching a top-level runtime function
		if funcName == "" ||
			funcName == "runtime.goexit" ||
			funcName == "runtime.main" ||
			funcName == "testing.tRunner" ||
			strings.Contains(funcName, "HandlerFunc.ServeHTTP") {
			break
		}
		// skip runtime/syscall functions, but continue the loop
		if strings.Contains(funcName, "zr.Callers") ||
			strings.Contains(funcName, "zr.CallerList") ||
			strings.Contains(funcName, "zr.Error") ||
			strings.Contains(funcName, "zr.Log") ||
			strings.Contains(funcName, "zr.logAsync") ||
			strings.HasPrefix(funcName, "runtime.") ||
			strings.HasPrefix(funcName, "syscall.") {
			continue
		}
		// let the file name's path use the right kind of OS path separator
		// (by default, the file name contains '/' on all platforms)
		if string(os.PathSeparator) != "/" {
			filename = strings.Replace(filename, "/", string(os.PathSeparator), -1)
		}
		// remove parent module/function names
		if index := strings.LastIndex(funcName, "/"); index != -1 {
			funcName = funcName[index+1:]
		}
		if strings.Count(funcName, ".") > 1 {
			funcName = funcName[strings.Index(funcName, ".")+1:]
		}
		// remove unneeded punctuation from function names
		for _, find := range []string{"(", ")", "*"} {
			if strings.Contains(funcName, find) {
				funcName = strings.Replace(funcName, find, "", -1)
			}
		}
		var line string
		if showSourceFileNames {
			line = fmt.Sprintf("func:%-30s  ln:%4d  file:%-30s",
				funcName, lineNo, filename)
		} else {
			line = fmt.Sprintf("%s:%d", funcName, lineNo)
		}
		ret = append(ret, line)
	}
	return ret
} //                                                                  CallerList

// TCaller returns the name of the unit test function.
func TCaller() string {
	for _, funcName := range CallerList() {
		if strings.HasPrefix(funcName, "TCaller") ||
			strings.HasPrefix(funcName, "TEqual") ||
			strings.HasPrefix(funcName, "TFail") ||
			strings.HasPrefix(funcName, "TFalse") ||
			strings.HasPrefix(funcName, "TTrue") ||
			strings.Contains(funcName, ".func") ||
			!strings.Contains(funcName, ".Test") {
			continue
		}
		return funcName
	}
	return "<no-caller>"
} //                                                                     TCaller

// TEqual asserts that result is equal to expect.
func TEqual(t *testing.T, result interface{}, expect interface{}) bool {
	var makeStr = func(val interface{}) string {
		switch val := val.(type) {
		case nil:
			return "nil"
		case bool:
			if val {
				return "true"
			}
			return "false"
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64, uintptr:
			return fmt.Sprintf("%d", val)
		case float32, float64:
			var s = fmt.Sprintf("%f", val)
			if strings.Contains(s, ".") {
				for strings.HasSuffix(s, "0") {
					s = s[:len(s)-1]
				}
				for strings.HasSuffix(s, ".") {
					s = s[:len(s)-1]
				}
			}
			return s
		case string:
			return val
		case time.Time: // use date part without time and time zone
			var s = val.Format(time.RFC3339)[:19] // "2006-01-02T15:04:05Z07:00"
			if strings.HasSuffix(s, "T00:00:00") {
				s = s[:10]
			}
			return s
		case fmt.Stringer:
			return val.String()
		}
		return fmt.Sprintf("(type: %v value: %v)", reflect.TypeOf(val), val)
	}
	if makeStr(result) != makeStr(expect) {
		t.Logf("\n LOCATION: %s \n EXPECTED: %s \n RETURNED: %s \n",
			TCaller(), makeStr(expect), makeStr(result))
		t.Fail()
		return false
	}
	return true
} //                                                                      TEqual

//end
