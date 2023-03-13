// Support fatal, error, warn, info, debug five log levels.
// Split log file by size and date.
// Clear expired log files.
package log

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
)

const (
    FatalLevel = iota
    ErrorLevel
    WarnLevel
    InfoLevel
    DebugLevel
)

const (
    FatalLevelConf = "fatal"
    ErrorLevelConf = "error"
    WarnLevelConf  = "warn"
    InfoLevelConf  = "info"
    DebugLevelConf = "debug"
)

const (
    FatalFlag = "[Fatal]"
    ErrorFlag = "[Error]"
    WarnFlag  = "[Warn]"
    InfoFlag  = "[Info]"
    DebugFlag = "[Debug]"
)

const (
    FlagTimeFmt = ".20060102."
)

const (
    DefaultMaxSize   = 100 // max size(MB) for one log file
    DefaultCheckDays = 1
    DefaultLogLevel  = "debug"
    DefaultTotalSize = 1024 // total size(MB) for all log files
)

var l *Logger
var onceSplit sync.Once
var onceClear sync.Once
var onceClearByTotalSizes sync.Once

type Logger struct {
    logger   *log.Logger
    file     *os.File
    config   *Config
    logLevel int
    mutex    sync.Mutex
    flagTime string
}

type Config struct {
    FileName   string `json:"file_name"`
    MaxSize    int    `json:"max_size"`
    ExpireDays int    `json:"expire_days"`
    LogLevel   string `json:"log_level"`
    TotalSize  int64  `json:"total_size"`
    StdOutput  bool   `json:"std_output"`
}

// LoadLogConfig initializes Logger struct.
// Load file for writing logs.
// Execute file splitting concurrently.
// Execute file clearing concurrently.
func LoadLogConfig(conf Config) error {
    l = new(Logger)
    l.config = &Config{
        FileName:   conf.FileName,
        MaxSize:    conf.MaxSize,
        ExpireDays: conf.ExpireDays,
        LogLevel:   conf.LogLevel,
        TotalSize:  conf.TotalSize,
        StdOutput:  conf.StdOutput,
    }
    if l.config.MaxSize <= 0 {
        l.config.MaxSize = DefaultMaxSize
    }
    // init flagTime
    l.flagTime = time.Now().Format(FlagTimeFmt)
    // init log level
    SetLogLevel(l.config.LogLevel)
    // std out
    if l.config.StdOutput {
        l.logger = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)
        return nil
    }
    // mkdir
    var err error
    var dir string
    dir, err = filepath.Abs(filepath.Dir(l.config.FileName))
    if err != nil {
        return err
    }
    err = os.MkdirAll(dir, os.ModePerm)
    if err != nil {
        return err
    }
    // open log file
    l.file, err = os.OpenFile(l.config.FileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        return err
    }
    // init logger
    l.logger = log.New(l.file, "", log.LstdFlags|log.Lmicroseconds)
    // split log file
    onceSplit.Do(func() {
        go splitLogFile()
    })
    // clear old files
    if l.config.ExpireDays > 0 {
        onceClear.Do(func() {
            go clearLogFile()
        })
    }
    // clear files when total size is too big
    if l.config.TotalSize > 0 {
        onceClearByTotalSizes.Do(func() {
            go clearLogFilesByTotalSize()
        })
    }
    return nil
}

// WriteLog is the real entry for writing log.
func writeLog(flag interface{}, detail ...interface{}) {
    detail = append([]interface{}{flag}, detail...)
    l.mutex.Lock()
    defer l.mutex.Unlock()
    // write log
    l.logger.Println(detail...)
}

// Fatal writes logs in fatal level.
func Fatal(detail ...interface{}) {
    if l.logLevel < FatalLevel {
        return
    }
    writeLog(FatalFlag, detail...)
}

// Error writes logs in error level.
func Error(detail ...interface{}) {
    if l.logLevel < ErrorLevel {
        return
    }
    writeLog(ErrorFlag, detail...)
}

// Warn writes logs in warn level.
func Warn(detail ...interface{}) {
    if l.logLevel < WarnLevel {
        return
    }
    writeLog(WarnFlag, detail...)
}

// Info writes logs in info level.
func Info(detail ...interface{}) {
    if l.logLevel < InfoLevel {
        return
    }
    writeLog(InfoFlag, detail...)
}

// Debug writes logs in debug level.
func Debug(detail ...interface{}) {
    if l.logLevel < DebugLevel {
        return
    }
    writeLog(DebugFlag, detail...)
}

// SetLogLevel sets log level for writing into file.
func SetLogLevel(level string) {
    switch level {
    case FatalLevelConf:
        l.logLevel = FatalLevel
    case ErrorLevelConf:
        l.logLevel = ErrorLevel
    case WarnLevelConf:
        l.logLevel = WarnLevel
    case InfoLevelConf:
        l.logLevel = InfoLevel
    case DebugLevelConf:
        l.logLevel = DebugLevel
    default:
        l.logLevel = DebugLevel
    }
}

// splitLogFile is used for splitting log file.
// Split file when time is next day.
// Split file when size is bigger than l.config.MaxSize.
func splitLogFile() {
    for {
        time.Sleep(time.Duration(1) * time.Second)
        nowTime := time.Now().Format(FlagTimeFmt)
        // split by time
        if l.flagTime != nowTime {
            l.mutex.Lock()
            fileBak := l.config.FileName + l.flagTime + strconv.Itoa(getFlagSplit())
            l.flagTime = nowTime
            _ = os.Rename(l.config.FileName, fileBak)
            l.file, _ = os.OpenFile(l.config.FileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
            l.logger.SetOutput(l.file)
            l.mutex.Unlock()
        } else {
            // split by size
            fileInfo, err := l.file.Stat()
            if err != nil {
                fmt.Println("get fileInfo of log file failed.", err)
                continue
            }
            if fileInfo.Size() >= int64(l.config.MaxSize)*1024*1024 {
                l.mutex.Lock()
                fileBak := l.config.FileName + l.flagTime + strconv.Itoa(getFlagSplit())
                _ = os.Rename(l.config.FileName, fileBak)
                l.file, _ = os.OpenFile(l.config.FileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
                l.logger.SetOutput(l.file)
                l.mutex.Unlock()
            }
        }
    }
}

// clearLogFile is used for clear expired log files.
// Clear logs before l.config.ExpireDays and delete DefaultCheckDays files.
func clearLogFile() {
    for {
        for i := l.config.ExpireDays + 1; i < l.config.ExpireDays+DefaultCheckDays+1; i++ {
            timeFmt := time.Now().Add(-time.Hour * time.Duration(24*i)).Format(FlagTimeFmt)
            fileTmp := l.config.FileName + timeFmt + "*"
            files, _ := filepath.Glob(fileTmp)
            for _, fileName := range files {
                os.Remove(fileName)
            }
        }
        time.Sleep(time.Duration(1) * time.Hour)
    }
}

// clearLogFileByTotalSize is used for clear log files when total size is too big.
// Delete oldest log files.
func clearLogFilesByTotalSize() {
    for {
        fileTmp := l.config.FileName + ".*"
        files, _ := filepath.Glob(fileTmp)
        fileInfoMap := make(map[string]os.FileInfo, 0)
        var totalSize int64 = 0
        for _, fileName := range files {
            fileInfo, err := os.Stat(fileName)
            if err != nil {
                continue
            }
            totalSize += fileInfo.Size()
            fileInfoMap[fileName] = fileInfo
        }
        for totalSize >= (l.config.TotalSize-(int64)(l.config.MaxSize))*1024*1024 && len(fileInfoMap) > 0 {
            key := ""
            timeMark := time.Now()
            for k, v := range fileInfoMap {
                if v.ModTime().UnixNano() < timeMark.UnixNano() {
                    timeMark = v.ModTime()
                    key = k
                }
            }
            totalSize -= fileInfoMap[key].Size()
            delete(fileInfoMap, key)
            os.Remove(key)
        }
        time.Sleep(time.Duration(1) * time.Minute)
    }
}

// getFlagSplit is used for get the suffix of the split log file
func getFlagSplit() int {
    suffix := 0
    prefix := l.config.FileName + l.flagTime
    fileTmp := prefix + "*"
    files, _ := filepath.Glob(fileTmp)
    for _, fileName := range files {
        a := strings.Split(fileName, l.flagTime)
        tmp, err := strconv.Atoi(a[len(a)-1])
        if err != nil {
            continue
        }
        if tmp >= suffix {
            suffix = tmp + 1
        }
    }
    return suffix
}
