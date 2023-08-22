package job

import (
	"PAN-USOM-XML2EDL/app"
	"crypto/tls"
	"net/http"

	"io"
	"os"
	"strings"
	"time"

	"github.com/antchfx/xmlquery"
)

func getIocRecords() (iocRecords iocRecordSlice) {

	var input io.Reader

	if appSett.FileTest == true {
		file, err := os.Open("./url-list.xml")
		if err != nil {
			app.LogErr.Fatalln(err)
		}
		defer file.Close()
		input = file

	} else {
		resp := fetchXMLFeed()
		input = strings.NewReader(string(resp))

	}

	xml, err := xmlquery.Parse(input)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	result := xmlquery.FindOne(xml, "/usom-data/url-list[url-info/url]")

	if result == nil {
		app.LogErr.Fatalln("Not an expected XML Data!")
	}

	for _, iter := range xmlquery.Find(result, "/url-info[url][date]") {

		url := iter.SelectElement("url").InnerText()
		dateString := iter.SelectElement("date").InnerText()
		date, err := time.Parse(time.DateTime, dateString)
		if err != nil {
			app.LogErr.Fatalln(err)
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

	timeout := time.Duration(30 * time.Second)
	client := &http.Client{Transport: transport, Timeout: timeout}

	url := appSett.FeedURL

	response, err := client.Get(url)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		app.LogErr.Fatalln(err)
	}

	return body

}
