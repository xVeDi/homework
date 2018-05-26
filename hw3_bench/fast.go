package main

import (
	"bufio"
	json "encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonAff6eb80DecodeCourseraHomeworkEasytest2(in *jlexer.Lexer, out *JSONData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "company":
			out.Company = string(in.String())
		case "country":
			out.Country = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "job":
			out.Job = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "phone":
			out.Phone = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonAff6eb80EncodeCourseraHomeworkEasytest2(out *jwriter.Writer, in JSONData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"company\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Company))
	}
	{
		const prefix string = ",\"country\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Country))
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"job\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Job))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"phone\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Phone))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JSONData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonAff6eb80EncodeCourseraHomeworkEasytest2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JSONData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonAff6eb80EncodeCourseraHomeworkEasytest2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JSONData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonAff6eb80DecodeCourseraHomeworkEasytest2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JSONData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonAff6eb80DecodeCourseraHomeworkEasytest2(l, v)
}

type JSONData struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	seenBrowsers := map[string]struct{}{}

	var isAndroid bool
	var isMSIE bool
	var email string

	var pool = sync.Pool{
		New: func() interface{} {
			return &JSONData{}

		},
	}
	fmt.Fprintln(out, "found users:")

	for i := 0; fileScanner.Scan(); i++ {

		user := pool.Get().(*JSONData)

		user.UnmarshalJSON(fileScanner.Bytes())
		if err != nil {
			panic(err)
		}

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {

			if strings.Contains(browser, "Android") {
				seenBrowsers[browser] = struct{}{}
				isAndroid = true

			}
			if strings.Contains(browser, "MSIE") {
				seenBrowsers[browser] = struct{}{}
				isMSIE = true

			}
		}
		if isAndroid && isMSIE {
			email = strings.Replace(user.Email, "@", " [at] ", 1)
			fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
		}

		pool.Put(user)
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))

}
