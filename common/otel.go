package common

import "os"

var ModuleName = "code-go"

func init() {
	m := os.Getenv("OTEL-ModuleName")
	if len(m) != 0 {
		ModuleName = m
	}
}
