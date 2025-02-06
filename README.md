# ⏳ reckontoken

reckontoken是使用go语言实现的BPE分词器，主要用于LLM模型调用中token的统计。

```python
package test

import (
	"encoding/base64"
	"fmt"
	"github.com/taweizhong/reckontoken"
	"strings"
	"testing"
)

func main() {
	encoder := reckontoken.GetEncoder("cl100k_base")
	text := "hello world"
	encodedTokens := encoder.EncodeOrdinary(text)
	fmt.Println(len(encodedTokens))
}
```

reckontoken可以直接安装使用：

```
go get github.com/taweizhong/reckontoken
```

