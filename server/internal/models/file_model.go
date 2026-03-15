package models 

type File struct {
	ID string `bson:"_id,omitempty" json:"id"`
	Filename string `bson:"filename"`
	Path string `bson:"path"`
	Size int64 `bson:"size" `
	MimeType string `bson:"mime_type"`
	Hash string `bson:"hash"`
	CreatedAt int64 `bson:"created_at"`
}