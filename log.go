package fastgo

import "log"

// 获取log记录器
func GetLogger() *log.Logger {
	return log.Default()
}
