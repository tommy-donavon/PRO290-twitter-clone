package data

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-playground/validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	Post struct {
		gorm.Model
		PostBody     string  `gorm:"not null" json:"post_body" validate:"required"`
		Author       string  `gorm:"not null" json:"author"`
		AuthorURI    string  `json:"author_uri"`
		ImageURI     string  `json:"image_uri"`
		Likes        int     `gorm:"default:0" json:"num_likes"`
		Reported     bool    `gorm:"default:0" json:"is_reported"`
		RefersToPost *uint   `json:"refers_to_post"`
		Comments     []*Post `gorm:"foreignkey:refers_to_post" json:"comments"`
	}

	FollowInformation struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	PostRepo struct {
		db *gorm.DB
	}
)

func NewPostRepo() *PostRepo {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Post{})
	return &PostRepo{db}
}

func (pr *PostRepo) CreatePost(post *Post, id uint) error {
	if id != 0 {
		p := pr.GetPost(id)
		if p.ID != 0 {
			post.RefersToPost = &id
			return pr.db.Create(post).Error
		}
	}
	return pr.db.Create(post).Error
}

func (pr *PostRepo) GetAllPosts() []*Post {
	posts := []*Post{}
	pr.db.Preload("Comments").Find(&posts)
	return posts
}

func (pr *PostRepo) GetPost(id uint) *Post {
	p := Post{}
	pr.db.Preload("Comments").Where("id = ?", id).First(&p)
	return &p
}

func (pr *PostRepo) DeletePost(id uint) error {
	return pr.db.Delete(&Post{}, id).Error
}

func (pr *PostRepo) GetFeed(username string, following []*FollowInformation) []*Post {
	feed := []*Post{}
	usernames := []string{}
	for _, val := range following {
		usernames = append(usernames, val.Username)
	}
	pr.db.Preload("Comments").Where("(author = ? OR author in ?) AND refers_to_post is null", username, usernames).Order("created_at DESC").Find(&feed)
	return feed
}
func (pr *PostRepo) LikePost(id uint) error {
	post := Post{}
	err := pr.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return err
	}
	post.Likes++
	return pr.db.Save(&post).Error
}

func (pr *PostRepo) UnlikePost(id uint) error {
	post := Post{}
	err := pr.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return err
	}
	post.Likes--
	return pr.db.Save(&post).Error
}

func (ur *PostRepo) UpdatePost(id uint, updateinfo map[string]string) error {
	post := ur.GetPost(id)
	if post.Author == "" {
		return fmt.Errorf("no post found")
	}
	postBytes, err := json.Marshal(post)
	if err != nil {
		return err
	}
	postMap := map[string]interface{}{}
	err = json.Unmarshal(postBytes, &postMap)
	if err != nil {
		return err
	}
	for key, value := range updateinfo {
		if _, ok := postMap[key]; ok {
			switch key {
			case "author_uri":
				post.AuthorURI = value
			default:
				return fmt.Errorf("provided field %s is not updateable", key)
			}
		}
	}
	if err := post.Validate(); err != nil {
		return err
	}
	return ur.db.Save(&post).Error

}

func (ur *PostRepo) UpdateAllAuthorUri(authorName, uri string) error {
	posts := []*Post{}
	err := ur.db.Where("author = ?").Find(&posts).Error
	if err != nil {
		return err
	}
	for _, value := range posts {
		value.AuthorURI = uri
	}
	return ur.db.Save(posts).Error
}

func (ur *PostRepo) GetUsersPostFeed(username string) ([]*Post, error) {
	posts := []*Post{}
	err := ur.db.Where("author = ?", username).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, err

}

func (p *Post) Validate() error {
	validator := validator.New()
	return validator.Struct(p)
}
