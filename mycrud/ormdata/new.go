package ormdata

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	//dsn := "spider_dev:spider_dev123@tcp(10.170.34.22:3307)/gteml?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:123456@tcp(172.29.107.199:3306)/gteml?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	//db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
	//	DisableForeignKeyConstraintWhenMigrating: true,
	//})
	if err != nil {
		panic(err)
	}

	//err = db.AutoMigrate(&User{}, &Todo{}, &Post{}, &Photo{}, &Comment{}, &Album{})
	//if err != nil {
	//	panic(err)
	//}

	//albums := make([]Album,0)
	//comments := make([]Comment,0)
	//photos := make([]Photo,0)
	//posts := make([]Post,0)
	//todos := make([]Todo,0)
	//users := make([]User,0)

	//albumsB, err := ioutil.ReadFile("./ormdata/albums.json")
	//if err != nil {
	//	panic(err)
	//}
	//err = json.Unmarshal(albumsB, &albums)
	//
	//commentsB, err := ioutil.ReadFile("./ormdata/comments.json")
	//if err != nil {
	//	panic(err)
	//}
	//err = json.Unmarshal(commentsB, &comments)
	//
	//photosB, err := ioutil.ReadFile("./ormdata/photos.json")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(photosB, &photos)
	//if err !=nil {
	//	panic(err)
	//}
	//
	//postsB, err := ioutil.ReadFile("./ormdata/posts.json")
	//if err != nil {
	//	panic(err)
	//}
	//err = json.Unmarshal(postsB, &posts)
	//if err != nil {
	//	panic(err)
	//}
	//
	//todosB, err := ioutil.ReadFile("./ormdata/todos.json")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(todosB, &todos)
	//if err != nil {
	//	panic(err)
	//}
	//
	//usersB, err := ioutil.ReadFile("./ormdata/users.json")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = json.Unmarshal(usersB, &users)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = db.CreateInBatches(&albums,10).Error
	//if err != nil {
	//	panic(err)
	//}
	//err = db.CreateInBatches(&comments,10).Error
	//if err != nil {
	//	panic(err)
	//}
	//err = db.CreateInBatches(&photos, 10).Error
	//if err != nil {
	//	panic(err)
	//}
	//err = db.CreateInBatches(&posts,10).Error
	//if err != nil {
	//	panic(err)
	//}
	//err = db.CreateInBatches(&todos,10).Error
	//if err != nil {
	//	panic(err)
	//}
	//err = db.CreateInBatches(&users,10).Error
	//if err != nil {
	//	panic(err)
	//}

	return db
}
