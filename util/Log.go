package util

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var Log = logrus.New()

func initLog() {
	Log.Out = os.Stdout
	var loglevel logrus.Level
	if err := loglevel.UnmarshalText([]byte(Config.GetString("log.level"))); err != nil {
		Log.Panicf("设置日志级别失败: %v", err)
	}
	Log.SetLevel(loglevel)
	if Config.GetBool("server.dev") {
		Log.SetReportCaller(true)
	}
	p, _ := filepath.Abs(Config.GetString("log.path"))
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if os.MkdirAll(p, os.ModePerm) != nil {
			Log.Warn("创建日志文件夹失败!")
		}
	}

	NewSimpleLogger(Log, p, Config.GetUint("log.saveFiles"))
}

func NewSimpleLogger(log *logrus.Logger, path string, save uint) {
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer(path, "debug", save),
		logrus.InfoLevel:  writer(path, "info", save),
		logrus.WarnLevel:  writer(path, "warn", save),
		logrus.ErrorLevel: writer(path, "error", save),
	}, &myFormatter{})
	log.AddHook(lfHook)
}

type myFormatter struct{}

func (mf *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	msg := fmt.Sprintf("[%s] [%s] %s\n", time.Now().Local().Format("2006-01-02 15:04:05"), strings.ToUpper(entry.Level.String()), entry.Message)

	return []byte(msg), nil
}

func writer(logPath, level string, save uint) io.Writer {
	logFullPath := path.Join(logPath, level)
	fileSuffix := time.Now().Local().Format("2006-01-02") + Config.GetString("log.suffix")

	logier, err := rotatelogs.New(
		logFullPath+"-"+fileSuffix,
		//rotatelogs.WithLinkName(logFullPath),
		rotatelogs.WithRotationCount(int(save)),
		rotatelogs.WithRotationTime(time.Hour*Config.GetDuration("log.rotateHours")),
	)

	if err != nil {
		panic(err)
	}
	return logier
}
