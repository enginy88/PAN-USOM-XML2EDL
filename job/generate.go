package job

import (
	"PAN-USOM-XML2EDL/app"

	"os"
	"regexp"
	"strconv"
	"sync"
)

func (iocRecords iocRecordSlice) generateSingleFileEDL(limit int) {

	var edlSlice []string

	for i, iocRecord := range iocRecords {
		if limit > 0 && i >= limit {
			app.LogInfo.Println("LIMITED RECORDS: " + strconv.Itoa(limit) + " record(s) limited by count.")
			break
		}

		edlSlice = append(edlSlice, iocRecord.URL)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	writeToFile(&wg, "edl.txt", edlSlice)

	wg.Wait()
}

func (iocRecords iocRecordSlice) generateMultiFileEDL(limit int) {

	var ipSlice []string
	var domainSlice []string
	var urlSlice []string

	regexPattern := `^(?P<leading_whitespace>[^\S\r\n]*)?(?P<scheme>http[s]?://)?(?:(?P<domain>[^\s/?#]+\.[^0-9\s\./?#:]+[^\s\./?#:]*)|(?P<ip>(?:[0-9]{1,3}\.){3}[0-9]{1,3}))(?P<port>:[0-9]{1,5})?(?:(?P<root>/+)(?P<path>[^\s?#]*))?(?:(?P<query>\?[^\s#]*)?(?P<fragment>#.*)?)?(?P<trailing_whitespace>[^\S\r\n]*)?$`

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		app.LogErr.Fatalln("FATAL ERROR: Error compiling regex pattern: '" + regexPattern + "'! (" + err.Error() + ")")
	}

	ipIndex := re.SubexpIndex("ip")
	domainIndex := re.SubexpIndex("domain")
	pathIndex := re.SubexpIndex("path")

	if ipIndex == -1 || domainIndex == -1 || pathIndex == -1 {
		app.LogErr.Fatalln("FATAL ERROR: Error getting regex subexpression index!")
	}

	matchCount := 0
	noMatchCount := 0
	for i, iocRecord := range iocRecords {
		if limit > 0 && i >= limit {
			app.LogInfo.Println("LIMITED RECORDS: " + strconv.Itoa(limit) + " record(s) limited by count.")
			break
		}

		matches := re.FindStringSubmatch(iocRecord.URL)

		if matches != nil && (matches[ipIndex] != "" || matches[domainIndex] != "") {

			if matches[domainIndex] != "" {
				domainSlice = append(domainSlice, matches[domainIndex])
			}

			if matches[ipIndex] != "" {
				ipSlice = append(ipSlice, matches[ipIndex])
			}

			urlSlice = append(urlSlice, matches[domainIndex]+matches[ipIndex]+"/"+matches[pathIndex])

			matchCount++

		} else {
			noMatchCount++

			if appSett.LogUnparsable {
				app.LogWarn.Println("UNPARSABLE RECORD: '" + iocRecord.URL + "'")
			}
		}

	}

	app.LogInfo.Println("PROCESSED RECORDS:  " + strconv.Itoa(matchCount) + " record(s) are parsed and processed.")

	if noMatchCount > 0 {
		app.LogWarn.Println("SKIPPED RECORDS: " + strconv.Itoa(noMatchCount) + " record(s) cannot be parsed and skipped.")
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go writeToFile(&wg, "edl-ip.txt", compact(ipSlice))
	go writeToFile(&wg, "edl-domain.txt", compact(domainSlice))
	go writeToFile(&wg, "edl-url.txt", compact(urlSlice))

	wg.Wait()

}

func compact[T comparable](slice []T) []T {
	keys := make(map[T]struct{})
	list := []T{}

	for _, s := range slice {
		if _, exists := keys[s]; !exists {
			keys[s] = struct{}{}
			list = append(list, s)
		}
	}

	return list
}

func writeToFile(wg *sync.WaitGroup, filename string, linesSlice []string) {

	defer wg.Done()

	filePath := appFlag.OutputDir + "/" + filename
	file, err := os.Create(filePath)
	if err != nil {
		app.LogErr.Fatalln("FATAL ERROR: Error creating output file: '" + filePath + "'! (" + err.Error() + ")")
	}
	defer file.Close()

	lineCount := 0
	for _, line := range linesSlice {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			app.LogErr.Fatalln("FATAL ERROR: Error writing to output file: '" + filePath + "'! (" + err.Error() + ")")
		}
		lineCount++
	}

	app.LogInfo.Println("WRITTEN RECORDS:  " + strconv.Itoa(lineCount) + " record(s) written to file: '" + filename + "'.")
}
