package logger

import "fmt"

var hooks []func(string)

func Println(v ...interface{}) {
	fmt.Println(v...)
	for _, hook := range hooks {
		hook(fmt.Sprintln(v...))
	}
}

func AddHook(h func(string)) {
	hooks = append(hooks, h)
}
