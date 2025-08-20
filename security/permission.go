package security

type Permission struct {
	ID   string `gorm:"id,primaryKey,size:255"`
	Name string `gorm:"name,size:255"`

	_ []*Role `gorm:"many2many:role_permissions"`
}

func (Permission) TableName() string {
	return "permissions"
}
