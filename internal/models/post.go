package models

import "mime/multipart"

type Post struct {
	ID      int
	UserID  int
	Title   string
	Content string
	Created string
	Files   []*multipart.FileHeader
}
