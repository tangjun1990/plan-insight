package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"

	"github.com/disintegration/imaging"
)

//func main() {
//	resizeImage("./boxbase.jpg", "./boxbase-3.jpg", 5000, 3500)
//}

func main() {
	var files []string

	root := "./base/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		targetFile := strings.Replace(file, "base", "formatted", -1)

		if strings.Index(file, ".png") > 0 || strings.Index(file, ".jpg") > 0 || strings.Index(file, ".jpeg") > 0 {
			resizeImage(file, targetFile, 110, 220)
		}
	}
}

func resizeImage(inputImage, targetImage string, width, height int) {
	//读取本地文件，本地文件尺寸300*400
	imgData, _ := ioutil.ReadFile(inputImage)
	buf := bytes.NewBuffer(imgData)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	image = imaging.Resize(image, width, height, imaging.Lanczos)
	_ = imaging.Save(image, targetImage)
}
