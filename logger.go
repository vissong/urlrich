package urlrich

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LOG_LEVEL_ALL   = 255
	LOG_LEVEL_TRACE = 6
	LOG_LEVEL_DEBUG = 5
	LOG_LEVEL_INFO  = 4
	LOG_LEVEL_WARN  = 3
	LOG_LEVEL_ERROR = 2
	LOG_LEVEL_FATAL = 1
	LOG_LEVE_OFF    = 0
)

type LeveledLogger interface {
	Debugf(string, ...interface{})
}

type Logger struct {
	l         *log.Logger
	w         *io.Writer
	level     int
	prefix    string
	calldepth int // 调用深度，用于打印原始的行号
}

func (log *Logger) log(level int, format string, vars ...interface{}) {
	if log.level >= level {
		log.l.Output(log.calldepth, fmt.Sprintf(format, vars...))
	}
}

func (log *Logger) GetLogger() *log.Logger {
	return log.l
}

func (log *Logger) IncrCalledppth(intValue int) {
	log.calldepth += intValue
}

func (log *Logger) SetPrefix(str string) {
	log.prefix = str
}

func (log *Logger) GetPrefix() string {
	if log.prefix != "" {
		return fmt.Sprintf("%s", log.prefix)
	}
	return ""
}

func (log *Logger) GetWriter() *io.Writer {
	return log.w
}

func (log *Logger) Println(vars ...interface{}) {
	log.Debug("%s", fmt.Sprint(vars...))
}

func (log *Logger) Trace(format string, vars ...interface{}) {
	if log.level >= LOG_LEVEL_TRACE {
		log.log(LOG_LEVEL_TRACE, "[TRACE] "+log.GetPrefix()+format, vars...)
	}
}

func (log *Logger) Debug(format string, vars ...interface{}) {
	if log.level >= LOG_LEVEL_DEBUG {
		log.log(LOG_LEVEL_DEBUG, "[DEBUG] "+log.GetPrefix()+format, vars...)
	}
}

func (log *Logger) Info(format string, vars ...interface{}) {
	if log.level >= LOG_LEVEL_INFO {
		log.log(LOG_LEVEL_INFO, "[INFO] "+log.GetPrefix()+format, vars...)
	}
}

func (log *Logger) Warn(format string, vars ...interface{}) {
	if log.level >= LOG_LEVEL_WARN {
		log.log(LOG_LEVEL_WARN, "[WARN] "+log.GetPrefix()+format, vars...)
	}
}

func (log *Logger) Error(format string, vars ...interface{}) {
	if log.level >= LOG_LEVEL_ERROR {
		log.log(LOG_LEVEL_ERROR, "[ERROR] "+log.GetPrefix()+format, vars...)
	}
}

func (log *Logger) Fatal(format string, vars ...interface{}) {
	if log.level >= LOG_LEVEL_FATAL {
		log.log(LOG_LEVEL_FATAL, "[FATAL] "+log.GetPrefix()+format, vars...)
	}
}

type RotateWriter struct {
	lock   sync.Mutex
	LogDir string
	Name   string
	fp     *os.File
	curLen int
	MaxLen int
}

func (w *RotateWriter) Write(b []byte) (int, error) {
	// w.lock.Lock()
	// defer w.lock.Unlock()
	l := len(b)
	if w.curLen+l > w.MaxLen {
		w.rotate()
	}
	w.curLen += l
	return w.fp.Write(b)
}

func NewRotateWriter(logDir, logName string, MaxLen int) *RotateWriter {
	var (
		err      error
		fileName string
		fileInfo os.FileInfo
		fp       *os.File
	)
	rw := &RotateWriter{LogDir: logDir, Name: logName, MaxLen: MaxLen}

	fileName = rw.LogDir + string(os.PathSeparator) + rw.Name

	// file exist
	if fileInfo, err = os.Stat(fileName); !os.IsNotExist(err) {
		rw.curLen = int(fileInfo.Size())
		if int(fileInfo.Size()) < rw.MaxLen {
			if fp, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
				panic(err)
			}
			rw.fp = fp
			return rw
		} else {
			// rename current file
			if err = os.Rename(fileName, fileName+"."+time.Now().Format(time.RFC3339)); err != nil {
				panic(err)
			}

		}

	}

	if fp, err = os.Create(fileName); err != nil {
		panic(err)
	}

	rw.fp = fp
	rw.curLen = 0

	return rw
}

func (w *RotateWriter) rotate() {
	w.lock.Lock()
	defer w.lock.Unlock()
	var (
		err      error
		fileName string
		fp       *os.File
	)

	// close current file
	if w.fp != nil {
		if err = w.fp.Close(); err != nil {
			panic(err)
		}
		w.fp = nil
	}

	fileName = w.LogDir + string(os.PathSeparator) + w.Name
	// file exist
	if _, err = os.Stat(fileName); !os.IsNotExist(err) {
		// rename current file
		if err = os.Rename(fileName, fileName+"."+time.Now().Format(time.RFC3339)); err != nil {
			panic(err)
		}
	}

	// create a new file
	if fp, err = os.Create(fileName); err != nil {
		panic(err)
	}
	w.fp = fp
	w.curLen = 0
}

func NewLogger(logDir, logName string, level int) *Logger {
	var (
		logger *log.Logger
		writer io.Writer
	)
	// 一个文件 500MB
	rw := NewRotateWriter(logDir, logName, 500<<20)
	if level >= LOG_LEVEL_DEBUG {
		writer = rw
	} else {
		writer = bufio.NewWriter(rw)
	}
	logger = log.New(writer, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return &Logger{
		l:         logger,
		w:         &writer,
		level:     level,
		prefix:    "",
		calldepth: 3}
}
