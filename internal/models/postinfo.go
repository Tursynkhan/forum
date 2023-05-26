package models

type PostInfo struct {
	ID           int
	Author       string
	Categories   []string
	Title        string
	Content      string
	Created      string
	Images       []string
	UserId       int
	Likes        int
	Dislikes     int
	Approved     string
	ReportStatus string
}
