package logger

import (
	"bytes"
	"fmt"
	"github.com/ifchange/botKit/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func Debugf(format string, v ...interface{}) {
	if config.GetConfig().Environment == "prod" {
		return
	}
	writer.Write(bytes.NewBufferString(fmt.Sprintln(fmt.Sprintf("[DEBUG]"+format, v...))).Bytes())
}

func Printf(format string, v ...interface{}) {
	writer.Write(bytes.NewBufferString(fmt.Sprintln(fmt.Sprintf(format, v...))).Bytes())
}

var (
	writer *Logger
)

func init() {
	cfg := config.GetConfig().Logger
	if cfg == nil {
		panic("logger config is nil")
	}

	now := time.Now()

	writer = &Logger{
		sourceFileName: cfg.LogFile,
		buffer:         make(chan []byte, 1000),
		mu:             new(sync.RWMutex),
	}
	writer.fileCheck(now)
	go writer.loopFileCheck()
	go writer.loopLogDump()
}

type Logger struct {
	sourceFileName string
	buffer         chan []byte
	fileName       string

	timestamp time.Time
	file      *os.File

	mu *sync.RWMutex
}

func (ins *Logger) Write(data []byte) (int, error) {
	dataLen := len(data)
	cp := make([]byte, dataLen)
	copy(cp, data)
	select {
	case ins.buffer <- cp:
		return dataLen, nil
	default:
		return 0, fmt.Errorf("try add data into logger buffer error, buffer is full")
	}
}

func (ins *Logger) loopFileCheck() {
	for now := range time.Tick(time.Duration(1) * time.Minute) {
		ins.mu.Lock()

		ins.fileCheck(now)
		ins.file.Sync()
		ins.cleanLogFile(now)

		ins.mu.Unlock()
	}
}

func (ins *Logger) fileCheck(now time.Time) {
	if now.Format("20060102") != ins.timestamp.Format("20060102") {
		ins.timestamp = now
		ins.fileName = ins.sourceFileName + ins.timestamp.Format("_2006_01_02")
	}

	if ins.file == nil || ins.file.Sync() != nil {
		ins.newLogFile()
	}

	if ins.file.Name() == ins.fileName {
		return
	}
	ins.file.Close()
	ins.newLogFile()
}

func (ins *Logger) newLogFile() {
	file, err := os.OpenFile(ins.fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("logger broken, can't open log file %s %v\n", ins.fileName, err)
		return
	}
	err = file.Sync()
	if err != nil {
		fmt.Printf("logger sync error %v\n", err)
		return
	}
	ins.file = file
}

func (ins *Logger) cleanLogFile(now time.Time) {
	filePath, fileName := filepath.Split(ins.fileName)
	dir, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Printf("cleanLogFile error %v %v\n", filePath, err)
		return
	}

	for _, file := range dir {
		subNames := strings.SplitAfter(file.Name(), fileName)
		if len(subNames) != 2 {
			continue
		}
		createDay, err := time.Parse("_2006_01_02", subNames[1])
		if err != nil {
			continue
		}
		if now.Sub(createDay) > 7*time.Duration(24)*time.Hour {
			os.Remove(filePath + file.Name())
		}
	}
}

func (ins *Logger) loopLogDump() {
	for data := range ins.buffer {
		ins.mu.RLock()
		_, err := ins.file.Write(data)
		if err != nil {
			fmt.Printf("logger try write file error %s %v\n", string(data), err)
		}
		ins.mu.RUnlock()
	}
}
