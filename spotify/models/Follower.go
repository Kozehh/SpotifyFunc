package models

type Follower struct {
	Name  string `json:"name"`
	Album string `json:"album"`
}

func (b *Follower) TableName() string {
	return "follower"
}
