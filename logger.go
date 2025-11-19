package gdrive

import "fmt"

// Logger 日志接口，用於備份過程的日志輸出
type Logger interface {
	// Infof 信息級別日志
	Infof(format string, v ...interface{})

	// Warningf 警告級別日志
	Warningf(format string, v ...interface{})

	// Errorf 錯誤級別日志
	Errorf(format string, v ...interface{})
}

// defaultLogger 默認日志實現（使用 fmt.Printf）
type defaultLogger struct{}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *defaultLogger) Warningf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// newDefaultLogger 創建默認日志實例
func newDefaultLogger() Logger {
	return &defaultLogger{}
}
