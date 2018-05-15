package main

import (
	"strconv"
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

	return func(mps []string) (a string) {
		for _, v := range mps {
			a += v
		}
		return
	}(maps)
}

func main() {
	println(MultiHash(SingleHash("0")))
}
