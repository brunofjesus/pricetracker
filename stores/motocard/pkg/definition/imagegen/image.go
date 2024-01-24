package imagegen

type Query struct {
	Bucket string     `json:"bucket"`
	Key    string     `json:"key"`
	Edits  QueryEdits `json:"edits"`
	V      string     `json:"v"`
}

type QueryEdits struct {
	Webp   QueryEditQuality `json:"webp"`
	Jpeg   QueryEditQuality `json:"jpeg"`
	Resize QueryEditResize  `json:"resize"`
}

type QueryEditQuality struct {
	Quality int `json:"quality"`
}

type QueryEditResize struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Fit    string `json:"fit"`
}
