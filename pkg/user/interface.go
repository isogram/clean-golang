package user

import "github.com/isogram/clean-golang/pkg/entity"

//Reader interface
type Reader interface {
	Me(ID int64) (*entity.User, error)
	Login(email, password string) (u *entity.User, token string, refreshToken string, err error)
	RefreshToken(UID int64, refreshToken string) (u *entity.User, token string, rt string, err error)
}

//Writer user writer
type Writer interface {
	Register(u *entity.UserRegister) (*entity.User, error)
}

//Repository repository interface
type Repository interface {
	GetByID(ID int64) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Store(u *entity.User) (*entity.User, error)
	CreateRefreshToken(UID int64) (refreshToken string, er error)
	GetRefreshToken(refreshToken string) (*entity.AuthRefreshToken, error)
	RevokeRefreshToken(refreshToken string) error
}

//UseCase use case interface
type UseCase interface {
	Reader
	Writer
}
