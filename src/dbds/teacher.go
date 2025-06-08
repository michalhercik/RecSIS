package dbds

import (
	"encoding/json"
	"fmt"
)

type TeacherSlice []Teacher

func (ts *TeacherSlice) Scan(val interface{}) error {
	switch v := val.(type) {
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
	SisID       string `json:"KOD"`
	LastName    string `json:"PRIJMENI"`
	FirstName   string `json:"JMENO"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}
