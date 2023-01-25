package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

type Post interface {
	CreatePost(post models.Post) (int, error)
	GetAllPosts() ([]models.PostInfo, error)
	GetPost(id int) (models.PostInfo, error)
	CreatePostCategory(postId int, categories []string) error
	GetAllCategories() ([]models.Category, error)
	GetPostsByMostLikes() ([]models.PostInfo, error)
	GetPostsByLeastLikes() ([]models.PostInfo, error)
	GetPostsByNewest() ([]models.PostInfo, error)
	GetPostsByOldest() ([]models.PostInfo, error)
	GetPostByCategory(category string) ([]models.PostInfo, error)
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post models.Post) (int, error) {
	res, err := r.db.Exec("INSERT INTO posts (Title,Content,UserId,Created) VALUES (?,?,?,?)", post.Title, post.Content, post.UserID, post.Created)
	if err != nil {
		return 0, fmt.Errorf("repository : create post : %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *PostRepository) CreatePostCategory(postId int, categories []string) error {
	for _, category := range categories {
		_, err := r.db.Exec("INSERT INTO post_categories (PostId,CategoryId) VALUES (?,?)", postId, category)
		if err != nil {
			return fmt.Errorf("repository : create post : %w", err)
		}
	}
	return nil
}

func (r *PostRepository) GetAllPosts() ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created from posts JOIN users ON users.Id = posts.UserId")
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		likes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=1 AND PostId=?", p.ID)
		countlike := 0
		err = likes.Scan(&countlike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPost(id int) (models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created from posts JOIN users ON users.Id = posts.UserId WHERE posts.Id=$1", id)
	if err != nil {
		return models.PostInfo{}, fmt.Errorf("repository : get all posts: %w", err)
	}
	var post models.PostInfo
	for rows.Next() {
		err := rows.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.UserId, &post.Created)
		if err == sql.ErrNoRows {
			return models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &post.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			post.Categories = append(post.Categories, category)
		}
	}
	return post, nil
}

func (r *PostRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT Id,Name FROM categories")
	if err != nil {
		return []models.Category{}, fmt.Errorf("repository : GetAllCategories: %w", err)
	}
	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Category{}, err
		} else if err != nil {
			return []models.Category{}, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *PostRepository) GetPostsByMostLikes() ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id,posts.Title, posts.Content, posts.UserId, users.Username,posts.Created, (SELECT COUNT(*) FROM posts_like WHERE posts_like.PostId = posts.Id AND Status = 1) as likes FROM posts JOIN users on posts.UserId = users.Id ORDER BY likes DESC")
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByMostLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserId, &p.Author, &p.Created, &p.Likes)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		likes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=1 AND PostId=?", p.ID)
		countlike := 0
		err = likes.Scan(&countlike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByLeastLikes() ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id,posts.Title, posts.Content, posts.UserId, users.Username, posts.Created, (SELECT COUNT(*) FROM posts_like WHERE posts_like.PostId = posts.Id AND Status = 1) as likes FROM posts JOIN users on posts.UserId = users.Id ORDER BY likes ASC")
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByLeastLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserId, &p.Author, &p.Created, &p.Dislikes)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		likes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=1 AND PostId=?", p.ID)
		countlike := 0
		err = likes.Scan(&countlike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostByCategory(category string) ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created FROM posts INNER JOIN users ON posts.UserId=users.Id INNER JOIN post_categories ON posts.Id=post_categories.PostId INNER JOIN categories ON categories.Id=post_categories.CategoryId WHERE categories.Name=?", category)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostByCategory : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		likes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=1 AND PostId=?", p.ID)
		countlike := 0
		err = likes.Scan(&countlike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}

	return posts, nil
}

func (r *PostRepository) GetPostsByNewest() ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created from posts JOIN users ON users.Id = posts.UserId ORDER BY posts.Created DESC")
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		likes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=1 AND PostId=?", p.ID)
		countlike := 0
		err = likes.Scan(&countlike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByOldest() ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created from posts JOIN users ON users.Id = posts.UserId ORDER BY posts.Created ")
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		categories_rows, _ := r.db.Query("SELECT categories.Name FROM post_categories JOIN categories ON categories.Id = post_categories.CategoryId WHERE PostId = ?", &p.ID)
		for categories_rows.Next() {
			category := ""
			categories_rows.Scan(&category)
			p.Categories = append(p.Categories, category)
		}
		likes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=1 AND PostId=?", p.ID)
		countlike := 0
		err = likes.Scan(&countlike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetAllposts : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}
	return posts, nil
}
