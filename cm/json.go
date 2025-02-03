package cm

import (
	"encoding/json"
)

func JsonMarshal(v any) string {
	res, _ := json.Marshal(v)
	return string(res)
}
