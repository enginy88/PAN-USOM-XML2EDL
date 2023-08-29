package job

import (
	"PAN-USOM-XML2EDL/app"

	"encoding/json"
	"sort"
	"strconv"
	"time"
)

var appFlag app.AppFlagStruct
var appSett app.AppSettStruct

type iocRecord struct {
	URL  string
	Date time.Time
}

type iocRecordSlice []iocRecord

func (iocRecords iocRecordSlice) sortByTime() {

	sort.Slice(iocRecords, func(i, j int) bool { return iocRecords[i].Date.After(iocRecords[j].Date) })

}

func RunAllJobs(appFlagParam app.AppFlagStruct, appSettParam app.AppSettStruct) {

	appFlag = appFlagParam
	appSett = appSettParam

	appSettJSON, err := json.Marshal(appSett)
	if err != nil {
		app.LogErr.Fatalln("FATAL ERROR: Cannot marshal JSON data! (" + err.Error() + ")")
	}
	app.LogInfo.Println("RUNNING CONFIG: '" + string(appSettJSON) + "'")

	iocRecords := getIocRecords()

	app.LogInfo.Println("FETCHED RECORDS: " + strconv.Itoa(len(iocRecords)) + " record(s) fetched.")

	if !appSett.NoSort {

		iocRecords.sortByTime()

	}

	filteredIocRecords := iocRecords.filterByDate(appSett.DaysOld)

	if appSett.DaysOld > 0 {
		app.LogInfo.Println("FILTERED RECORDS: " + strconv.Itoa(len(filteredIocRecords)) + " record(s) filtered by date.")
	}

	if appSett.SingleOutput {
		filteredIocRecords.generateSingleFileEDL(appSett.LimitCount)
	} else {
		filteredIocRecords.generateMultiFileEDL(appSett.LimitCount)
	}

}

func (iocRecords iocRecordSlice) filterByDate(days int) (filteredIocRecords iocRecordSlice) {

	if days < 1 {
		return iocRecords
	}

	lastUpdated := iocRecords[0].Date
	comperableDate := lastUpdated.AddDate(0, 0, 0-days)

	for _, iocRecord := range iocRecords {

		if iocRecord.Date.After(comperableDate) {
			filteredIocRecords = append(filteredIocRecords, iocRecord)
			continue
		}

	}

	return filteredIocRecords

}
