package helper

import (
	"os"

	"github.com/followgo/ND-Tester/public/errors"
)

// HasFile 判断文件存在
func HasFile(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return !fi.IsDir(), nil
}

// HasDir 判断目录存在
func HasDir(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return fi.IsDir(), nil
}
