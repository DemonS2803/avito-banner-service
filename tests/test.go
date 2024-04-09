package tests

import (
	"fmt"
	"strconv"
)

func main() {
	s := ""
	n, e := strconv.Atoi(s)
	if e != nil {
		fmt.Println("not parsed")
	} else {
		fmt.Println(n)
	}

}
