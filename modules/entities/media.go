package entities

type Image struct {
	Id       string `db:"id"`
	FileName string `db:"filename" json:"filename"`
	Url      string `db:"url" json:"url"`
	
}
