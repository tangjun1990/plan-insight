package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/tangjun1990/flygo/core/kcfg"
	"github.com/tangjun1990/flygo/core/kcfg/manager"
	"github.com/tangjun1990/flygo/core/utils/xgo"
)

type fileDataSource struct {
	path        string
	dir         string
	enableWatch bool
	changed     chan struct{}
}

func init() {
	manager.Register(manager.DefaultScheme, &fileDataSource{})
}

func (fp *fileDataSource) Parse(path string, watch bool) kcfg.ConfigType {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		panic("new datasource:" + err.Error())
	}
	dir := checkAndGetParentDir(absolutePath)
	fp.path = absolutePath
	fp.dir = dir
	fp.enableWatch = watch

	if watch {
		fp.changed = make(chan struct{}, 1)
		xgo.Go(fp.watch)
	}

	return extParser(path)
}

func extParser(configAddr string) kcfg.ConfigType {
	ext := filepath.Ext(configAddr)
	switch ext {
	case ".json":
		return kcfg.ConfigTypeJSON
	case ".toml":
		return kcfg.ConfigTypeToml
	default:
		panic("data source: invalid configuration type")
	}
	return ""
}

func (fp *fileDataSource) ReadConfig() (content []byte, err error) {
	return ioutil.ReadFile(fp.path)
}

func (fp *fileDataSource) Close() error {
	close(fp.changed)
	return nil
}

func (fp *fileDataSource) IsConfigChanged() <-chan struct{} {
	return fp.changed
}

func (fp *fileDataSource) watch() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		panic("new file watcher error:" + err.Error())
	}
	defer w.Close()

	configFile := filepath.Clean(fp.path)
	realConfigFile, _ := filepath.EvalSymlinks(fp.path)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.Events:
				currentConfigFile, _ := filepath.EvalSymlinks(fp.path)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
				if (filepath.Clean(event.Name) == configFile &&
					event.Op&writeOrCreateMask != 0) ||
					(currentConfigFile != "" && currentConfigFile != realConfigFile) {
					realConfigFile = currentConfigFile
					select {
					case fp.changed <- struct{}{}:
					default:
					}
				}
			case err := <-w.Errors:
				fmt.Println("read watch error:" + err.Error())
			}
		}
	}()
	err = w.Add(fp.dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func checkAndGetParentDir(path string) string {
	// check path is the directory
	isDir, err := isDirectory(path)
	if err != nil || isDir {
		return path
	}
	return getParentDirectory(path)
}

func isDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}

func getParentDirectory(dirctory string) string {
	if runtime.GOOS == "windows" {
		dirctory = strings.Replace(dirctory, "\\", "/", -1)
	}
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
