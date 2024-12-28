package logger

import (
	"container/list"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/logger"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

var levelArry = [4]string{DEBUG, INFO, WARN, ERROR}
var colorArry = [4]string{MagentaBold, BlueBold, YellowBold, RedBold}

var levelMap = map[string]int{
	DEBUG: 0,
	INFO:  1,
	WARN:  2,
	ERROR: 3,
}

type LoggerStruct struct {
	level string
	// logger   *log.Logger
	model    string
	function string
}

type Interface interface {
	LogMode(logger.LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type logetContent struct {
	format string
	v      []interface{}
}

type logQueueStruct struct {
	*list.List
	Mu     sync.Mutex
	RWLock sync.RWMutex
}

var logQueue logQueueStruct
var logLogger *log.Logger

func init() {
	logQueue = logQueueStruct{
		List: list.New(),
	}
}

func logOutput() {
	if !logQueue.Mu.TryLock() {
		return
	}
	for logQueue.Len() > 0 {
		logQueue.RWLock.RLock()
		e := logQueue.Front()
		logQueue.RWLock.RUnlock()
		switch e.Value.(type) {
		case logetContent:
			content := e.Value.(logetContent)
			logLogger.Printf(content.format, content.v...)
		default:
		}
		logQueue.RWLock.Lock()
		logQueue.Remove(e)
		logQueue.RWLock.Unlock()
	}
	logQueue.Mu.Unlock()
}

func NewLogger(level string, args ...any) (*LoggerStruct, error) {
	level = strings.ToUpper(level)
	if _, ok := levelMap[level]; !ok {
		return nil, errors.New("level not found")
	}
	var err error
	filename := "data/logs/backend.log"
	if len(args) > 0 {
		if args[0] == "stdout" {
			logLogger = log.New(io.MultiWriter(os.Stdout), "", log.Ldate|log.Ltime|log.Lmicroseconds)
			return &LoggerStruct{
				level: level,
				// logger: log.New(io.MultiWriter(os.Stdout), "", log.Ldate|log.Ltime|log.Lmicroseconds),
			}, nil
		}
		filename = args[0].(string)
	}
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, os.FileMode(0775))

	if os.IsNotExist(err) {
		err = os.MkdirAll("data/logs", os.FileMode(0775))
		if err != nil {
			return nil, err
		}
		logFile, _ = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, os.FileMode(0775))
	}

	logLogger = log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime|log.Lmicroseconds)
	return &LoggerStruct{
		level: level,
		// logger: log.New(io.MultiWriter(logPipe, os.Stdout), "", log.Ldate|log.Ltime),
		// logger: log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime|log.Lmicroseconds),
	}, nil
}

func (l *LoggerStruct) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	for logQueue.Len() > 0 {
		select {
		case <-ctx.Done():
			cancel()
			return
		default:
			continue
		}
	}
}

func (l *LoggerStruct) Clone() *LoggerStruct {
	return &LoggerStruct{
		level: l.level,
		// logger:   log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime|log.Lmicroseconds),
		model:    l.model,
		function: l.function,
	}
}

func (l *LoggerStruct) SetModel(model string) {
	l.model = model
}

func (l *LoggerStruct) SetFunction(function string) {
	l.function = function
}

func (l *LoggerStruct) SetLevel(level string) {
	l.level = strings.ToUpper(level)
}

func (l *LoggerStruct) Fatal(v ...interface{}) {
	l.Error("%v", v...)
	l.Close()
	os.Exit(1)
}

func (l *LoggerStruct) Print(level, format string, v ...interface{}) {
	if levelMap[level] >= levelMap[l.level] {
		logQueue.RWLock.Lock()
		if l.model != "" && l.function != "" {
			// l.logger.Printf("- "+colorArry[levelMap[level]]+level+Reset+" - "+Blue+"["+l.model+"."+l.function+"]"+Reset+" - "+format, v...)
			logQueue.PushBack(logetContent{
				format: "- " + colorArry[levelMap[level]] + level + Reset + " - " + Blue + "[" + l.model + "." + l.function + "]" + Reset + " - " + format,
				v:      v,
			})
		} else if l.model != "" {
			// l.logger.Printf("- "+colorArry[levelMap[level]]+level+Reset+" - "+Blue+"["+l.model+"]"+Reset+" - "+format, v...)
			logQueue.PushBack(logetContent{
				format: "- " + colorArry[levelMap[level]] + level + Reset + " - " + Blue + "[" + l.model + "]" + Reset + " - " + format,
				v:      v,
			})
		} else {
			// l.logger.Printf("- "+colorArry[levelMap[level]]+level+Reset+" - "+format, v...)
			logQueue.PushBack(logetContent{
				format: "- " + colorArry[levelMap[level]] + level + Reset + " - " + format,
				v:      v,
			})
		}
		logQueue.RWLock.Unlock()
		go logOutput()
	}
}

func (l *LoggerStruct) Println(level string, v ...interface{}) {
	l.Print(level, "%v", v...)
}

func (l *LoggerStruct) Debug(format string, v ...interface{}) {
	l.Print(DEBUG, format, v...)
}

func (l *LoggerStruct) Info(format string, v ...interface{}) {
	l.Print(INFO, format, v...)
}

func (l *LoggerStruct) Warn(format string, v ...interface{}) {
	l.Print(WARN, format, v...)
}

func (l *LoggerStruct) Error(format string, v ...interface{}) {
	l.Print(ERROR, format, v...)
}
