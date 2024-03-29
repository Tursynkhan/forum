package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"strings"
)

type PostRepository struct {
	db *sql.DB
}

type Post interface {
	CreatePost(post models.Post) (int, error)
	GetAllPosts() ([]models.PostInfo, error)
	GetPost(id int) (models.PostInfo, error)
	CreatePostCategory(postId int, categories []string) error
	EditPostCategory(postId int, categories []string) error
	GetAllCategories() ([]models.Category, error)
	GetPostsByMostLikes() ([]models.PostInfo, error)
	GetPostsByLeastLikes() ([]models.PostInfo, error)
	GetPostsByNewest() ([]models.PostInfo, error)
	GetPostsByOldest() ([]models.PostInfo, error)
	GetPostByCategory(categoryId int) ([]models.PostInfo, error)
	GetLenAllPost() (int, error)
	SaveImageForPost(postId int, filePath string) error
	DeletePostById(postId int) error
	EditPost(newPost models.Post, postId int) error
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

func (r *PostRepository) EditPostCategory(postId int, categories []string) error {
	_, err := r.db.Exec("DELETE FROM post_categories WHERE PostId=?", postId)
	if err != nil {
		return fmt.Errorf("repo : EditPostCategory : %w", err)
	}
	for _, category := range categories {
		_, err := r.db.Exec("INSERT INTO post_categories (PostId,CategoryId) VALUES (?,?)", postId, category)
		if err != nil {
			return fmt.Errorf("repository : EditPostCategory : %w", err)
		}
	}
	return nil
}

func (r *PostRepository) GetAllPosts() ([]models.PostInfo, error) {
	rows, err := r.db.Query(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created
	HAVING Approved="YES";
	`)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : get all posts : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		var category string
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created, &p.Likes, &p.Dislikes, &category)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		p.Categories = strings.Split(category, ",")
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPost(id int) (models.PostInfo, error) {
	row := r.db.QueryRow(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created HAVING posts.Id=?;
	`, id)
	var post models.PostInfo
	var category string
	err := row.Scan(&post.ID, &post.Author, &post.Title, &post.Content, &post.UserId, &post.Created, &post.Likes, &post.Dislikes, &category)
	if errors.Is(err, sql.ErrNoRows) {
		return models.PostInfo{}, err
	} else if err != nil {
		return models.PostInfo{}, err
	}
	post.Categories = strings.Split(category, ",")
	images_rows, err := r.db.Query("SELECT images.ImageName FROM images JOIN posts ON images.PostId=posts.Id WHERE posts.Id=?", &post.ID)
	if err != nil {
		return models.PostInfo{}, fmt.Errorf("repository: GetPost: images: %w", err)
	}
	for images_rows.Next() {
		imageName := ""
		images_rows.Scan(&imageName)
		post.Images = append(post.Images, imageName)
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
	rows, err := r.db.Query(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created ORDER BY Likes DESC;
	`)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByMostLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		var category string
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created, &p.Likes, &p.Dislikes, &category)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		p.Categories = strings.Split(category, ",")
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByLeastLikes() ([]models.PostInfo, error) {
	rows, err := r.db.Query(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created ORDER BY Likes ASC;
	`)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByMostLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		var category string
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created, &p.Likes, &p.Dislikes, &category)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		p.Categories = strings.Split(category, ",")
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostByCategory(categoryId int) ([]models.PostInfo, error) {
	rows, err := r.db.Query(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
		WHERE categories.Id=?
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created ORDER BY Likes DESC;
	`, categoryId)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByMostLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		var category string
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created, &p.Likes, &p.Dislikes, &category)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		p.Categories = strings.Split(category, ",")
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByNewest() ([]models.PostInfo, error) {
	rows, err := r.db.Query(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created ORDER BY posts.Created DESC;
	`)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByMostLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		var category string
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created, &p.Likes, &p.Dislikes, &category)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		p.Categories = strings.Split(category, ",")
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetPostsByOldest() ([]models.PostInfo, error) {
	rows, err := r.db.Query(`SELECT posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created,COUNT(DISTINCT CASE WHEN posts_like.Status = 1 THEN posts_like.Id END) AS Likes,COUNT(DISTINCT CASE WHEN posts_like.Status = -1 THEN posts_like.Id END) AS Dislikes,GROUP_CONCAT(DISTINCT categories.Name) AS category
	FROM posts
		JOIN users ON users.Id = posts.UserId
		INNER JOIN post_categories ON posts.Id = post_categories.PostId
		INNER JOIN categories ON categories.Id = post_categories.CategoryId
		LEFT JOIN posts_like ON posts.Id = posts_like.PostId
	GROUP BY posts.Id, users.Username, posts.Title, posts.Content, posts.UserId, posts.Created ORDER BY posts.Created ASC;
	`)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostsByMostLikes : %w", err)
	}
	var posts []models.PostInfo
	for rows.Next() {
		p := models.PostInfo{}
		var category string
		err := rows.Scan(&p.ID, &p.Author, &p.Title, &p.Content, &p.UserId, &p.Created, &p.Likes, &p.Dislikes, &category)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.PostInfo{}, errors.New("No posts")
		} else if err != nil {
			return []models.PostInfo{}, err
		}
		p.Categories = strings.Split(category, ",")
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *PostRepository) GetLenAllPost() (int, error) {
	row := r.db.QueryRow("SELECT COUNT (*) FROM posts")
	count := 0
	err := row.Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, err
		} else {
			return 0, fmt.Errorf("repository: GetLenAllPost : %w", err)
		}
	}
	return count, nil
}

func (r *PostRepository) SaveImageForPost(postId int, filePath string) error {
	stmt, err := r.db.Prepare("INSERT INTO images (ImageName,PostId) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(filePath, postId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) DeletePostById(postId int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE Id=?", postId)
	if err != nil {
		return fmt.Errorf("repository: DeletePostById: %w", err)
	}
	return nil
}

func (r *PostRepository) EditPost(newPost models.Post, postId int) error {
	_, err := r.db.Exec("UPDATE posts SET title=?,content=? WHERE Id=?", newPost.Title, newPost.Content, postId)
	if err != nil {
		return fmt.Errorf("repo : EditPost : %w", err)
	}
	return nil
}
