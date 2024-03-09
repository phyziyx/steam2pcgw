package types

import "encoding/json"

type Requirement map[string]interface{}

func (req *Requirement) UnmarshalJSON(data []byte) error {
	stringified := string(data)
	if stringified == `''` ||
		stringified == `""` ||
		stringified == `{}` ||
		stringified == `[]` {
		return nil
	}

	type requirement Requirement
	return json.Unmarshal(data, (*requirement)(req))
}
