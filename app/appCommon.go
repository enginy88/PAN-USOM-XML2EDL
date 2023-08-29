package app

import (
	"fmt"
	"log"
	"os"
)

var (
	LogErr    *log.Logger
	LogWarn   *log.Logger
	LogInfo   *log.Logger
	LogAlways *log.Logger
)

func init() {
	LogErr = log.New(os.Stderr, "(PAN-USOM-XML2EDL) ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	LogWarn = log.New(os.Stdout, "(PAN-USOM-XML2EDL) WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	LogInfo = log.New(os.Stdout, "(PAN-USOM-XML2EDL) INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	LogAlways = log.New(os.Stdout, "(PAN-USOM-XML2EDL) ALWAYS: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func Typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func FindString(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
