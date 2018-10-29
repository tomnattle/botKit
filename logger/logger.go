package logger

import (
	"bytes"
	"fmt"
	"github.com/ifchange/botKit/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

	writer = &Logger{
		sourceFileName: cfg.LogFile,
		buffer:         make(chan []byte, 1000),
		timestamp:      time.Now(),
	}
	writer.generateFileName()
	go writer.dumpStart()
}

type Logger struct {
	file     *os.File
	fileName string

	sourceFileName string
	buffer         chan []byte
	timestamp      time.Time
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

func (ins *Logger) dumpStart() {
	ticker := time.Tick(time.Duration(1) * time.Second)

	fileCheckTickerTimes, logFileCleanTickerTimes := 0, 0

	fileCheckDuration := time.Duration(1) * time.Minute
	logFileCleanDuration := time.Duration(24) * time.Hour

	ins.checkFile()

	defer func() {
		// using a closure
		// keep ins pointer
		ins.file.Close()
	}()

	for {
		<-ticker

		fileCheckTickerTimes++
		if fileCheckTickerTimes > int(fileCheckDuration.Seconds()) {
			ins.checkFile()
			fileCheckTickerTimes = 0
		}

		logFileCleanTickerTimes++
		if logFileCleanTickerTimes > int(logFileCleanDuration.Seconds()) {
			ins.cleanLogFile()
			logFileCleanTickerTimes = 0
		}

		ins.dumpAll()
	}
}

func (ins *Logger) cleanLogFile() {
	filePath, fileName := filepath.Split(ins.fileName)
	rd, err := ioutil.ReadDir(filePath)
	if err != nil {
		fmt.Printf("cleanLogFile error %v %v", filePath, err)
		return
	}

	now := time.Now()

	for _, fi := range rd {
		subNames := strings.SplitAfter(fi.Name(), fileName)
		if len(subNames) != 2 {
			continue
		}
		createDay, err := time.Parse("_2006_01_02", subNames[1])
		if err != nil {
			continue
		}
		if now.Sub(createDay) > 7*time.Duration(24)*time.Hour {
			os.Remove(filePath + fi.Name())
		}
	}
}

func (ins *Logger) checkFile() {
	if time.Now().Format("20060102") != ins.timestamp.Format("20060102") {
		ins.timestamp = time.Now()
		ins.generateFileName()
	}

	if ins.file == nil || ins.file.Sync() != nil {
		file, err := os.OpenFile(ins.fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("logger broken, can't open log file %s %v", ins.fileName, err)
			return
		}
		err = file.Sync()
		if err != nil {
			fmt.Printf("logger sync error %v", err)
			return
		}
		ins.file = file
	}

	if ins.file.Name() == ins.fileName {
		return
	}
	ins.file.Close()

	file, err := os.OpenFile(ins.fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("logger broken, can't open log file %s %v", ins.fileName, err)
		return
	}
	err = file.Sync()
	if err != nil {
		fmt.Printf("logger sync error %v", err)
		return
	}
	ins.file = file
}

func (ins *Logger) generateFileName() {
	ins.fileName = ins.sourceFileName + ins.timestamp.Format("_2006_01_02")
}

func (ins *Logger) dumpAll() {
allClear:
	for {
		select {
		case data := <-ins.buffer:
			_, err := ins.file.Write(data)
			if err != nil {
				fmt.Printf("logger try write file error %s %v", string(data), err)
			}
		default:
			err := ins.file.Sync()
			if err != nil {
				fmt.Printf("logger sync error %v", err)
			}
			break allClear
		}
	}
}
