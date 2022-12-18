package hello

type HelloRequest struct {
	ID    string `json:"id" param:"path,id"`
	Query Query  `json:"query" param:",query"`
}

type HelloResponse struct {
	ID string `json:"id"`
}

type Query struct {
	Age         string   `json:"age" param:"query,age" query:"eq" gorm:"column:age"`
	Email       string   `json:"email" param:"query,email" query:"ne" gorm:"column:email"`
	IDIn        []string `json:"idIn" param:"query,idIn" query:"in" gorm:"column:id"`
	BetweenTime []string `json:"between_time" param:"query,between_time" query:"between" gorm:"column:time"`
}

type Pm struct {
	UUID  string `json:"uuid" param:"path,uuid" query:"eq" gorm:"column:age"`
	Brand struct {
	} `json:"brand" gorm:"table:brand;foreignKey:BrandUUID" query:"in"`
}
