package job

import (
	"PAN-USOM-XML2EDL/app"

	"encoding/json"
	"os"
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
		app.LogErr.Fatalln(err)
	}
	app.LogInfo.Println("RUNNING CONFIG: " + string(appSettJSON))

	iocRecords := getIocRecords()

	app.LogInfo.Println("GET RECORDS: " + strconv.Itoa(len(iocRecords)) + " record(s) fetched.")

	if appSett.NoSort == false {

		iocRecords.sortByTime()

	}

	filteredIocRecords := iocRecords.filterByDate(appSett.DaysOld)

	app.LogInfo.Println("FILTER RECORDS: " + strconv.Itoa(len(filteredIocRecords)) + " record(s) filtered by date.")

	filteredIocRecords.generateEDL(appSett.LimitCount)

}

func (iocRecords iocRecordSlice) filterByDate(days int) (filteredIocRecords iocRecordSlice) {

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

func (iocRecords iocRecordSlice) generateEDL(limit int) {

	file, err := os.Create(appFlag.OutputDir + "/edl.txt")
	if err != nil {
		app.LogErr.Fatalln(err)
	}
	defer file.Close()

	for i, iocRecord := range iocRecords {
		if limit > 0 && i >= limit {
			app.LogInfo.Println("LIMIT RECORDS: " + strconv.Itoa(limit) + " record(s) limited by count.")
			break
		}
		_, err := file.WriteString(iocRecord.URL + "\n")
		if err != nil {
			app.LogErr.Fatalln(err)
		}
	}

}
