package log 

import (
    "fmt"
    "github.com/fsnotify/fsnotify"
    "log"
    "os"
    "runtime"
    "strings"
)

var (
    logger *log.Logger
    file   *os.File
)

// Init logger with a local file path. Open and write log file
func InitLogger(path string) {
    if path == "" {
        logger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
        return
    }   
    file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
    if err != nil {
        panic(err)
    }   
    logger = log.New(file, "", log.Ldate|log.Lmicroseconds)
    go fileMonitor(path)
}

// Debug log
func Debug(tid, format string, args ...interface{}) {
    writeFile("DEBG", tid, format, args...)
}

// Trace log
func Trace(tid, format string, args ...interface{}) {
    writeFile("TRAC", tid, format, args...)
}

// Information log
func Info(tid, format string, args ...interface{}) {
    writeFile("INFO", tid, format, args...)
}

// Warning log
func Warn(tid, format string, args ...interface{}) {
    writeFile("WARN", tid, format, args...)
}

// Error log
func Error(tid, format string, args ...interface{}) {
    writeFile("EROR", tid, format, args...)
}

// Critical log
func Critical(tid, format string, args ...interface{}) {
    writeFile("CRIT", tid, format, args...)
}

// Close log file
func Close() {
    if file != nil {
        file.Close()
    }
}

// Write message to log file as a line
func writeFile(level, tid, format string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(2)
    format = fmt.Sprintf("%s:%d: [%s] [%s] ", strings.Split(file, "/src/")[1], line, level, tid) + format
    logger.Printf(format, args...)
}

// setOutput sets the local path as log file
func setOutput(path string) {
    f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
    if err != nil {
        panic(err)
    }
    logger.SetOutput(f)
    Close()
    file = f
}

// monitor the local path
func fileMonitor(path string) {
    dir := path[0:strings.LastIndex(path, "/")]
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    done := make(chan struct{})
    defer func() {
        watcher.Close()
        close(done)
    }()
    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if (event.Op&fsnotify.Rename == fsnotify.Rename ||
                    event.Op&fsnotify.Remove == fsnotify.Remove) &&
                    (event.Name == path) {
                    log.Println("file has been removed or renamed", path)
                    setOutput(path)
                }
            case err := <-watcher.Errors:
                log.Println("watcher error:", err)
                done <- struct{}{}
                return
            }
        }
    }()
     
     err = watcher.Add(dir)
    if err != nil {
        log.Fatal(err)
    }
    <-done
}
