package rbac

import (
	"context"
	"errors"
	"time"

	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/database"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/utils"
)

type UserHelper[T any, IDTYPE int64 | string] interface {
	database.CrudHelper[T, User[IDTYPE], IDTYPE]
	Login(ctx context.Context, username string, password string, condition database.Condition[T]) (string, error)
}

type userHelper[T any, IDTYPE int64 | string] struct {
	helper database.CrudHelper[T, User[IDTYPE], IDTYPE]
	secret string
}

func NewUserHelper[T any, IDTYPE int64 | string](helper database.CrudHelper[T, User[IDTYPE], IDTYPE], secret string) UserHelper[T, IDTYPE] {
	return &userHelper[T, IDTYPE]{
		helper: helper,
		secret: secret,
	}
}

func (u *userHelper[T, IDTYPE]) Create(ctx context.Context, user *User[IDTYPE]) (*User[IDTYPE], error) {
	if user.Username == "" {
		return nil, errors.New("username is required")
	}

	if user.Password == "" {
		return nil, errors.New("password is required")
	}

	pass, err := encryptPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = pass

	user, err = u.helper.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.Password = "" // removing password as it should not be available in reponse
	return user, nil
}

func (u *userHelper[T, IDTYPE]) Update(ctx context.Context, user *User[IDTYPE], project []string, condition database.Condition[T]) error {
	if len(project) == 0 || project[0] == "" || project[0] == "*" {
		return errors.New("all field updates are not allowed")
	}

	projectSet := utils.NewSetFromSlice(project)

	if user.Password == "" && projectSet.Contains("password") {
		return errors.New("password is required")
	}

	if user.Username == "" && projectSet.Contains("username") {
		return errors.New("username is required")
	}

	if user.Password != "" {
		pass, err := encryptPassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = pass
	}

	return u.helper.Update(ctx, user, project, condition)
}

func (u *userHelper[T, IDTYPE]) Get(ctx context.Context, project []string, condition database.Condition[T]) ([]User[IDTYPE], error) {
	if len(project) == 0 || project[0] == "" || project[0] == "*" {
		return nil, errors.New("all field get are not allowed")
	}

	projectSet := utils.NewSetFromSlice(project)
	if projectSet.Contains("password") {
		return nil, errors.New("password field is not allowed")
	}

	return u.helper.Get(ctx, project, condition)
}

func (u *userHelper[T, IDTYPE]) GetTableName() string {
	return u.helper.GetTableName()
}

func (u *userHelper[T, IDTYPE]) GetColumns(project []string, withoutID bool) []string {
	return u.helper.GetColumns(project, withoutID)
}

func (u *userHelper[T, IDTYPE]) Delete(ctx context.Context, condition database.Condition[T]) error {
	return u.helper.Delete(ctx, condition)
}

func (u *userHelper[T, IDTYPE]) Login(ctx context.Context, username string, password string, condition database.Condition[T]) (string, error) {
	users, err := u.helper.Get(ctx, []string{"id", "username", "password"}, condition.New().Set("username", database.ConditionOperationEqual, username))
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "", errors.New("username not found")
	}

	user := users[0]

	if !comparePassword(password, user.Password) {
		return "", errors.New("invalid password")
	}

	user.Password = "" // removing this as it should not be available in token

	return GenerateJWT(user, time.Hour, u.secret)
}
