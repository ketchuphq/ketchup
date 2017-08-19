package enginebase

import (
	"html/template"
	"time"

	"github.com/Masterminds/sprig"
)

var sprigWhitelist = []string{
	"date",
	"dateModify",
	"dateInZone",
	"now",
}

func FuncMap() template.FuncMap {
	funcMap := map[string]interface{}{}
	genericFuncs := sprig.GenericFuncMap()
	for _, k := range sprigWhitelist {
		funcMap[k] = genericFuncs[k]
	}
	funcMap["dateParseMillis"] = func(millis int64) time.Time {
		return time.Unix(millis/1000, 0)
	}
	return funcMap
}
