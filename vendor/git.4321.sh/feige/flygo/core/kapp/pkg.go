package kapp

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"git.4321.sh/feige/flygo/core/utils/xcolor"
	"git.4321.sh/feige/flygo/core/utils/xtime"
)

var (
	startTime    string
	goVersion    string
	flygoVersion string
)

var (
	appName         string // appName 、serviceName
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildStatus     string
	buildTime       string
)

func init() {
	if appName == "" {
		appName = os.Getenv(EnvAppName)
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = xtime.TS.Format(time.Now())
	setBuildTime(buildTime)
	goVersion = runtime.Version()
	initEnv()

	flygoVersion = "unknown version"
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, value := range info.Deps {
			if value.Path == "git.4321.sh/feige/flygo" {
				flygoVersion = value.Version
			}
		}
	}
}

// SetAppName ...
func SetAppName(name string) {
	appName = name
}

func Name() string {
	return appName
}

func AppVersion() string {
	return buildAppVersion
}

func FlygoVersion() string {
	return flygoVersion
}

func BuildTime() string {
	return buildTime
}

func BuildUser() string {
	return buildUser
}

func BuildHost() string {
	return buildHost
}

func setBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

func HostName() string {
	return hostName
}

func StartTime() string {
	return startTime
}

func GoVersion() string {
	return goVersion
}

func PrintVersion() {
	fmt.Printf("%-20s : %s\n", xcolor.Green("flygo"), xcolor.Blue("I am flygo"))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppName"), xcolor.Blue(appName))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppHost"), xcolor.Blue(HostName()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("Region"), xcolor.Blue(AppRegion()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("Zone"), xcolor.Blue(AppZone()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("AppVersion"), xcolor.Blue(buildAppVersion))
	fmt.Printf("%-20s : %s\n", xcolor.Green("FlygoVersion"), xcolor.Blue(flygoVersion))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildUser"), xcolor.Blue(buildUser))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildHost"), xcolor.Blue(buildHost))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildTime"), xcolor.Blue(BuildTime()))
	fmt.Printf("%-20s : %s\n", xcolor.Green("BuildStatus"), xcolor.Blue(buildStatus))
}
