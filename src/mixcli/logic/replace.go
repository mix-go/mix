package logic

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func ReplaceAll(root, old, new string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 替换内容
			text, err := ReadAll(path)
			if err != nil {
				return err
			}
			str := string(text)
			reg := regexp.MustCompile(old)
			str = reg.ReplaceAllString(str, new)
			if err := WriteToFile(path, str); err != nil {
				return err
			}
		}

		return err
	})
	return err
}

func ReplaceMod(root string) error {
	path := fmt.Sprintf("%s/go.mod", root)
	text, err := ReadAll(path)
	if err != nil {
		return err
	}
	str := string(text)
	reg := regexp.MustCompile(`(replace \([\s\S]*?\))`)
	str = reg.ReplaceAllString(str, "")
	if err := WriteToFile(path, str); err != nil {
		return err
	}
	return nil
}

func ReplaceMain(root, old, new string) error {
	path := fmt.Sprintf("%s/main.go", root)
	text, err := ReadAll(path)
	if err != nil {
		return err
	}
	str := string(text)
	reg := regexp.MustCompile(old)
	str = reg.ReplaceAllString(str, new)
	if err := WriteToFile(path, str); err != nil {
		return err
	}
	return nil
}
