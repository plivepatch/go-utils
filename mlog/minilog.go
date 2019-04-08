package mlog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type LEVEL int32

var LevelLogFile LEVEL = DEBUG
var LevelConsole LEVEL = OFF
var maxFileSize int64
var maxFileCount int32
var dailyRolling bool = true
var consoleAppender bool = true
var RollingFile bool = false
var logObj *_FILE
var timeIntervel uint32 = 60 //seconds

const DATEFORMAT = "2006-01-02"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type _FILE struct {
	dir      string
	filename string
	_suffix  int
	isCover  bool
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

func SetFlags(flag int) {
	logObj.lg.SetFlags(flag)
}

func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

func SetLogLevelLogFile(_level LEVEL) {
	LevelLogFile = _level
}

func SetLogLevelConsole(_level LEVEL) {
	LevelConsole = _level
}

func SetTimeIntervel(_timeIntervel uint32) {
	timeIntervel = _timeIntervel
}

func SetRollingFile(fileDir string, fileName string, maxNumber int32, maxSize int64) {
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(MB)
	RollingFile = true
	dailyRolling = false
	mkdirlog(fileDir)
	logObj = &_FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()
	for i := 1; i <= int(maxNumber); i++ {
		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
			logObj._suffix = i
		} else {
			break
		}
	}
	if !logObj.isMustRename() {
		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
		logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}
	go fileMonitor()
}

func SetRollingDaily(fileDir, fileName string) {
	RollingFile = false
	dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	mkdirlog(fileDir)
	logObj = &_FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if !logObj.isMustRename() {
		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}
}

func mkdirlog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0755); err != nil {
			if os.IsPermission(err) {
				fmt.Println("create log dir error:", err.Error())
				e = err
			}
		}
	}
	return
}

func console(level string, s ...interface{}) {
	if consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(level, file, strconv.Itoa(line), fmt.Sprint(s...))
	}
}

func consolef(level string, s string) {
	if consoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Printf("%s %s %s %s\n", level, file, strconv.Itoa(line), s)
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func Debug(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= DEBUG {
		if logObj != nil {
			logObj.lg.Output(2, "DEBUG "+fmt.Sprintln(v))
		}
	}
	if LevelConsole <= DEBUG {
		console("DEBUG", v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= DEBUG {
		if logObj != nil {
			logObj.lg.Output(2, "DEBUG "+fmt.Sprintf(format, v...))
		}
	}
	if LevelConsole <= DEBUG {
		consolef("DEBUG", fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LevelLogFile <= INFO {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintln("INFO", v))
		}
	}
	if LevelConsole <= INFO {
		console("INFO", v...)
	}
}

func Infof(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= INFO {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintf("INFO "+format, v...))
		}
	}
	if LevelConsole <= INFO {
		consolef("INFO", fmt.Sprintf(format, v...))
	}
}

func Warn(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= WARN {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintln("WARN", v))
		}
	}
	if LevelConsole <= WARN {
		console("WARN", v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= WARN {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintf("WARN "+format, v...))
		}
	}
	if LevelConsole <= WARN {
		consolef("WARN", fmt.Sprintf(format, v...))
	}
}

func Error(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LevelLogFile <= ERROR {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintln("ERORR", v))
		}
	}
	if LevelConsole <= ERROR {
		console("ERORR", v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= ERROR {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintf("ERROR "+format, v...))
		}
	}
	if LevelConsole <= ERROR {
		consolef("ERORR", fmt.Sprintf(format, v...))
	}
}

func Fatal(v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}
	if LevelLogFile <= FATAL {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintln("FATAL", v))
		}
	}
	if LevelConsole <= FATAL {
		console("FATAL", v...)
	}
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if LevelLogFile <= FATAL {
		if logObj != nil {
			logObj.lg.Output(2, fmt.Sprintf("FATALT "+format, v...))
		}
	}
	if LevelConsole <= FATAL {
		consolef("FATALT", fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func (f *_FILE) isMustRename() bool {
	if dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	} else {
		if maxFileCount > 1 {
			if fileSize(f.dir+"/"+f.filename) >= maxFileSize {
				return true
			}
		}
	}
	return false
}

func (f *_FILE) rename() {
	if dailyRolling {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
			f.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
		}
	} else {
		f.coverNextOne()
	}
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(maxFileCount) + 1)
}

func (f *_FILE) coverNextOne() {
	f._suffix = f.nextSuffix()
	if f.logfile != nil {
		f.logfile.Close()
	}
	if isExist(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix))) {
		os.Remove(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix)))
	}
	os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+"."+strconv.Itoa(int(f._suffix)))
	f.logfile, _ = os.Create(f.dir + "/" + f.filename)
	f.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		return 1 * int64(TB)
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func fileMonitor() {
	timer := time.NewTicker((time.Duration)(timeIntervel) * time.Second)
	for {
		select {
		case <-timer.C:
			fileCheck()
		}
	}
}

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
	}
}
