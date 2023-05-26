package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

type User interface {
	GetPostByUsername(username string) ([]models.PostInfo, error)
	GetLikedPostByUsername(usernaem string) ([]models.PostInfo, error)
	GetCommentedPostByUsername(username string) ([]models.PostInfo, error)
	GetProfileByUsername(username string) (models.ProfileUser, error)
	GetAllRoles() ([]models.Role, error)
	CreateCategory(category string) error
	GetNameCategoryById(categoryId int) (string, error)
	DeleteCategoryByName(category string) error
	GetRoleIdByName(nameRole string) (int, error)
	UpdateUserRole(username string, roleId int) error
}

func (r *UserRepository) GetPostByUsername(username string) ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created from posts JOIN users ON users.Id = posts.UserId Where users.Username=?", username)
	if err != nil {
		return []models.PostInfo{}, fmt.Errorf("repository : GetPostByUsername : %w", err)
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
			return []models.PostInfo{}, fmt.Errorf("GetPostByUsername : GetAllLikesByPostId : %w", err)
		}

		dislikes := r.db.QueryRow("SELECT COUNT(*) FROM posts_like WHERE Status=-1 AND PostId=?", p.ID)
		countdislike := 0
		err = dislikes.Scan(&countdislike)
		if err != nil {
			return []models.PostInfo{}, fmt.Errorf("GetPostByUsername : GetAllLikesByPostId : %w", err)
		}
		p.Likes = countlike
		p.Dislikes = countdislike
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *UserRepository) GetLikedPostByUsername(username string) ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created FROM posts JOIN posts_like ON posts_like.PostId=posts.Id JOIN users on users.Id=posts_like.UserId WHERE users.Username=? AND posts_like.Status=1", username)
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

func (r *UserRepository) GetCommentedPostByUsername(username string) ([]models.PostInfo, error) {
	rows, err := r.db.Query("SELECT posts.Id, users.Username, posts.Title, posts.Content,posts.UserId,posts.Created FROM posts JOIN comments ON comments.PostId=posts.Id JOIN users on users.Id=comments.UserId WHERE users.Username=?", username)
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

func (r *UserRepository) GetProfileByUsername(username string) (models.ProfileUser, error) {
	row := r.db.QueryRow("SELECT username,email,(SELECT count(*) FROM posts JOIN posts_like ON posts_like.PostId=posts.Id JOIN users on users.Id=posts_like.UserId WHERE users.Username=? AND posts_like.Status=1) as countoflikedpost,(SELECT COUNT(*) from posts JOIN users ON users.Id = posts.UserId Where users.Username=?) as countofpost,(SELECT count(*) FROM posts JOIN comments ON comments.PostId=posts.Id JOIN users on users.Id=comments.UserId WHERE users.Username=?) as countofcomment from users WHERE username=?", username, username, username, username)
	var user models.ProfileUser
	err := row.Scan(&user.Username, &user.Email, &user.CountOfLikes, &user.CountOfPosts, &user.CountOfComments)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProfileUser{}, fmt.Errorf("user : GetUuserByUsername : %w", err)
		} else {
			return models.ProfileUser{}, err
		}
	}
	return user, nil
}

func (r *UserRepository) GetAllRoles() ([]models.Role, error) {
	rows, err := r.db.Query("SELECT Id,Name FROM roles")
	if err != nil {
		return []models.Role{}, fmt.Errorf("repository : GetAllRoles: %w", err)
	}
	var roles []models.Role
	for rows.Next() {
		var r models.Role
		err := rows.Scan(&r.ID, &r.Name)
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Role{}, err
		} else if err != nil {
			return []models.Role{}, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}
func (r *UserRepository) CreateCategory(category string) error {
	_, err := r.db.Exec("INSERT OR IGNORE INTO categories(Name) VALUES (?) ", category)
	if err != nil {
		return fmt.Errorf("repo : CreateCategory : %w", err)
	}
	return nil
}
func (r *UserRepository) DeleteCategoryByName(category string) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE Name=?", category)
	if err != nil {
		return fmt.Errorf("repo : DeleteCategoryByName : %w", err)
	}
	return nil
}
func (r *UserRepository) GetRoleIdByName(nameRole string) (int, error) {
	row := r.db.QueryRow("SELECT Id FROM roles WHERE Name=?", nameRole)
	var roleId int
	err := row.Scan(&roleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user : GetRoleIdByName : %w", err)
		} else {
			return 0, err
		}
	}
	return roleId, nil
}
func (r *UserRepository) UpdateUserRole(username string, roleId int) error {
	_, err := r.db.Exec("UPDATE users SET RoleId=? WHERE Username=?", roleId, username)
	if err != nil {
		return fmt.Errorf("repo : UpdateUserRole : %w", err)
	}
	return nil
}

func (r *UserRepository) GetNameCategoryById(categoryId int) (string, error) {
	row := r.db.QueryRow("SELECT Name FROM categories WHERE Id=?", categoryId)
	var category string
	err := row.Scan(&category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user : GetRoleIdByName : %w", err)
		} else {
			return "", err
		}
	}
	return category, nil
}
