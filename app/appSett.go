package app

import (
	"io"
	"strconv"

	"github.com/JeremyLoy/config"
)

const (
	defaultFeedURL       = "https://www.usom.gov.tr/url-list.xml"
	defaultDaysOld       = 0
	defaultLimitCount    = 0
	defaultSingleOutput  = false
	defaultNoSort        = false
	defaultLogSilence    = false
	defaultLogUnparsable = false
	defaultSkipVerify    = false
	defaultFileTest      = false
)

type AppSettStruct struct {
	FeedURL       string `config:"PANUSOMXML2EDL_FEED_URL"`
	DaysOld       int    `config:"PANUSOMXML2EDL_DAYS_OLD"`
	LimitCount    int    `config:"PANUSOMXML2EDL_LIMIT_COUNT"`
	SingleOutput  bool   `config:"PANUSOMXML2EDL_SINGLE_OUTPUT"`
	NoSort        bool   `config:"PANUSOMXML2EDL_NO_SORT"`
	LogSilence    bool   `config:"PANUSOMXML2EDL_LOG_SILENCE"`
	LogUnparsable bool   `config:"PANUSOMXML2EDL_LOG_UNPARSABLE"`
	SkipVerify    bool   `config:"PANUSOMXML2EDL_SKIP_VERIFY"`
	FileTest      bool   `config:"PANUSOMXML2EDL_FILE_TEST"`
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
		LogErr.Fatalln("FATAL ERROR: Cannot find/load 'appsett.env' file! (" + err.Error() + ")")
	}

}

func checkAppSett() {

	if appSett.LogSilence {
		LogWarn.SetOutput(io.Discard)
		LogInfo.SetOutput(io.Discard)
	}

	if appSett.FileTest {

		LogWarn.Println("CONFIG MSG: FileTest object value set to ('" + strconv.FormatBool(appSett.FileTest) + "'), no live data will be fetched.")
	}

	if appSett.SkipVerify {
		LogWarn.Println("CONFIG MSG: SkipVerify object value set to ('" + strconv.FormatBool(appSett.SkipVerify) + "'), which may be led to security risks.")
	}

	if appSett.LogUnparsable {
		LogWarn.Println("CONFIG MSG: LogUnparsable object value set to ('" + strconv.FormatBool(appSett.LogUnparsable) + "'), each record that cannot be parsed will be logged.")
	}

	if appSett.DaysOld < 0 || appSett.DaysOld > 10000 {
		LogErr.Fatalln("FATAL ERROR: DaysOld object value ('" + strconv.Itoa(appSett.DaysOld) + "') should be between 0 & 10.000!")
	}

	if appSett.LimitCount < 0 || appSett.LimitCount > 10000000 {
		LogErr.Fatalln("FATAL ERROR: LimitCount object value ('" + strconv.Itoa(appSett.DaysOld) + "') should be between 0 & 10.000.000!")
	}

	if appSett.SingleOutput {
		LogWarn.Println("CONFIG MSG: SingleOutput object value set to ('" + strconv.FormatBool(appSett.SingleOutput) + "'), records will be written into a single file without being parsed or processed.")
	}

	if appSett.NoSort {
		LogWarn.Println("CONFIG MSG: NoSort object value set to ('" + strconv.FormatBool(appSett.NoSort) + "'), no sorting will be done and the order will be preserved.")
	}

	if len(appSett.FeedURL) > 2000 {
		LogErr.Fatalln("FATAL ERROR: FeedUrl object ('" + appSett.FeedURL + "') should has less than 2000 chars!")
	}

	if len(appSett.FeedURL) < 1 {
		appSett.FeedURL = defaultFeedURL
	} else {
		if !appSett.FileTest {
			LogWarn.Println("CONFIG MSG: Using custom value ('" + appSett.FeedURL + "') for FeedUrl object.")
		}
	}

}
