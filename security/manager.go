package security

import (
	"errors"
	"fmt"
	"regexp"
	"slices"

	"github.com/ramoncl001/comet/data"
)

type OperationResult struct {
	Success bool
	Errors  []error
}

func success() OperationResult {
	return OperationResult{
		Success: true,
		Errors:  []error{},
	}
}

func failure(errs ...string) OperationResult {
	var v []error
	for _, e := range errs {
		v = append(v, errors.New(string(e)))
	}
	return OperationResult{
		Success: false,
		Errors:  v,
	}
}

type UserManager interface {
	Create(user *User) OperationResult
	GetByEmail(email string) *User
	GetByID(id string) *User
	SetPassword(user *User, password string) OperationResult
	CheckPassword(user *User, password string) bool
	AddRole(user *User, roles ...*Role) OperationResult
	RetrieveRole(user *User, roles ...*Role) OperationResult
}

type defaultUserManager struct {
	UserManager
	config *UserConfig
	db     *data.DatabaseContext
}

func NewDefaultUserManager(config *UserConfig, db *data.DatabaseContext) UserManager {
	return &defaultUserManager{
		config: config,
		db:     db,
	}
}

func (mg *defaultUserManager) Create(user *User) OperationResult {
	err := mg.db.Create(&user).Error
	if err != nil {
		return failure(err.Error())
	}

	return success()
}

func (mg *defaultUserManager) GetByEmail(email string) *User {
	var user User
	err := mg.db.Preload("Roles").First(&user, "email = ?", email).Error
	if err != nil {
		return nil
	}

	return &user
}

func (mg *defaultUserManager) GetByID(id string) *User {
	var user User
	err := mg.db.Preload("Roles").First(&user, "id = ?", id).Error
	if err != nil {
		return nil
	}

	return &user
}

func (mg *defaultUserManager) SetPassword(user *User, password string) OperationResult {
	errs := mg.validatePassword(password)
	if len(errs) != 0 {
		return OperationResult{
			Success: false,
			Errors:  errs,
		}
	}

	passwordHash := SHA256(password)
	user.PasswordHash = passwordHash
	if user.ID != "" {
		mg.db.Save(&user)
	}

	return success()
}

func (mg *defaultUserManager) CheckPassword(user *User, password string) bool {
	return SHA256(password) == user.PasswordHash
}

func (mg *defaultUserManager) AddRole(user *User, roles ...*Role) OperationResult {
	if user.ID == "" {
		return failure("user does not exists")
	}

	user.Roles = append(user.Roles, roles...)

	if err := mg.db.Save(&user).Error; err != nil {
		return failure(err.Error())
	}

	return success()
}

func (mg *defaultUserManager) RetrieveRole(user *User, roles ...*Role) OperationResult {
	if user.ID == "" {
		return failure("user does not exists")
	}

	user.Roles = slices.DeleteFunc(user.Roles, func(r *Role) bool {
		return slices.Contains(roles, r)
	})

	if err := mg.db.Save(&user).Error; err != nil {
		return failure(err.Error())
	}

	return success()
}

func (mg *defaultUserManager) validatePassword(password string) []error {
	var errs []error
	if len(password) < mg.config.PasswordConfig.MinimumChars {
		errs = append(errs, fmt.Errorf("password must have at least %d characters", mg.config.PasswordConfig.MinimumChars))
	}

	passwordBytes := []byte(password)

	if mg.config.PasswordConfig.NeedDigits {
		rgx := regexp.MustCompile(digitRegex)
		if !rgx.Match(passwordBytes) {
			errs = append(errs, errors.New("password must have at least one digit"))
		}
	}

	if mg.config.PasswordConfig.NeedLowercase {
		rgx := regexp.MustCompile(lowerRegex)
		if !rgx.Match(passwordBytes) {
			errs = append(errs, errors.New("password must have at least a lower case letter"))
		}
	}

	if mg.config.PasswordConfig.NeedUppercase {
		rgx := regexp.MustCompile(upperRegex)
		if !rgx.Match(passwordBytes) {
			errs = append(errs, errors.New("password must have at least an upper case letter"))
		}
	}

	if mg.config.PasswordConfig.NeedSpecialChars {
		rgx := regexp.MustCompile(specialRegex)
		if !rgx.Match(passwordBytes) {
			errs = append(errs, errors.New("password must have at least a special character"))
		}
	}

	return errs
}
