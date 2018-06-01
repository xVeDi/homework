package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

type JSONDataFS struct {
	browsers [10][]byte `json:"browsers"`
	email1   []byte
	email2   []byte
	name     []byte `json:"name"`
	bcount   int    // счетчик браузеров, поскольку используется пул для структуры
}

func (v *JSONDataFS) UnmarshalJSON2(data []byte) {

	var splitN, field, st int
	v.bcount = 0
	for i := 0; i < len(data); i++ {
		// `]` = 93 - граница Browsers и остальных данных
		if data[i] == 44 && data[i-1] == 93 {
			splitN = i - 1
		}

		// `,"` = 44 34 - граница полей
		// последнее поле не нужно, обработчика дла него нет
		if splitN > 0 {
			if data[i] == 34 && data[i-1] == 44 {
				switch field {
				case 3:

					// 3 поле - email - разобрать на 2 части - до и после @
					// чтобы не применять replace
					for ii := (st + 9); ii < (i - 2); ii++ {

						if data[ii] == 64 {
							v.email1 = data[st+9 : ii]
							v.email2 = data[ii+1 : i-2]
							break
						}
					}

				case 5:
					v.name = data[st+8 : i-2]
				}
				st = i
				field++
			}
		}
	}

	// записать браузеры

	st = 13
	for i := 13; i < splitN; i++ {
		if data[i] == 44 && data[i-1] == 34 && data[i+1] == 34 {
			v.browsers[v.bcount] = make([]byte, len(data[st+1:i-1]), len(data[st+1:i-1]))
			copy(v.browsers[v.bcount], data[st+1:i-1])
			v.bcount++
			st = i + 1
		}
		if i == (splitN - 1) {
			v.browsers[v.bcount] = make([]byte, len(data[st+1:i]), len(data[st+1:i]))
			copy(v.browsers[v.bcount], data[st+1:i])
			v.bcount++
		}
	}

}

var pool = sync.Pool{
	New: func() interface{} {
		return &JSONDataFS{}

	},
}

// вам надо написать более быструю оптимальную этой функции
func SuperFastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	//seenBrowsers := map[string]struct{}{}
	seenBrowsers := [][]byte{}

	var isAndroid bool
	var isMSIE bool
	var ii int
	// var email string
	var browser []byte

	var sliceAndroid = []byte{65, 110, 100, 114, 111, 105, 100}
	var sliceMSIE = []byte{77, 83, 73, 69}
	//var isSeen bool

	fmt.Fprintln(out, "found users:")

	for i := 0; fileScanner.Scan(); i++ {

		user := pool.Get().(*JSONDataFS)

		user.UnmarshalJSON2(fileScanner.Bytes())

		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false
		//isSeen := false

		for i = 0; i < user.bcount; i++ {

			browser = user.browsers[i]

			if bytes.Contains(browser, sliceAndroid) {
				isAndroid = true
				isSeen := false
				for _, val := range seenBrowsers {
					if bytes.Equal(val, browser) {
						isSeen = true
						break
					}
				}
				if !isSeen{
					seenBrowsers = append(seenBrowsers, browser)
				}
			}

			if bytes.Contains(browser, sliceMSIE) {
				isMSIE = true
				isSeen := false
				for _, val := range seenBrowsers {
					if bytes.Equal(val, browser) {
						isSeen = true
						break
					}

				}
				if !isSeen {
					seenBrowsers = append(seenBrowsers, browser)
				}
			}

		}
		if isAndroid && isMSIE {
			//email = strings.Replace(user.email, "@", " [at] ", 1)
			fmt.Fprintf(out, "[%d] %s <%s [at] %s>\n", ii, user.name, user.email1, user.email2)
		}
		// user = &JSONDataFS{}
		pool.Put(user)
		ii++
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))

}
