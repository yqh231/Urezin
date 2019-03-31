package main

import "encoding/json"
import "fmt"

func test(v interface{}) {
	json.Unmarshal([]byte(`{"test": 4}`), v)
}

func test2(v *string) {
	//a := "abc"
	*v = "abc"
}

type ruby struct {
	lang string
}

func test3(v *ruby) {
	n := ruby{
		"ruby",
	}
	v = &n
}

func main() {
	v1 := ruby{
		"en",
	}
	test3(&v1)
	fmt.Println(v1)
}
