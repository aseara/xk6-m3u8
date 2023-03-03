package m3u8

import "fmt"

func ReportError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}
