package dbds

import (
	"encoding/json"
	"fmt"
)

type TeacherSlice []Teacher

func (ts *TeacherSlice) Scan(val interface{}) error {
	switch v := val.(type) {
	case nil:
		*ts = nil
		return nil
	case []byte:
		json.Unmarshal(v, &ts)
		return nil
	case string:
		json.Unmarshal([]byte(v), &ts)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

type Teacher struct {
	SisID       string `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	TitleBefore string `json:"title_before"`
	TitleAfter  string `json:"title_after"`
}
