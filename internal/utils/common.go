package utils

import (
	"fmt"
	"log"
	"os"
)

func GetEnvDefault(name, def string) string {
	val, ok := os.LookupEnv(name)
	if ok {
		return val
	}
	return def
}

func GetValue(v interface{}, def string) string {
	switch v.(type) {
	case string:
		if v != "" {
			return v.(string)
		}
	case *string:
		if *(v.(*string)) != "" {
			return *(v.(*string))
		}
	}
	return def
}

// Must halts the application if error is not nil
func Must(err error) {
	if err != nil {
		log.Fatalf("ERROR : %s\n", err)
	}
}

// Die halts the execution with a message
func Die(msg string) {
	Must(fmt.Errorf("%s", msg))
}
