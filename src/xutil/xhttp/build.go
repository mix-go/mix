package xhttp

import (
	"encoding/json"
	"net/url"
)

func BuildJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func BuildQuery(m map[string]string) string {
	values := &url.Values{}
	for k, v := range m {
		values.Add(k, v)
	}
	return values.Encode()
}
