package golangutils

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"time"
)

func CreateLogWriter(logType string) (*bufio.Writer, func()) {
	// 获取当前时间
	now := time.Now()

	// 格式化日期和时间
	dateStr := now.Format("2006-01-02")
	timeStr := now.Format("15_04_05_000")
	filePath := filepath.Join(os.TempDir(), "logs", logType, dateStr, timeStr+".log")

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		log.Fatalf("创建目录失败：%s\nerror:%v", filePath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("创建文件失败：%s\nerror:%v", filePath, err)
	}

	writer := bufio.NewWriter(file)
	deferFunc := func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Fatalf("关闭文件失败：%s\nerror:%v", filePath, closeErr)
		}
	}
	return writer, deferFunc
}
