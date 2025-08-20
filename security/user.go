package security

import "time"

type ApplicationUser interface{}

type User struct {
	ApplicationUser
	ID           string  `gorm:"id,primaryKey,size:255"`
	Username     string  `gorm:"username,not null,unique,size:255"`
	Email        string  `gorm:"email,not null,size:255"`
	PhoneNumber  *string `gorm:"phone_number,size:255"`
	PasswordHash string  `gorm:"password_hash,not null,size:256"`
	IsActive     bool    `gorm:"is_active,default:false"`

	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt int64     `gorm:"updated_at,autoUpdate:milli"`

	Roles []*Role `gorm:"many2many:user_roles"`
}

func (User) TableName() string {
	return "users"
}
