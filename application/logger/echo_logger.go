package logger

import (
	"github.com/labstack/echo/v4"
)

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// EchoLogger は `Logger`と同じメソッドを持つことでインターフェースを実装
type EchoLogger struct {
	echoLogger echo.Logger
}

// NewEchoLogger は `EchoLogger` のインスタンスを作成
func NewEchoLogger(e *echo.Echo) Logger {
	return &EchoLogger{echoLogger: e.Logger}
}

// `Info` メソッド（Echo の Logger を使用）
func (l *EchoLogger) Info(msg string, args ...interface{}) {
	l.echoLogger.Infof(msg, args...)
}

// `Error` メソッド（Echo の Logger を使用）
func (l *EchoLogger) Error(msg string, args ...interface{}) {
	l.echoLogger.Errorf(msg, args...)
}
