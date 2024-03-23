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

func GetQueryOnlyFriend(param *entity.ListFriendPayload, joinState string, userOnlyStatement string) string {
	var conditions []string

	joinStatement := ""
	if param.OnlyFriend && param.UserId != "" {
		conditions = append(conditions, userOnlyStatement)
		joinStatement = joinState
	}

	if param.Search != "" {
		conditions = append(conditions, fmt.Sprintf("u.name LIKE '%%%s%%'", param.Search))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	if param.SortBy == entity.SortByCreatedAt {
		param.SortBy = "u.created_at"
	} else {
		param.SortBy = "u.friend_count"
	}

	query := fmt.Sprintf(`
        SELECT u.id, u.name, u.image_url, u.friend_count, u.created_at
        FROM users u
				%s
				%s
				ORDER BY %s %s
				LIMIT %d OFFSET %d`, joinStatement, whereClause, param.SortBy, param.OrderBy, param.Limit, param.Offset)

	return query
}

func GetQueryParam(param *entity.ListFriendPayload) string {
	var conditions []string

	if param.Search != "" {
		conditions = append(conditions, fmt.Sprintf("u.name LIKE '%%%s%%'", param.Search))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	if param.SortBy == entity.SortByCreatedAt {
		param.SortBy = "u.created_at"
	} else {
		param.SortBy = "u.friend_count"
	}

	query := fmt.Sprintf(`
        SELECT u.id, u.name, u.image_url, u.friend_count, u.created_at
        FROM users u
				%s
				ORDER BY %s %s
				LIMIT %d OFFSET %d`, whereClause, param.SortBy, param.OrderBy, param.Limit, param.Offset)

	return query
}

func (r *friendRepo) GetListFriends(param *entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error) {
	var listUser []entity.UserList
	var query string

	if param.OnlyFriend {
		query = GetQueryOnlyFriend(param, "JOIN friends f ON u.id = f.user_id1", fmt.Sprintf("f.user_id2 = '%s'", param.UserId))

		err := r.db.Select(&listUser, query)
		if err != nil {
			log.Println("Error executing query:", err)
			return listUser, nil, err
		}

		var listUser1 []entity.UserList
		query2 := GetQueryOnlyFriend(param, "JOIN friends f ON u.id = f.user_id2", fmt.Sprintf("f.user_id1 = '%s'", param.UserId))
		err1 := r.db.Select(&listUser1, query2)
		if err1 != nil {
			log.Println("Error executing query:", err1)
			return listUser1, nil, err1
		}

		listUser = append(listUser, listUser1...)
	} else {
		query = GetQueryParam(param)

		err := r.db.Select(&listUser, query)
		if err != nil {
			log.Println("Error executing query:", err)
			return listUser, nil, err
		}
	}

	pagination := &entity.Meta{
		Limit:  param.Limit,
		Offset: param.Offset,
		Total:  len(listUser),
	}

	return listUser, pagination, nil
}
