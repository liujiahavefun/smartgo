package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	PER_QUARTER = iota
	HALF_HOUR
	PER_HOUR
	PER_DAY
)

func GetFileNameAndNextTime(policy string, prefix string) (filename string, nexttime time.Time, err error) {
	var unit int64 = 60 * 60
	var format string = "2006-01-02_15:04:05"
	switch policy {
	case "perminute":
		unit = 60
	case "perquarter":
		unit = 15 * 60
	case "halfhour":
		unit = 30 * 60
	case "perhour":
		unit = 60 * 60
	case "perday":
		unit = 24 * 60 * 60
	}

	now := time.Now().Unix()
	now = now - (now % unit)
	next := now + unit

	//fmt.Println("now: ", time.Unix(now, 0).Format("2006-01-02 15:04:05"))
	//fmt.Println("next: ", time.Unix(next, 0).Format("2006-01-02 15:04:05"))

	filename = prefix + "_" + time.Unix(now, 0).Format(format) + ".log"

	return filename, time.Unix(next, 0), nil
}

func PathJoin(path string, fn string) (absPath string, err error) {
	if abs := filepath.IsAbs(path); abs == true {
		return filepath.Join(path, fn), nil
	}

	if wd, err := os.Getwd(); err == nil {
		return filepath.Join(wd, path, fn), nil
	}

	return "", err
}

func IsDirExist(path string) (bool, error) {
	fi, err := os.Stat(path) //path存在？
	if err == nil {
		if fi.IsDir() {
			return true, nil
		} else {
			return true, fmt.Errorf("path \"%s\" exist, but not directory", path)
		}
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func IsFileExist(path string) (bool, error) {
	fi, err := os.Stat(path) //path存在？
	if err == nil {
		if fi.IsDir() {
			return true, fmt.Errorf("path \"%s\" exist, but not file", path)
		} else {
			return true, nil
		}
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
