package pdf

import (
	"fmt"
	"github.com/zituocn/gow/lib/logy"
	"github.com/zituocn/gow/lib/pdf/core"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	pdf           *core.Report
	path          string
	width, height float64
	margin        core.Scope
	tempFilePath  string
}

//NewServerImage 获取远程 图片
func NewServerImage(url, path string, pdf *core.Report) (img *Image, err error) {
	var (
		fileName string
	)
	//fileName, _ = utils.GetUUID()
	if url != "" {
		fileNameList := strings.Split(url, "/")
		urlFileName := fileNameList[len(fileNameList)-1]
		fileNameList2 := strings.Split(urlFileName, ".")
		fileName = fileNameList2[0]
	}

	tempFilePath := fmt.Sprintf("%s/%s.png", path, fileName)

	//判断下载的图片是否存在，如果存在不需要重复下载，如果不存在则需要下载
	imagesJpeg := fmt.Sprintf("%s/%s.jpeg", path, fileName)

	if exists(imagesJpeg) {
		//logy.Info("有这个图片")
		w, h := GetImageWidthAndHeight(imagesJpeg)
		img = &Image{
			pdf:          pdf,
			path:         imagesJpeg,
			width:        float64(w / 10),
			height:       float64(h / 10),
			tempFilePath: imagesJpeg,
		}
		return
	}

	if !strings.Contains(url, "http") {
		img = NewImageFromServer(tempFilePath, pdf)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(tempFilePath)
	if err != nil {
		return
	}
	io.Copy(file, resp.Body)
	defer file.Close()

	img = NewImageFromServer(tempFilePath, pdf)
	return
}

//exists 图片是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//NewImageFromServer 从远程请求图片
func NewImageFromServer(path string, pdf *core.Report) *Image {
	var tempFilePath string
	picturePath, _ := filepath.Abs(path)
	imageType, _ := GetImageType(picturePath)
	if imageType == "png" {
		index := strings.LastIndex(picturePath, ".")
		tempFilePath = path[0:index+1] + "jpeg"
		err := ConvertPNG2JPEG(picturePath, tempFilePath)
		if err != nil {
			logy.Errorf("[PDF]转换图片出错:%v", err)
		}
		picturePath = tempFilePath
	}

	w, h := GetImageWidthAndHeight(picturePath)
	image := &Image{
		pdf:          pdf,
		path:         picturePath,
		width:        float64(w / 10),
		height:       float64(h / 10),
		tempFilePath: tempFilePath,
	}
	if tempFilePath != "" {
		pdf.AddCallBack(image.delTempImage)
	}
	if strings.Contains(path, "temp") {
		if _, err := os.Stat(path); err != nil {
		}
		os.Remove(path)
	}

	return image
}

func NewImage(path string, pdf *core.Report) *Image {
	if _, err := os.Stat(path); err != nil {
		fmt.Println(path, err)
		panic("the path error")
	}

	var tempFilePath string
	picturePath, _ := filepath.Abs(path)
	imageType, _ := GetImageType(picturePath)
	if imageType == "png" {
		index := strings.LastIndex(picturePath, ".")
		tempFilePath = picturePath[0:index] + ".jpeg"
		err := ConvertPNG2JPEG(picturePath, tempFilePath)
		if err != nil {
			panic(err)
		}
		picturePath = tempFilePath
	}

	w, h := GetImageWidthAndHeight(picturePath)
	image := &Image{
		pdf:          pdf,
		path:         picturePath,
		width:        float64(w / 10),
		height:       float64(h / 10),
		tempFilePath: tempFilePath,
	}
	if tempFilePath != "" {
		pdf.AddCallBack(image.delTempImage)
	}

	return image
}

func NewImageWithWidthAndHeight(path string, width, height float64, pdf *core.Report) *Image {
	contentWidth, contentHeight := pdf.GetContentWidthAndHeight()
	if width > contentWidth {
		width = contentWidth
	}
	if height > contentHeight {
		height = contentHeight
	}

	if _, err := os.Stat(path); err != nil {
		panic("the path error")
	}

	var tempFilePath string
	picturePath, _ := filepath.Abs(path)
	imageType, _ := GetImageType(picturePath)

	if imageType == "png" {
		index := strings.LastIndex(picturePath, ".")
		tempFilePath = picturePath[0:index] + ".jpeg"
		err := ConvertPNG2JPEG(picturePath, tempFilePath)
		if err != nil {
			panic(err.Error())
		}
		picturePath = tempFilePath
	}

	w, h := GetImageWidthAndHeight(picturePath)
	if float64(h)*width/float64(w) > height {
		width = float64(w) * height / float64(h)
	} else {
		height = float64(h) * width / float64(w)
	}
	image := &Image{
		pdf:          pdf,
		path:         picturePath,
		width:        width,
		height:       height,
		tempFilePath: tempFilePath,
	}

	if tempFilePath != "" {
		pdf.AddCallBack(image.delTempImage)
	}

	return image
}

func (image *Image) SetMargin(margin core.Scope) *Image {
	margin.ReplaceMarign()
	image.margin = margin
	return image
}

func (image *Image) GetHeight() float64 {
	return image.height
}
func (image *Image) GetWidth() float64 {
	return image.width
}

// 自动换行
func (image *Image) GenerateAtomicCell() error {
	var (
		sx, sy = image.pdf.GetXY()
	)

	x, y := sx+image.margin.Left, sy+image.margin.Top
	_, pageEndY := image.pdf.GetPageEndXY()
	if y < pageEndY && y+image.height > pageEndY {
		image.pdf.AddNewPage(false)
	}

	image.pdf.Image(image.path, x, y, x+image.width, y+image.height)
	sx, _ = image.pdf.GetPageStartXY()
	image.pdf.SetXY(sx, y+image.height+image.margin.Bottom)
	return nil
}

func (image *Image) delTempImage(report *core.Report) {
	if image.tempFilePath == "" {
		return
	}

	if _, err := os.Stat(image.tempFilePath); err != nil {
		return
	}

	os.Remove(image.tempFilePath)
}
