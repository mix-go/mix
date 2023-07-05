package xhttp

import (
	"encoding/json"
	"net/url"
)

func BuildJSON(v interface{}) Body {
	b, _ := json.Marshal(v)
	return Body(b)
}

func BuildQuery(m map[string]string) Body {
	values := &url.Values{}
	for k, v := range m {
		values.Add(k, v)
	}
	return Body(values.Encode())
}
