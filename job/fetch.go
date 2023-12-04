package job

import (
	"PAN-USOM-XML2EDL/app"

	"crypto/tls"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

func getIocRecords() (iocRecords iocRecordSlice) {

	var input io.Reader

	if appSett.FileTest {
		file, err := os.Open("./url-list.xml")
		if err != nil {
			app.LogErr.Fatalln("FATAL ERROR: Cannot find/open 'url-list.xml' file! (" + err.Error() + ")")
		}
		defer file.Close()
		input = file

	} else {
		resp := fetchXMLFeed()
		input = strings.NewReader(string(resp))

	}

	xml, err := xmlquery.Parse(input)
	if err != nil {
		app.LogErr.Fatalln("FATAL ERROR: Cannot parse XML format! (" + err.Error() + ")")
	}

	result := xmlquery.FindOne(xml, "/usom-data/url-list[url-info/url]")

	if result == nil {
		app.LogErr.Fatalln("FATAL ERROR: Not an expected XML Data!")
	}

	for _, iter := range xmlquery.Find(result, "/url-info[url][date]") {

		url := iter.SelectElement("url").InnerText()
		dateString := iter.SelectElement("date").InnerText()
		date, err := time.Parse(time.DateTime, dateString)
		if err != nil {
			app.LogErr.Fatalln("FATAL ERROR: Cannot parse DateTime format! (" + err.Error() + ")")
		}

		iocRecords = append(iocRecords, iocRecord{URL: url, Date: date})

	}

	return iocRecords

}

func fetchXMLFeed() (xml []byte) {

	transport := http.DefaultTransport.(*http.Transport)
	if appSett.SkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	timeout := time.Second * time.Duration(appSett.TimeoutDuration)
	client := &http.Client{Transport: transport, Timeout: timeout}

	url := appSett.FeedURL

	response, err := client.Get(url)
	if err != nil {
		app.LogErr.Fatalln("FATAL ERROR: Cannot download requested data! (" + err.Error() + ")")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		app.LogErr.Fatalln("FATAL ERROR: Cannot read downloaded data! (" + err.Error() + ")")
	}

	return body

}
