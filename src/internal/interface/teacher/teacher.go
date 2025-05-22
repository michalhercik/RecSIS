package teacher

type Teacher struct {
	SISID       any    `json:"KOD"`
	LastName    string `json:"PRIJMENI"`
	FirstName   string `json:"JMENO"`
	TitleBefore string `json:"TITULPRED"`
	TitleAfter  string `json:"TITULZA"`
}
