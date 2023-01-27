package base

type Model struct {
	Id        int    `json:"id" bson:"id" gorm:"primarykey"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	UpdatedAt string `json:"UpdatedAt" bson:"UpdatedAt"`
}
