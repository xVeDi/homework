package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func SingleHash(data string) string {
	md := DataSignerMd5(data)
	return DataSignerCrc32(data) + "~" + DataSignerCrc32(md)
}

func MultiHash(data string) string {

	var maps []string

	for i := 0; i < 6; i++ {

		// 4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
		maps = append(maps, DataSignerCrc32(strconv.Itoa(i)+data))
		// fmt.Println(data, " MultiHash: crc32(th+step1))", i, " ", maps[i])
	}

	return func(maps []string) (a string) {
		for _, v := range maps {
			a += v
		}
		return
	}(maps)
}

func CombineResults(in []string) (a string) {

	sort.Strings(in)

	for _, v := range in {
		a = a + "_" + v

	}
	return strings.TrimLeft(a, "_")
}

// func main() {
// 	// res := []string{"4958044192186797981418233587017209679042592862002427381542", "29568666068035183841425683795340791879727309630931025356555"}
// 	res := []string{}
// 	res = append(res, MultiHash(SingleHash("1")))
// 	res = append(res, MultiHash(SingleHash("0")))
// 	res = append(res, MultiHash(SingleHash("1")))
// 	res = append(res, MultiHash(SingleHash("2")))
// 	res = append(res, MultiHash(SingleHash("3")))
// 	res = append(res, MultiHash(SingleHash("5")))
// 	res = append(res, MultiHash(SingleHash("8")))

// 	fmt.Println(CombineResults(res))
// }
