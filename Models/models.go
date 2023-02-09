package Models

type Article struct {
	Id       int32  `json:"Id"`
	Title    string `json:"Title"`
	Desc     string `json:"Desc"`
	Content  string `json:"Content"`
	Category string `json:"Category"`
}

type Articles []Article
