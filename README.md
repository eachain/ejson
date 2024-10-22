# ejson

ejson提供了一种json的快速访问和生成的方法。



## 查询

```go
package main

import (
	"fmt"

	"github.com/eachain/ejson"
)

func main() {
	fmt.Println(ejson.FromString(`null`).IsNull())
	// Output:
	// true

	fmt.Println(ejson.FromString(`true`).Bool())
	// Output:
	// true

	fmt.Println(ejson.FromString(`123`).Int())
	// Output:
	// 123

	fmt.Println(ejson.FromString(`[456]`).Get("[0]").Int())
	// Output:
	// 456

	fmt.Println(ejson.FromString(`{"a":789}`).Get("a").Int())
	// Output:
	// 789

	fmt.Println(ejson.FromString(`[{"x": {"y": ["eachain"]}}]`).Get("[0].x.y[0]").Str())
	// Output:
	// eachain

	fmt.Println(ejson.FromString(`[0]`).Get("a.b").Exists())
	// Output:
	// false

	fmt.Println(ejson.FromString(`{"b":999}`).Any("a", "b").Int())
	// Output:
	// 999

	fmt.Println(ejson.FromString(`{"c":3, "b":2, "a": 1}`).Keys())
	// Output:
	// [c b a]

	fmt.Println(ejson.FromString(`{"c":3, "b":2, "a": 1}`).Len())
	// Output:
	// 3
}
```



## 写入

```go
package main

import (
	"fmt"

	"github.com/eachain/ejson"
)

func main() {
	js := new(ejson.JSON)
	js.Get("a[0].b.c").Set(123)
	fmt.Println(js)
	// Output:
	// {"a":[{"b":{"c":123}}]}

	js = ejson.FromString(`{"a":123, "b":456}`)
	js.Get("a").Remove()
	fmt.Println(js)
	// Output:
	// {"b":456}

	js = ejson.FromString(`{"a":123, "b":{"x":true}}`)
	js.Get("a").Set(js.Get("b"))
	js.Get("b").Remove()
	fmt.Println(js)
	// Output:
	// {"a":{"x":true}}

	js = ejson.FromString(`{"a":123, "b":{"x":true}}`)
	b := js.Get("b")
	b.Remove()
	js.Get("a").Set(b)
	fmt.Println(js)
	// Output:
	// {"a":{"x":true}}
}
```



