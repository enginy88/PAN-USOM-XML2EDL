package app

import (
	"io"
	"strconv"

	"github.com/JeremyLoy/config"
)

const (
	defaultFeedURL     = "https://www.usom.gov.tr/url-list.xml"
	defaultDaysOld     = 0
	defaultLimitCount  = 0
	defaultNoSort      = false
	defaultLogSilence  = false
	defaultSkipVerify  = false
	defaultMultiThread = false
	defaultFileTest    = false
)

type AppSettStruct struct {
	FeedURL     string `config:"PANUSOMXML2EDL_FEED_URL"`
	DaysOld     int    `config:"PANUSOMXML2EDL_DAYS_OLD"`
	LimitCount  int    `config:"PANUSOMXML2EDL_LIMIT_COUNT"`
	NoSort      bool   `config:"PANUSOMXML2EDL_NO_SORT"`
	LogSilence  bool   `config:"PANUSOMXML2EDL_LOG_SILENCE"`
	SkipVerify  bool   `config:"PANUSOMXML2EDL_SKIP_VERIFY"`
	MultiThread bool   `config:"PANUSOMXML2EDL_MULTI_THREAD"`
	FileTest    bool   `config:"PANUSOMXML2EDL_FILE_TEST"`
}

var appSett *AppSettStruct

func GetAppSett() *AppSettStruct {

	appSettObject := new(AppSettStruct)
	appSett = appSettObject

	loadAppSett()
	checkAppSett()

	return appSett

}

func loadAppSett() {

	err := config.From("./appsett.env").FromEnv().To(appSett)
	if err != nil {
		LogErr.Fatalln("Cannot find/load 'appsett.env' file! (" + err.Error() + ")")
	}

}

func checkAppSett() {

	if appSett.LogSilence == true {
		LogWarn.SetOutput(io.Discard)
		LogInfo.SetOutput(io.Discard)
	}

	if appSett.FileTest == true {

		LogWarn.Println("CONFIG MSG: FileTest object value set to ('" + strconv.FormatBool(appSett.FileTest) + "'), no live data will be fetched.")
	}

	if appSett.SkipVerify == true {
		LogWarn.Println("CONFIG MSG: SkipVerify object value set to ('" + strconv.FormatBool(appSett.SkipVerify) + "'), which may be led to security risks.")
	}

	if appSett.MultiThread == true {
		LogWarn.Println("CONFIG MSG: MultiThread object value set to ('" + strconv.FormatBool(appSett.MultiThread) + "'), parallel processing is enabled.")
	}

	if appSett.DaysOld < 0 || appSett.DaysOld > 10000 {
		LogErr.Fatalln("DaysOld object value ('" + strconv.Itoa(appSett.DaysOld) + "') should be between 0 & 10.000!")
	}
	if appSett.LimitCount < 0 || appSett.LimitCount > 10000000 {
		LogErr.Fatalln("LimitCount object value ('" + strconv.Itoa(appSett.DaysOld) + "') should be between 0 & 10.000.000!")
	}
	if appSett.NoSort == true {
		LogWarn.Println("CONFIG MSG: NoSort object value set to ('" + strconv.FormatBool(appSett.NoSort) + "'), no sorting will be done and the order will be preserved.")
	}

	if len(appSett.FeedURL) > 2000 {
		LogErr.Fatalln("FeedUrl object ('" + appSett.FeedURL + "') should has less than 2000 chars!")
	}

	if len(appSett.FeedURL) < 1 {
		appSett.FeedURL = defaultFeedURL
		LogWarn.Println("CONFIG MSG: Using default value ('" + appSett.FeedURL + "') for FeedUrl object.")
	}

}
