package reckontoken

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// 获取当前包文件所在的目录
func getPackageDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}

// 读取词典
func LoadTokens(url string) (map[string]int, error) {
	response, _ := http.Get(url)
	tokens := make(map[string]int)
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(parts[0]))
			if err != nil {
				log.Fatal(err)
			}
			token := string(decoded)
			value := strings.TrimSpace(parts[1])
			tokens[token] = atoi(value)
		} else {
			log.Fatalf("file format error")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}
func LoadLocalTokens(path string) (map[string]int, error) {
	baseDir := getPackageDir()
	fullPath := filepath.Join(baseDir, path)

	file, err := os.Open(fullPath) // 使用 os.Open 而不是 os.ReadFile
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %w", fullPath, err)
	}
	defer file.Close() // 确保文件正确关闭
	tokens := make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(parts[0]))
			if err != nil {
				log.Fatal(err)
			}
			token := string(decoded)
			value := strings.TrimSpace(parts[1])
			tokens[token] = atoi(value)
		} else {
			log.Fatalf("file format error")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

// 转换
func atoi(str string) int {
	var result int
	for i := 0; i < len(str); i++ {
		result = result*10 + int(str[i]-'0')
	}
	return result
}
