package reckontoken

import (
	"bufio"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

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
			// 将解码后的字节数组转换为字符串
			token := string(decoded)
			value := strings.TrimSpace(parts[1])
			// 可以转换 value 为 int
			tokens[token] = atoi(value)
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func atoi(str string) int {
	var result int
	for i := 0; i < len(str); i++ {
		result = result*10 + int(str[i]-'0')
	}
	return result
}
