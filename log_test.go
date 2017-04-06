package log 

import (
    "github.com/satori/go.uuid"
    "testing"
)

func TestLogger(t *testing.T) {
    InitLogger("./test.log")
    Debug(uuid.NewV1().String(), "this is %s", "debug")
    Trace(uuid.NewV1().String(), "this is %s", "trace")
    Info(uuid.NewV1().String(), "this is %s", "info")
    Warn(uuid.NewV1().String(), "this is %s", "warn")
    Error(uuid.NewV1().String(), "this is %s", "error")
    Critical(uuid.NewV1().String(), "this is %s", "critical")
    Close()
}
