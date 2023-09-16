package util_fn

import "encoding/json"

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func JsonString(v any) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}
