package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cockroachdb/errors"
)

func GetHeadLogFile() *os.File {
	headLogFile, err := os.OpenFile("logs/head.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if os.IsNotExist(err) {
		headLogFile, err = os.Create("logs/head.log")
		if err != nil {
			panic(errors.Wrap(err, "failed to create head log file"))
		}
		return headLogFile
	}
	if err != nil {
		panic(errors.Wrap(err, "failed to open head log file"))
	}
	return headLogFile
}

func GetHeadLogFileSize() (int64, error) {
	info, err := os.Stat("logs/head.log")
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func RotateHeadLogFile() (retErr error) {
	currentDate := time.Now()

	topFile, err := GetTopStackLogFile()
	if err != nil {
		return err
	}
	if topFile == nil {
		topFileName := fmt.Sprintf("%d_%d_%d_%d.log", 0, currentDate.Year(), int(currentDate.Month()), currentDate.Day())
		if err := os.Rename("logs/head.log", "logs/"+topFileName); err != nil {
			return errors.Wrap(err, "failed to rename head.log")
		}
		return nil
	}
	defer func() {
		if err := topFile.Close(); err != nil && retErr == nil {
			retErr = errors.Wrap(err, "failed to close top file")
		}
	}()

	var index int
	var year, month, day int

	// Extract just the filename from the full path
	_, fileName := filepath.Split(topFile.Name())
	_, err = fmt.Sscanf(fileName, "%d_%d_%d_%d.log", &index, &year, &month, &day)
	if err != nil {
		return err
	}

	nextIndex := index + 1
	nextFileName := fmt.Sprintf("%d_%d_%d_%d.log", nextIndex, currentDate.Year(), int(currentDate.Month()), currentDate.Day())
	if err := os.Rename("logs/head.log", "logs/"+nextFileName); err != nil {
		return errors.Wrap(err, "failed to rename head.log")
	}
	return nil
}

func GetTopStackLogFile() (*os.File, error) {
	dirPath := "./logs"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}

	format := "%d_%d_%d_%d.log"
	topFileName := ""
	topFileIndex := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "head.log" {
			continue
		}

		var index int
		var day int
		var month int
		var year int

		_, err = fmt.Sscanf(file.Name(), format, &index, &year, &month, &day)
		if err != nil {
			continue
		}

		if topFileName == "" || index > topFileIndex {
			topFileIndex = index
			topFileName = file.Name()
		}
	}

	// If no valid log files found, return nil
	if topFileName == "" {
		return nil, nil
	}

	topFile, err := os.OpenFile(dirPath+"/"+topFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return topFile, nil
}

func BeforeLogFileWrite(incomingData []byte) error {
	// Lock to prevent concurrent rotation attempts
	rotationMutex.Lock()
	defer rotationMutex.Unlock()

	size, err := GetHeadLogFileSize()
	if err != nil {
		return errors.Wrap(err, "failed to get head log file size, log is not recorded into the file")
	}
	// 1MB
	// If size is less than 1MB, don't rotate
	if size+int64(len(incomingData)) <= 1024*1024 {
		return nil
	}
	// Do the rotation
	err = RotateHeadLogFile()
	if err != nil {
		return errors.Wrap(err, "failed to rotate head log file")
	}

	// Reopen the file handle to point to new head.log
	if rotatingFile != nil {
		err = rotatingFile.Reopen()
		if err != nil {
			return errors.Wrap(err, "failed to reopen log file after rotation")
		}
	}

	return nil
}
