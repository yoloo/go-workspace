package util

import (
	"bytes"
	"encoding/json"
)

func JsonStringify(obj interface{}) string  {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(&obj); err != nil{
		return ""
	}
	return string(buf.Bytes())
}

