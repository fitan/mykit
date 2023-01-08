package ormdata

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"column:name;notnull;comment:'姓名'" json:"name"`
	Username string `json:"username" gorm:"column:username;notnull;comment:'用户名'"`
	Email    string `json:"email" gorm:"column:email;notnull;comment:'邮箱'"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address" gorm:"serializer:json;column:address;notnull;comment:'地址'"`
	Phone   string `json:"phone" gorm:"column:phone;notnull;comment:'电话'"`
	Website string `json:"website" gorm:"column:website;notnull;comment:'网站'"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	} `json:"company" gorm:"serializer:json;column:company;notnull;comment:'公司'"`

	Posts  *[]Post  `json:"posts,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Albums *[]Album `json:"albums,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Todos  *[]Todo  `json:"todos,omitempty" gorm:"foreignKey:UserId;references:ID"`
}

type Post struct {
	gorm.Model
	UserId   int        `json:"userId" gorm:"column:user_id;notnull;comment:'用户id'"`
	Title    string     `json:"title" gorm:"column:title;notnull;comment:'标题'"`
	Body     string     `json:"body" gorm:"column:body;notnull;comment:'内容'"`
	User     *User      `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
	Comments *[]Comment `json:"comments,omitempty" gorm:"foreignKey:PostId;references:ID"`
}

type Comment struct {
	gorm.Model
	PostId int    `json:"postId" gorm:"column:post_id;notnull;comment:'帖子id'"`
	Name   string `json:"name" gorm:"column:name;notnull;comment:'姓名'"`
	Email  string `json:"email" gorm:"column:email;notnull;comment:'邮箱'"`
	Body   string `json:"body" gorm:"column:body;notnull;comment:'内容'"`
	Post   *Post  `json:"post,omitempty" gorm:"foreignKey:ID;references:PostId"`
}

type Album struct {
	gorm.Model
	UserId int      `json:"userId" gorm:"column:user_id;notnull;comment:'用户id'"`
	Title  string   `json:"title" gorm:"column:title;notnull;comment:'标题'"`
	User   *User    `json:"user" gorm:"foreignKey:UserId"`
	Photos *[]Photo `json:"photos,omitempty" gorm:"foreignKey:AlbumId;references:ID"`
}

type Photo struct {
	gorm.Model
	AlbumId      int    `json:"albumId" gorm:"column:album_id;notnull;comment:'相册id'"`
	Title        string `json:"title" gorm:"column:title;notnull;comment:'标题'"`
	Url          string `json:"url" gorm:"column:url;notnull;comment:'url'"`
	ThumbnailUrl string `json:"thumbnailUrl" gorm:"column:thumbnail_url;notnull;comment:'缩略图url'"`
	Album        *Album `json:"album,omitempty" gorm:"foreignKey:ID;references:AlbumId"`
}

type Todo struct {
	gorm.Model
	UserId    int    `json:"userId" gorm:"column:user_id;notnull;comment:'用户id'"`
	Title     string `json:"title" gorm:"column:title;notnull;comment:'标题'"`
	Completed bool   `json:"completed" gorm:"column:completed;notnull;comment:'是否完成'"`
	User      *User  `json:"user,omitempty" gorm:"foreignKey:UserId;references:ID"`
}

func (*Todo) Crud() {

}
