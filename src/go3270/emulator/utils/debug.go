package utils

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(fn any) (pkg string, nm string) {
	var pc uintptr
	// 👇 for a specific function
	if fn != nil {
		v := reflect.ValueOf(fn)
		if v.Kind() != reflect.Func {
			return "", ""
		}
		pc = v.Pointer()
		// 👇 for the caller
	} else {
		pc, _, _, _ = runtime.Caller(2)
	}
	// 👇 get the function object for the PC
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "", ""
	}
	full := f.Name()
	// 👇 eg: "net/http.(*Server).Serve" or "main.myFunc"
	slashParts := strings.Split(full, "/")
	lastPart := slashParts[len(slashParts)-1]
	// Split remaining by "."
	dotParts := strings.Split(lastPart, ".")
	if len(dotParts) == 1 {
		// no package
		return "", dotParts[0]
	}
	// 👇 put it all together
	pkg = strings.Join(dotParts[:len(dotParts)-1], ".")
	nm = dotParts[len(dotParts)-1]
	return pkg, nm
}
