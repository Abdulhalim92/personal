package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func GetLogger() Logger {
	return Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
	return Logger{l.WithField(k, v)}
}

func init() {
	logg := logrus.New()
	logg.SetReportCaller(true)
	logg.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0777)
	if err != nil {
		log.Println(err)
	}

	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	logg.SetOutput(io.Discard)

	logg.AddHook(&writerHook{
		Writer:    []io.Writer{allFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	logg.SetLevel(logrus.TraceLevel)
	logg.SetLevel(logrus.InfoLevel)
	logg.SetLevel(logrus.FatalLevel)
	logg.SetLevel(logrus.ErrorLevel)
	logg.SetLevel(logrus.PanicLevel)

	e = logrus.NewEntry(logg)
}

//package logging
//
//import (
//	"fmt"
//	"github.com/sirupsen/logrus"
//	"io"
//	"os"
//	"path"
//	"runtime"
//)
//
//type writerHook struct {
//	Writer    []io.Writer
//	LogLevels []logrus.Level
//}
//
//func (hook *writerHook) Fire(entry *logrus.Entry) error {
//	line, err := entry.Bytes()
//	if err != nil {
//		return err
//	}
//	for _, w := range hook.Writer {
//		w.Write(line)
//	}
//	return err
//}
//
//func (hook *writerHook) Levels() []logrus.Level {
//	return hook.LogLevels
//}
//
//var e *logrus.Entry
//
//type Logger struct {
//	*logrus.Entry
//}
//
//func GetLogger() Logger {
//	return Logger{e}
//}
//
//func (l *Logger) GetLoggerWithField(k string, v interface{}) Logger {
//	return Logger{l.WithField(k, v)}
//}
//
//func init() {
//	log := logrus.New()
//	log.SetReportCaller(true)
//	log.Formatter = &logrus.JSONFormatter{
//		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
//			fileName := path.Base(frame.File)
//			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s: %d", fileName, frame.Line)
//		},
//		DisableHTMLEscape: false,
//		DisableTimestamp:  true,
//	}
//
//	err := os.MkdirAll("logs", 0777)
//	if err != nil {
//		panic(err)
//	}
//
//	allFiles, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_APPEND, 0777)
//	if err != nil {
//		panic(err)
//	}
//
//	log.SetOutput(io.Discard)
//
//	log.AddHook(&writerHook{
//		Writer:    []io.Writer{allFiles, os.Stdout},
//		LogLevels: logrus.AllLevels,
//	})
//
//	log.SetLevel(logrus.InfoLevel)
//	log.SetLevel(logrus.FatalLevel)
//	log.SetLevel(logrus.ErrorLevel)
//	log.SetLevel(logrus.TraceLevel)
//	//log.SetLevel(logrus.PanicLevel)
//	e = logrus.NewEntry(log)
//}
