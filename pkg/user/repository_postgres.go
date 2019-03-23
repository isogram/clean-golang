package user

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/isogram/clean-golang/pkg/entity"
	"github.com/isogram/clean-golang/pkg/utils"
)

const (
	DefaultUserStatus = "inactive"
)

// UserRepoPostgres data structure
type UserRepoPostgres struct {
	readDB  *sql.DB
	writeDB *sql.DB
}

// NewRepoPostgres function for initializing user repo
func NewRepoPostgres(readDB, writeDB *sql.DB) *UserRepoPostgres {
	return &UserRepoPostgres{readDB, writeDB}
}

// GetByID get user detail
func (r *UserRepoPostgres) GetByID(ID int64) (*entity.User, error) {
	var (
		user entity.User
		err  error
	)
	err = r.readDB.QueryRow("select id, fullname, username, email, password, status, created_at from users where id = $1", ID).Scan(&user.ID, &user.FullName, &user.UserName, &user.Email, &user.Password, &user.Status, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail get user detail by email
func (r *UserRepoPostgres) GetByEmail(email string) (*entity.User, error) {
	var (
		user entity.User
		err  error
	)
	err = r.readDB.QueryRow("select id, fullname, username, email, password, status, created_at from users where email = $1", email).Scan(&user.ID, &user.FullName, &user.UserName, &user.Email, &user.Password, &user.Status, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Store a user
func (r *UserRepoPostgres) Store(u *entity.User) (*entity.User, error) {

	sql := fmt.Sprintf(`INSERT INTO users(fullname, username, password, email, status) VALUES('%s', '%s', '%s', '%s', '%s') returning id, created_at`, u.FullName, u.UserName, u.Password, u.Email, DefaultUserStatus)

	var (
		lastID    int64
		createdAt time.Time
	)
	err := r.writeDB.QueryRow(sql).Scan(&lastID, &createdAt)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	u.ID = lastID
	u.Status = DefaultUserStatus
	u.CreatedAt = createdAt

	return u, nil
}

// CreateRefreshToken generating refresh token for a user
func (r *UserRepoPostgres) CreateRefreshToken(UID int64) (rt string, e error) {

	strTime := fmt.Sprintf("%v", time.Now().UnixNano())
	refreshToken := utils.StringMD5(strTime)

	sql := fmt.Sprintf(`INSERT INTO refresh_token(user_id, refresh_token) VALUES('%v', '%s') returning id`, UID, refreshToken)

	var (
		lastID int64
	)
	err := r.writeDB.QueryRow(sql).Scan(&lastID)
	if err != nil {
		return "", err
	}

	// REVOKED ALL OLD REFRESH TOKEN
	// stmt, err := r.writeDB.Prepare("UPDATE refresh_token SET revoked=1 WHERE user_id=$1 AND id < $2")
	// if err != nil {
	// 	return "", err
	// }
	// _, err = stmt.Exec(UID, lastID)
	// if err != nil {
	// 	return "", err
	// }

	return refreshToken, nil
}

// GetRefreshToken get refresh token by specify "refresh token"
func (r *UserRepoPostgres) GetRefreshToken(refreshToken string) (*entity.AuthRefreshToken, error) {
	var (
		rt  entity.AuthRefreshToken
		err error
	)
	err = r.readDB.QueryRow("select id, user_id, refresh_token, revoked from refresh_token where refresh_token = $1", refreshToken).
		Scan(&rt.ID, &rt.UserID, &rt.RefreshToken, &rt.Revoked)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

// RevokeRefreshToken set refresh token revoked = 1
func (r *UserRepoPostgres) RevokeRefreshToken(refreshToken string) error {
	stmt, err := r.writeDB.Prepare("UPDATE refresh_token SET revoked=1 WHERE refresh_token=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(refreshToken)
	if err != nil {
		return err
	}

	return nil
}
