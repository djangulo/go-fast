package main

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

// Player struct to hold base data
type Player struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string `gorm:"type:varchar(50);unique;not null;" json:"name"`
	Wins      int    `gorm:"not null;default:0;" json:"wins"`
}

func (p *Player) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid)
}
