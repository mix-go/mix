package logic

import (
	"errors"
	"fmt"
	"github.com/mix-go/console"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyPath(src, dst string) bool {
	debug := console.App.Debug
	if debug {
		fmt.Println("")
	}

	src = strings.Replace(src, "\\", "/", -1)
	srcFileInfo := GetFileInfo(src)
	if srcFileInfo == nil || !srcFileInfo.IsDir() {
		return false
	}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.Replace(path, "\\", "/", -1)
		relationPath := strings.Replace(path, src, "", -1)
		dstPath := strings.TrimRight(strings.TrimRight(strings.Replace(dst, "\\", "/", -1), "/"), "\\") + relationPath

		if debug {
			fmt.Println(fmt.Sprintf(" Copy %s to %s", path, dstPath))
		}

		if !info.IsDir() {
			if CopyFile(path, dstPath) {
				return nil
			} else {
				return errors.New(path + " copy fail")
			}
		} else {
			if _, err := os.Stat(dstPath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
						return err
					} else {
						return nil
					}
				} else {
					return err
				}
			} else {
				return nil
			}
		}
	})

	if err != nil {
		return false
	}
	return true
}

func CopyFile(src, dst string) bool {
	if len(src) == 0 || len(dst) == 0 {
		return false
	}
	src = strings.Replace(src, "\\", "/", -1)
	srcFile, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	dst = strings.Replace(dst, "\\", "/", -1)
	dstPathArr := strings.Split(dst, "/")
	dstPathArr = dstPathArr[0 : len(dstPathArr)-1]
	dstPath := strings.Join(dstPathArr, "/")

	dstFileInfo := GetFileInfo(dstPath)
	if dstFileInfo == nil {
		if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
			return false
		}
	}

	//这里要把O_TRUNC 加上，否则会出现新旧文件内容出现重叠现象
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return false
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return false
	} else {
		return true
	}
}

func GetFileInfo(src string) os.FileInfo {
	if fileInfo, e := os.Stat(src); e != nil {
		if os.IsNotExist(e) {
			return nil
		}
		return nil
	} else {
		return fileInfo
	}
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// os.O_TRUNC 覆盖写入，不加则追加写入
func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, io.SeekEnd)
		_, err = f.WriteAt([]byte(content), n)
		defer f.Close()
	}
	return err
}
