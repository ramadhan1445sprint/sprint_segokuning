package repo

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type FriendRepo interface {
	AddFriend(string, string) error
	DeleteFriend(string, string) error
	GetFriendPair(string, string) (*entity.FriendPair, error)
	GetListFriends(*entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error)
}

type friendRepo struct {
	db *sqlx.DB
}

func NewFriendRepo(db *sqlx.DB) FriendRepo {
	return &friendRepo{db: db}
}

func (r *friendRepo) GetFriendPair(userId, friendId string) (*entity.FriendPair, error) {
	var friendPair entity.FriendPair

	statement := `
		SELECT id, user_id1, user_id2 FROM friends
		WHERE
			(user_id1 = $1 AND user_id2 = $2)
		OR
			(user_id1 = $2 AND user_id2 = $1)
	`

	err := r.db.Get(&friendPair, statement, userId, friendId)
	if err != nil {
		return nil, err
	}

	return &friendPair, nil
}

func (r *friendRepo) AddFriend(userId, friendId string) (err error) {
	tx := r.db.MustBegin()

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()

	// add friend
	_, err = tx.Exec("INSERT INTO friends (user_id1, user_id2) VALUES ($1, $2)", userId, friendId)
	if err != nil {
		return err
	}

	// increment friendcount
	_, err = tx.Exec("UPDATE users SET friend_count = friend_count + 1 WHERE id IN ($1, $2)",
		userId,
		friendId,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *friendRepo) DeleteFriend(userId, friendId string) (err error) {
	tx := r.db.MustBegin()

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()

			if rollbackErr != nil {
				err = rollbackErr
			}
		}
	}()

	// delete friend
	_, err = tx.Exec(`
		DELETE FROM friends 
		WHERE 
			(user_id1 = $1 AND user_id2 = $2)
		OR
			(user_id1 = $2 AND user_id2 = $1)
		`, userId, friendId)
	if err != nil {
		return err
	}

	// decrement friendcount
	_, err = tx.Exec("UPDATE users SET friend_count = friend_count - 1 WHERE id IN ($1, $2)",
		userId,
		friendId,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (r *friendRepo) GetListFriends(param *entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error) {
	var conditions []string

	joinStatement := ""
	if param.OnlyFriend && param.UserId != "" {
		conditions = append(conditions, fmt.Sprintf("users.id = '%s'", param.UserId))
		joinStatement = fmt.Sprintf("JOIN friends ON users.id = friends.user_id1 OR users.id = friends.user_id2")
	}

	if param.Search != "" {
		conditions = append(conditions, fmt.Sprintf("users.name LIKE '%%%s%%'", param.Search))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	if param.SortBy == entity.SortByCreatedAt {
		param.SortBy = "users.created_at"
	} else {
		param.SortBy = "users.friend_count"
	}

	query := fmt.Sprintf(`
        SELECT users.id, users.name, users.image_url, users.friend_count, users.created_at
        FROM users
				%s
				%s
				ORDER BY %s %s
				LIMIT %d OFFSET %d`, joinStatement, whereClause, param.SortBy, param.OrderBy, param.Limit, param.Offset)

	var listUser []entity.UserList
	err := r.db.Select(&listUser, query)
	if err != nil {
		log.Println("Error executing query:", err)
		return listUser, nil, err
	}

	pagination := &entity.Meta{
		Limit:  param.Limit,
		Offset: param.Offset,
		Total:  len(listUser),
	}

	return listUser, pagination, nil
}
