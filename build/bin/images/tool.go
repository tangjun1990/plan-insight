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

func main2() {
	resizeImage("./boxbase-4.jpg", "./boxbase-5.jpg", 2200, 2760)
}

func main() {
	var files []string

	root := "./base3/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		targetFile := strings.Replace(file, "base3/", "formatted3/6-1-", -1)

		// 服装270*480，汽车480*270，椅子和生活方式300*400, 图案和材质350*350
		if strings.Index(file, ".png") > 0 || strings.Index(file, ".jpg") > 0 || strings.Index(file, ".jpeg") > 0 {
			resizeImage(file, targetFile, 350, 350)
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
