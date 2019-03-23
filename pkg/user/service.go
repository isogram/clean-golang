package user

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/isogram/clean-golang/pkg/entity"
	"github.com/isogram/clean-golang/pkg/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//Service service interface
type Service struct {
	repo Repository
}

//NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

// Me show detail user
func (s *Service) Me(ID int64) (*entity.User, error) {
	user, err := s.repo.GetByID(ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Register an user to system
func (s *Service) Register(ur *entity.UserRegister) (*entity.User, error) {
	var u entity.User

	u.UserName = ur.UserName
	u.FullName = ur.FullName
	u.Email = ur.Email

	password := []byte(ur.Password)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u.Password = string(hashedPassword)

	user, err := s.repo.Store(&u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Login for user
func (s *Service) Login(email, password string) (u *entity.User, tok string, refreshToken string, e error) {
	var err error

	user, err := s.repo.GetByEmail(email)
	if err != nil {

		if err == sql.ErrNoRows {
			err = errors.New("Your email and password combination does not match")
			return nil, "", "", err
		}

		return nil, "", "", err
	}

	// check matched password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		err = errors.New("Your email and password combination does not match")
		return nil, "", "", err
	}

	// check status
	if user.Status != "active" {
		err = errors.New("Your account is not activated or banned. Please contact administrator")
		return nil, "", "", err
	}

	// creating JWT payload
	encodedUserID := utils.Int64ToHashed(user.ID)

	durationStr := os.Getenv("JWT_EXPIRE_DURATION")
	durationInt, _ := strconv.Atoi(durationStr)
	duration := time.Duration(int64(time.Hour) * int64(durationInt))
	expiresAt := time.Now().Add(duration).Unix()
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &entity.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		encodedUserID,
	}

	log.Println(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		err = errors.New("Problem on generating token. Please try again")
		return nil, "", "", err
	}

	refreshToken, err = s.repo.CreateRefreshToken(user.ID)
	if err != nil {
		err = errors.New("Problem on generating refresh token. Please try again")
		return nil, "", "", err
	}

	return user, tokenString, refreshToken, nil

}

// RefreshToken get new token based on given refresh token
func (s *Service) RefreshToken(UID int64, r string) (u *entity.User, tok string, refreshToken string, e error) {
	var err error

	rToken, err := s.repo.GetRefreshToken(r)
	if err != nil {

		if err == sql.ErrNoRows {
			err = errors.New("No refresh token matched")
		}

		return nil, "", "", err
	}

	// check UID matched with refresh token
	if UID != rToken.UserID {
		err = errors.New("Refresh token not belong this user")
		return nil, "", "", err
	}

	// check refresh token is already revoked
	if rToken.Revoked != 0 {
		err = errors.New("Refresh token has been revoked")
		return nil, "", "", err
	}

	user, err := s.repo.GetByID(UID)
	if err != nil {

		if err == sql.ErrNoRows {
			err = errors.New("No such user in database")
			return nil, "", "", err
		}

		return nil, "", "", err
	}

	// creating JWT payload
	encodedUserID := utils.Int64ToHashed(user.ID)

	durationStr := os.Getenv("JWT_EXPIRE_DURATION")
	durationInt, _ := strconv.Atoi(durationStr)
	duration := time.Duration(int64(time.Hour) * int64(durationInt))
	expiresAt := time.Now().Add(duration).Unix()
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &entity.AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		encodedUserID,
	}

	log.Println(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		err = errors.New("Problem on generating token. Please try again")
		return nil, "", "", err
	}

	refreshToken, err = s.repo.CreateRefreshToken(user.ID)
	if err != nil {
		err = errors.New("Problem on generating refresh token. Please try again")
		return nil, "", "", err
	}

	err = s.repo.RevokeRefreshToken(r)
	if err != nil {
		err = errors.New("Problem on revoked refresh token. Please try again")
		return nil, "", "", err
	}

	return user, tokenString, refreshToken, nil
}
