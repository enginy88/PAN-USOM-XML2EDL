package main

import (
	"PAN-USOM-XML2EDL/app"
	"PAN-USOM-XML2EDL/job"

	"fmt"
	"time"
)

var appFlag *app.AppFlagStruct
var appSett *app.AppSettStruct

func main() {

	start := time.Now()
	app.LogAlways.Println("HELLO MSG: Welcome to PAN-USOM-XML2EDL v2.5 by EY!")

	appFlag = app.GetAppFlag()
	appSett = app.GetAppSett()

	job.RunAllJobs(*appFlag, *appSett)

	duration := fmt.Sprintf("%.1f", time.Since(start).Seconds())
	app.LogAlways.Println("BYE MSG: All done in " + duration + "s, bye!")

}
