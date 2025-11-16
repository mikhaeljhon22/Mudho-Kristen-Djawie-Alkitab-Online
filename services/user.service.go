package services

import (
	"alkitab/entitys"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db    *gorm.DB
	mutex sync.Mutex
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *UserService) FindUserID(username string) int {
	var user entitys.UsersLetstalk
	s.db.Where("username = ?", username).First(&user)
	fmt.Print(user.ID)
	return user.ID
}

func (s *UserService) SignUpAddUser(user entitys.UsersLetstalk) bool {
	var existingUser entitys.UsersLetstalk
	findUsername := s.db.Where("username = ?", user.Username).Or("email = ?", user.Email).Find(&existingUser)
	if findUsername.RowsAffected == 1 {
		return false
	} else {
		s.mutex.Lock()
		s.db.Create(&user)
		s.mutex.Unlock()
		return true
	}
}

func (s *UserService) GetUserByUsername(username string) (*entitys.UsersLetstalk, error) {
	var user entitys.UsersLetstalk
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) StoreCodeVerif(username string, code string) {
	var user entitys.UsersLetstalk
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.db.Model(&user).Where("username = ?", username).Update("code", code)
}
func (s *UserService) SigninUser(user entitys.UsersLetstalk) bool {
	var foundUsers entitys.UsersLetstalk
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.db.Where("username = ? AND is_active = ?", user.Username, true).First(&foundUsers)
	checKpW := CheckPasswordHash(user.Password, foundUsers.Password)
	return checKpW
}

func (s *UserService) ProfileUser(username string) entitys.UsersLetstalk {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var foundUser entitys.UsersLetstalk
	s.db.Where("username = ?", username).First(&foundUser)
	return foundUser
}

func (s *UserService) VerifyCode(code string) bool {
	var user entitys.UsersLetstalk
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.db.Model(&user).Where("code = ?", code).Update("is_active", true)
	return true
}
