package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
type Channel struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	DeletedStatus bool      `json:"deleted_status" gorm:"default:false"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type UserChannelDetails struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	UserName      string    `json:"user_name" gorm:"forien_key"`
	ChannelName   string    `json:"channel_name" gorm:"forien_key"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedStatus bool      `json:"deleted_status" gorm:"default:false"`
}
