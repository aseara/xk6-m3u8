package redis

import "fmt"

func ReportError(err error, msg string) {
	if err != nil {
		fmt.Println(err, msg)
	}
}
