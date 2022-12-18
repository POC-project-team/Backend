package postgres

import (
	response "backend/internal/dto/responseDto"
	u "backend/internal/entity"
	"backend/internal/repository/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Client struct {
	db *gorm.DB
}

func NewClient() (*Client, error) {
	client := &Client{}
	db, err := client.Connect()
	if err != nil {
		return nil, err
	}
	client.db = db
	return client, nil
}

func (c *Client) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=admin password=admin dbname=poc_db port=5432"))
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Tag{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Note{}); err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) GetUserID(login, password string) (int, error) {
	var user model.User
	if err := c.db.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		return 0, err
	}
	return int(user.UserID), nil
}

func (c *Client) GetAllUsers() ([]string, error) {
	var users []model.User
	if err := c.db.Find(&users).Error; err != nil {
		return nil, err
	}
	// return all id's of users
	var usersId []string
	for _, user := range users {
		usersId = append(usersId, fmt.Sprint(user.UserID))
	}
	return usersId, nil
}

func (c *Client) CreateUser(login, password string) (u.User, error) {
	var user model.User

	// create new user with unique id
	// put new id to user

	if err := c.db.Create(&model.User{Login: login, Password: password}).Error; err != nil {
		return u.User{
			UserID: 0,
		}, err
	}

	// add and return new user
	if err := c.db.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		return u.User{
			UserID: 0,
		}, err
	}

	return u.User{UserID: int(user.UserID)}, nil
}

func (c *Client) ChangePassword(userId int, password string) error {
	if err := c.db.Model(&model.User{}).Where("user_id = ?", userId).Update("password", password).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) ChangeLogin(userId int, login string) error {
	if err := c.db.Model(&model.User{}).Where("user_id = ?", userId).Update("login", login).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteUser(userId int) error {
	if err := c.db.Where("user_id = ?", userId).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}

//	func (c *Client) GetUserTags(userId int) ([]model.Tag, error) {
//		var tags []model.Tag
//		if err := c.db.Where("user_id = ?", userId).Find(&tags).Error; err != nil {
//			return nil, err
//		}
//		return tags, nil
//	}
func (c *Client) GetUserTags(userId int) ([]response.TagNoUserNotes, error) {
	var tags []model.Tag
	if err := c.db.Where("user_id = ?", userId).Find(&tags).Error; err != nil {
		return nil, err
	}
	var tagsNoUserNotes []response.TagNoUserNotes
	for _, tag := range tags {
		tagsNoUserNotes = append(tagsNoUserNotes, response.TagNoUserNotes{
			TagID:   tag.TagID,
			TagName: tag.TagName,
		})
	}
	return tagsNoUserNotes, nil
}

//func (c *Client) GetTag(userId int, tagName string) (model.Tag, error) {
//	var tag model.Tag
//	if err := c.db.Where("user_id = ? AND tag_name = ?", userId, tagName).First(&tag).Error; err != nil {
//		return tag, err
//	}
//	return tag, nil
//}

func (c *Client) GetTag(userId int, tagId string) (response.TagNoUserNotes, error) {
	var tag model.Tag

	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).First(&tag).Error; err != nil {
		return response.TagNoUserNotes{}, err
	}
	return response.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}, nil
}

func (c *Client) DeleteTag(userId int, tagId string) error {
	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).Delete(&model.Tag{}).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) GetUserNotes(userId int, tagId string) ([]u.Note, error) {
	var notes []u.Note
	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (c *Client) UpdateTag(userId int, tagId, tagName string) (response.TagNoUserNotes, error) {
	var tag response.TagNoUserNotes
	if err := c.db.Model(&model.Tag{}).Where("user_id = ? AND tag_id = ?", userId, tagId).Update("tag_name", tagName).Error; err != nil {
		return tag, err
	}
	return tag, nil
}

func (c *Client) CreateTag(userId int, tagId, tagName string) (response.TagNoUserNotes, error) {
	var tag model.Tag
	if err := c.db.Create(&model.Tag{UserID: uint(userId), TagID: tagId, TagName: tagName}).Error; err != nil {
		return response.TagNoUserNotes{}, err
	}
	// return the created tag
	if err := c.db.Where("user_id = ? AND tag_id = ?", userId, tagId).First(&tag).Error; err != nil {
		return response.TagNoUserNotes{}, err
	}
	return response.TagNoUserNotes{
		TagID:   tag.TagID,
		TagName: tag.TagName,
	}, nil
}

func (c *Client) TransferTag(userId int, tagId, login string) error {
	var user model.User
	if err := c.db.Where("login = ?", login).First(&user).Error; err != nil {
		return err
	}
	if err := c.db.Model(&model.Tag{}).Where("user_id = ? AND tag_id = ?", userId, tagId).Update("user_id", user.UserID).Error; err != nil {
		return err
	}
	return nil
}

func (c *Client) AddNote(userId int, tagId, noteInfo string) (u.Tag, error) {
	var tag u.Tag

	if err := c.db.Create(&model.Note{UserId: uint(userId), TagID: tagId, Note: noteInfo, Time: time.Now()}).Error; err != nil {
		return tag, err
	}
	return tag, nil
}
