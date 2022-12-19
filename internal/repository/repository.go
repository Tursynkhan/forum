package repository

type (
	Autorization interface{}
	Post         interface{}
	Comment      interface{}
)

type Repository struct {
	Autorization
	Post
	Comment
}

func NewRepository() *Repository {
	return &Repository{}
}
