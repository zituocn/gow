/*

file.go
日志文件存储实现

*/

package logy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type StorageType int

const (

	// StorageTypeMinutes 按分钟存储
	StorageTypeMinutes StorageType = iota

	// StorageTypeHour 按小时存储
	StorageTypeHour

	// StorageTypeDay 按天存储
	StorageTypeDay

	// StorageTypeMonth 按月存储
	StorageTypeMonth
)

func (s StorageType) getFileFormat() string {
	return formats[s]
}

var (
	formats = map[StorageType]string{
		StorageTypeMinutes: "2006-01-02-15-04",
		StorageTypeHour:    "2006-01-02-15",
		StorageTypeDay:     "2006-01-02",
		StorageTypeMonth:   "2006-01",
	}

	// defaultMaxDay 日志文件的默认最大保存天数
	// 7天之外的文件，会被自动清理
	defaultMaxDay = 7
)

// FileWriterOptions FileWrite 参数
type FileWriterOptions struct {

	// StorageType 存储时间类型
	StorageType StorageType

	// StorageMaxDay 日志最大保存天数
	StorageMaxDay int

	// Dir 保存目录
	Dir string

	// Prefix 前缀
	Prefix string

	// date 日期
	date string
}

// FileWriter 日志文件存储实现
type FileWriter struct {
	FileWriterOptions
	mu   *sync.Mutex
	file *os.File
}

// NewFileWriter return a new FileWriter
//	param: opts
func NewFileWriter(opts ...FileWriterOptions) *FileWriter {
	opt := prepareFileWriterOption(opts)
	fw := &FileWriter{
		FileWriterOptions: opt,
		mu:                &sync.Mutex{},
	}
	if fw.Dir[len(fw.Dir)-1:] != "/" {
		fw.Dir = fw.Dir + "/"
	}

	fw.clearLogFile()
	go fw.startTimer()
	return fw
}

// WriteLog write []byte to file
func (fw *FileWriter) WriteLog(now time.Time, level int, b []byte) {
	fw.mu.Lock()
	fw.writeFile(now)
	fw.file.Write(b)
	fw.mu.Unlock()
}

// writeFile write file to directory
func (fw *FileWriter) writeFile(now time.Time) {
	date := now.Format(fw.StorageType.getFileFormat())
	if fw.date != date && fw.file != nil {
		fw.file.Close()
		fw.file = nil
	}

	if fw.file == nil {
		dir := filepath.Dir(fw.Dir)
		err := os.MkdirAll(dir, 755)
		if err != nil {
			panic(err)
		}

		fileName := fmt.Sprintf("%s.%s.log", fw.Prefix, date)
		file, err := os.OpenFile(filepath.Join(fw.Dir, fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			panic(err)
		}
		fw.file = file
		fw.date = date
	}

}

//=========================timer and clear file===================

func (fw *FileWriter) startTimer() {
	now := time.Now()
	next := now.Add(time.Second * 3600)
	second := time.Duration(next.Sub(now).Seconds())
	fw.timer(second)
}

func (fw *FileWriter) timer(second time.Duration) {
	timer := time.NewTicker(second * time.Second)
	for {
		select {
		case <-timer.C:
			{
				fw.clearLogFile()
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

// clearLogFile delete the log file in directory fw.Dir
func (fw *FileWriter) clearLogFile() {
	files := getDirFiles(fw.Dir)
	now := time.Now()
	for _, item := range files {
		modTime := item.ModTime
		flag := modTime.Add(time.Hour * 24 * time.Duration(fw.StorageMaxDay-1)).Before(now)
		if flag {
			os.Remove(fw.Dir + item.Name)
		}
	}
}

// FileInfo file info
type FileInfo struct {
	Name    string
	ModTime time.Time
	Size    int64
}

// getDirFiles return files
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
