package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type FriendRepo interface {
	AddFriend(string, string) error
	DeleteFriend(string, string) error
	GetFriendPair(string, string) (*entity.FriendPair, error)
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
