package logic

import (
    "fmt"
    "os"
    "path/filepath"
    "regexp"
)

func ReplaceName(root, name string) error {
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
            reg := regexp.MustCompile(`github.com/mix-go/mix-skeleton/console`)
            str = reg.ReplaceAllString(str, name)
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
