package categories

type PcDigaMenuItem struct {
	Typename    string           `json:"__typename"`
	Children    []PcDigaMenuItem `json:"childrens"`
	Name        string           `json:"name"`
	Type        string           `json:"type"`
	URLPath     string           `json:"url_path"`
	Highlight   int              `json:"highlight"`
	Content     any              `json:"content"`
	Columns     int              `json:"columns"`
	WhichColumn int              `json:"which_column"`
}
