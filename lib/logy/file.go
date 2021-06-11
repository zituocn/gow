package logy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// StorageType 文件存储类型
type StorageType int

const (

	// StorageTypeMinutes 按分钟存储
	StorageTypeMinutes StorageType = iota

	// StorageTypeHour 按小时存储
	StorageTypeHour

	// StorageTypeDay 按天存储
	StorageTypeDay

	// StorageTypeMonth 按月份存储
	StorageTypeMonth
)

var (
	//formats 文件存储
	formats = map[StorageType]string{
		StorageTypeMinutes: "2006-01-02-15-04",
		StorageTypeHour:    "2006-01-02-15",
		StorageTypeDay:     "2006-01-02",
		StorageTypeMonth:   "2006-01",
	}

	//defaultMaxDay 默认最大保存天数
	defaultMaxDay = 7
)

// getFileFormat 获取文件存储格式
func (s StorageType) getFileFormat() string {
	return formats[s]
}

// SetFileFormat 设置文件存储格式
func SetFileFormat(s StorageType, format string) {
	formats[s] = format
}

// FileWriterOptions FileWriterOptions
type FileWriterOptions struct {
	StorageType   StorageType //存储类型
	StorageMaxDay int         //最大保存天数
	Dir           string      //目录
	Prefix        string      //前缀
	date          string      //日期
}

// FileWriter 文件存储实现
type FileWriter struct {
	FileWriterOptions
	sync.Mutex
	file *os.File
}

// WriteLog 接口实现
func (fw *FileWriter) WriteLog(t time.Time, level int, b []byte) {
	fw.Lock()
	fw.writeFile(t)
	fw.file.Write(b)
	fw.Unlock()
}

func (fw *FileWriter) writeFile(t time.Time) {
	newDate := t.Format(fw.StorageType.getFileFormat())
	if fw.date != newDate && fw.file != nil {
		fw.file.Close()
		fw.file = nil
	}
	if fw.file == nil {
		//目录
		dir := filepath.Dir(fw.Dir)
		err := os.MkdirAll(dir, 755)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%s.%s.log", fw.Prefix, newDate)
		file, err := os.OpenFile(filepath.Join(fw.Dir, fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(err)
		}
		fw.file = file
		fw.date = newDate
	}
}

// NewFileWriter return a new FileWriter
//	struct FileWriterOptions
func NewFileWriter(opts ...FileWriterOptions) *FileWriter {
	opt := prepareFileWriterOption(opts)
	fw := &FileWriter{
		FileWriterOptions: opt,
	}
	if fw.Dir[len(fw.Dir)-1:] != "/" {
		fw.Dir = fw.Dir + "/"
	}
	fw.clearLog()
	go fw.startTimer()
	return fw
}

//===============private=================

func prepareFileWriterOption(opts []FileWriterOptions) FileWriterOptions {
	var opt FileWriterOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Dir == "" {
		opt.Dir = "./"
	}
	if opt.StorageMaxDay <= 0 {
		opt.StorageMaxDay = defaultMaxDay
	}

	return opt
}

func (fw *FileWriter) startTimer() {
	now := time.Now()
	next := now.Add(time.Second * 3600)
	second := time.Duration(next.Sub(now).Seconds())
	fw.timer(second)
}

// timer
func (fw *FileWriter) timer(seconds time.Duration) {
	timer := time.NewTicker(seconds * time.Second)
	for {
		select {
		case <-timer.C:
			{
				fw.clearLog()
				nextTimer := time.NewTicker(3600 * time.Second)
				for {
					select {
					case <-nextTimer.C:
						{
							fw.startTimer()
							return
						}
					}
				}
			}
		}
	}
}

// clearLog() remove dir logs file
func (fw *FileWriter) clearLog() {
	files := getDirFiles(fw.Dir)
	now := time.Now()
	for _, item := range files {
		modTime := item.ModTime
		flag := modTime.Add(time.Hour * 24 * time.Duration(fw.StorageMaxDay-1)).Before(now)
		if flag {
			os.Remove(fmt.Sprintf("%s%s", fw.Dir, item.Name))
		}
	}
}

// FileInfo file info
type FileInfo struct {
	Name    string
	ModTime time.Time
	Size    int64
}

// getDirFiles
func getDirFiles(path string) (files []*FileInfo) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	files = make([]*FileInfo, 0)
	for _, fi := range dir {
		if !fi.IsDir() {
			files = append(files, &FileInfo{
				Name:    fi.Name(),
				ModTime: fi.ModTime(),
				Size:    fi.Size(),
			})
		}
	}

	return
}
