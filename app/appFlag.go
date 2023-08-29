package app

import (
	"flag"
	"os"
	"path/filepath"
)

type AppFlagStruct struct {
	WorkingDir    string
	OutputDir     string
	workingDirRaw string
	outputDirRaw  string
}

var appFlag *AppFlagStruct

func GetAppFlag() *AppFlagStruct {

	appFlagObject := new(AppFlagStruct)
	appFlag = appFlagObject

	parseAppFlag()
	changeWorkingDir()

	return appFlag

}

func parseAppFlag() {

	workingDir := flag.String("dir", "", "Path of directory which contains 'appsett.env' file and also EDL file(s) are created.")
	outputgDir := flag.String("out", "", "Path of directory which EDL file(s) are created. (Overrides '-dir' option.)")
	flag.Parse()

	appFlag.workingDirRaw = *workingDir
	appFlag.outputDirRaw = *outputgDir

}

func changeWorkingDir() {

	origDir, err := os.Getwd()
	if err != nil {
		LogErr.Fatalln("FATAL ERROR: Cannot get working directory! (" + err.Error() + ")")
	}

	workingDir := origDir
	outputDir := origDir

	if appFlag.workingDirRaw != "" {

		err := os.Chdir(appFlag.workingDirRaw)
		if err != nil {
			LogErr.Fatalln("FATAL ERROR: Cannot change working directory! (" + err.Error() + ")")
		}

		newDir, err := os.Getwd()
		if err != nil {
			LogErr.Fatalln("FATAL ERROR: Cannot get working directory! (" + err.Error() + ")")
		}

		workingDir = newDir

		LogInfo.Println("CONFIG MSG: Flag 'dir' set, changing working directory from '" + origDir + "' to '" + newDir + "'.")
	}

	if appFlag.outputDirRaw != "" {

		outputDir = filepath.Join(workingDir, appFlag.outputDirRaw)

		LogInfo.Println("CONFIG MSG: Flag 'out' set, going to write the EDL file(s) to '" + outputDir + "' directory.")
	}

	appFlag.WorkingDir = workingDir
	appFlag.OutputDir = outputDir

}
