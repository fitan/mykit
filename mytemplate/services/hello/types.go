package hello

type HelloRequest struct {
	ID    string `json:"id" param:"path,id"`
	Query Query  `json:"query" param:",query"`
}

type HelloResponse struct {
	ID string `json:"id"`
}

type Query struct {
	Age   string `json:"age" param:"query,age" query:"eq" gorm:"column:age"`
	Email string `json:"email" param:"query,email" query:"ne" gorm:"column:email"`
}
