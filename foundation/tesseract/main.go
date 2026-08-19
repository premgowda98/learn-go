package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/otiai10/gosseract/v2"
	"github.com/pkg/errors"
)

var MaxContentReadTimeout = time.Duration(15) * time.Second
var MaxImageSize int64 = 10 * 1024 * 1024 // 40MB

func main() {

	fPath := "./40mb.jpg"

	_, err := readContentWithTimeout(fPath)

	if err != nil {
		slog.Error(fmt.Sprintf("readContentWithTimeout failed:%v", err))
	}

	// slog.Info("Data", "content", data)

}

func readImageFile(fPath string) (string, error) {
	time.Sleep(5 * time.Second)
	slog.Info(fmt.Sprintf("reading fpath:%s, type:image", fPath))

	file, err := os.Open(fPath)
	if err != nil {
		slog.Error(fmt.Sprintf("readImageFile:open file error:%v", err))
		return "", errors.Wrapf(err, "failed to open image file:%s", fPath)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()

	if fileInfo.Size() > MaxImageSize {
		slog.Error(fmt.Sprintf("readImageFile:file size too large: %s", humanize.Bytes(uint64(fileInfo.Size()))))
		return "", errors.Wrapf(err, "file size too large:%s", fPath)
	}

	data := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(data)
	if err != nil {
		slog.Error(fmt.Sprintf("readImageFile:read file error:%v", err))
		return "", errors.Wrapf(err, "failed to read image file:%s", fPath)
	}

	client := gosseract.NewClient()
	defer client.Close()
	err = client.SetImageFromBytes(data)
	if err != nil {
		slog.Error(fmt.Sprintf("readImageFile:image bytes read error:%v", err))
		return "", errors.Wrapf(err, "failed to read image bytes:%s", fPath)
	}
	text, err := client.Text()
	if err != nil {
		slog.Error(fmt.Sprintf("readImageFile:text extract failed:%v", err))
		return "", errors.Wrapf(err, "failed to extract image text:%s", fPath)
	}

	gosseract.ClearPersistentCache()

	slog.Info(fmt.Sprintf("read fpath:%s, type:image", fPath))
	slog.Debug(fmt.Sprintf("readImageFile:%s, content:%s", fPath, text))

	return text, nil
}

type ContentResult struct {
	Content string
	Error   error
}

func readContent(fPath string) ContentResult {
	slog.Info(fmt.Sprintf("processing file:%s read request", fPath))

	var err error
	var result ContentResult

	// record time taken
	startTime := time.Now()
	fType := strings.TrimPrefix(filepath.Ext(fPath), ".")
	defer func(err error) {
		timeTaken := time.Since(startTime).Seconds()
		if err == nil {
			slog.Info(fmt.Sprintf("fpath:%s, type:%s, read took %.2fs", fPath, fType, timeTaken))
		} else {
			slog.Info(fmt.Sprintf("fpath:%s, type:%s, read took %.2fs with err=%v", fPath, fType, timeTaken, result.Error))
		}
	}(err)

	result.Content, err = readImageFile(fPath)
	if err != nil {
		result.Error = errors.Wrapf(err, "failed to read file:%s", fPath)
	}

	// read successful
	slog.Info(fmt.Sprintf("processed file:%s read request", fPath))
	return result
}

func readContentWithTimeout(fPath string) (*ContentResult, error) {
	slog.Info(fmt.Sprintf("processing file:%s read request with %s timeout", fPath, MaxContentReadTimeout.String()))

	ctx, cancel := context.WithTimeout(context.Background(), MaxContentReadTimeout)
	defer cancel() // Ensure resources are released

	done := make(chan ContentResult, 1)

	go func() {
		done <- readContent(fPath)
	}()

	select {
	case result := <-done:
		slog.Info(fmt.Sprintf("successfully read the content of fpath:%s", fPath))
		return &result, nil

	case <-ctx.Done():
		// handle error when context is cancelled or timeout
		slog.Warn(fmt.Sprintf("fpath:%s content read has timed out.", fPath))
		return &ContentResult{}, fmt.Errorf("timed out")
	}
}
