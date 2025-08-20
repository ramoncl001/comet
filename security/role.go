package security

type Role struct {
	ID   string `gorm:"id,primaryKey,size:255"`
	Name string `gorm:"name,not null,size:255"`

	_           []User        `gorm:"many2many:user_roles"`
	Permissions []*Permission `gorm:"many2many:role_permissions"`
}

func (Role) TableName() string {
	return "roles"
}
