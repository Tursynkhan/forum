package models

import "mime/multipart"

type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	Created      string
	Approved     string
	ReportStatus string
	Files        []*multipart.FileHeader
}
