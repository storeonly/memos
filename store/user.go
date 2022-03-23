package store

import (
	"fmt"
	"memos/api"
	"memos/common"
	"strings"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(create *api.UserCreate) (*api.User, error) {
	user, err := createUser(s.db, create)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) PatchUser(patch *api.UserPatch) (*api.User, error) {
	user, err := patchUser(s.db, patch)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) FindUser(find *api.UserFind) (*api.User, error) {
	list, err := findUserList(s.db, find)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	} else if len(list) > 1 {
		return nil, &common.Error{Code: common.Conflict, Err: fmt.Errorf("found %d users with filter %+v, expect 1. ", len(list), find)}
	}

	return list[0], nil
}

func createUser(db *DB, create *api.UserCreate) (*api.User, error) {
	row, err := db.Db.Query(`
		INSERT INTO user (
			name,
			password_hash,
			open_id
		)
		VALUES (?, ?, ?)
		RETURNING id, name, password_hash, open_id, created_ts, updated_ts
	`,
		create.Name,
		create.PasswordHash,
		create.OpenId,
	)
	if err != nil {
		return nil, FormatError(err)
	}
	defer row.Close()

	row.Next()
	var user api.User
	if err := row.Scan(
		&user.Id,
		&user.Name,
		&user.PasswordHash,
		&user.OpenId,
		&user.CreatedTs,
		&user.UpdatedTs,
	); err != nil {
		return nil, FormatError(err)
	}

	return &user, nil
}

func patchUser(db *DB, patch *api.UserPatch) (*api.User, error) {
	set, args := []string{}, []interface{}{}

	if v := patch.Name; v != nil {
		set, args = append(set, "name = ?"), append(args, v)
	}
	if v := patch.PasswordHash; v != nil {
		set, args = append(set, "password_hash = ?"), append(args, v)
	}
	if v := patch.OpenId; v != nil {
		set, args = append(set, "open_id = ?"), append(args, v)
	}

	args = append(args, patch.Id)

	row, err := db.Db.Query(`
		UPDATE user
		SET `+strings.Join(set, ", ")+`
		WHERE id = ?
		RETURNING id, name, password_hash, open_id, created_ts, updated_ts
	`, args...)
	if err != nil {
		return nil, FormatError(err)
	}
	defer row.Close()

	if row.Next() {
		var user api.User
		if err := row.Scan(
			&user.Id,
			&user.Name,
			&user.PasswordHash,
			&user.OpenId,
			&user.CreatedTs,
			&user.UpdatedTs,
		); err != nil {
			return nil, FormatError(err)
		}

		return &user, nil
	}

	return nil, &common.Error{Code: common.NotFound, Err: fmt.Errorf("user ID not found: %d", patch.Id)}
}

func findUserList(db *DB, find *api.UserFind) ([]*api.User, error) {
	where, args := []string{"1 = 1"}, []interface{}{}

	if v := find.Id; v != nil {
		where, args = append(where, "id = ?"), append(args, *v)
	}
	if v := find.Name; v != nil {
		where, args = append(where, "name = ?"), append(args, *v)
	}
	if v := find.OpenId; v != nil {
		where, args = append(where, "open_id = ?"), append(args, *v)
	}

	rows, err := db.Db.Query(`
		SELECT 
			id,
			name,
			password_hash,
			open_id,
			created_ts,
			updated_ts
		FROM user
		WHERE `+strings.Join(where, " AND "),
		args...,
	)
	if err != nil {
		return nil, FormatError(err)
	}
	defer rows.Close()

	list := make([]*api.User, 0)
	for rows.Next() {
		var user api.User
		if err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.PasswordHash,
			&user.OpenId,
			&user.CreatedTs,
			&user.UpdatedTs,
		); err != nil {
			fmt.Println(err)
			return nil, FormatError(err)
		}

		list = append(list, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, FormatError(err)
	}

	return list, nil
}
