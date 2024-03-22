package repo

import (
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type FriendRepo interface {
	GetListFriends(*entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error)
}

type friendRepo struct {
	db *sqlx.DB
}

func NewFriendRepo(db *sqlx.DB) FriendRepo {
	return &friendRepo{db}
}

func (r *friendRepo) GetListFriends(param *entity.ListFriendPayload) ([]entity.UserList, *entity.Meta, error) {
	var conditions []string

	fmt.Print(param.UserId)

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

	fmt.Printf(query)

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
