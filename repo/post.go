package repo

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type PostRepo interface {
	CreatePost(post *entity.Post) error
	GetPost(filter *entity.PostFilter) ([]entity.PostData, error)
	GetTotalPost() (int, error)
	CheckPostById(postId string) (bool, error)
	CheckFriendPost(postId string, userId string) (bool, error)
}

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepo {
	return &postRepo{db}
}

func (r *postRepo) CreatePost(post *entity.Post) error {
	_, err := r.db.Exec("INSERT INTO posts (post_in_html, tags, user_id) VALUES ($1, $2, $3)", post.PostInHtml, post.Tags, post.UserID)

	if err != nil {
		return err
	}

	return nil
}

func (r *postRepo) GetTotalPost() (int, error) {
	var total int

	err := r.db.Get(total, "SELECT count(*) from posts")

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *postRepo) CheckPostById(postId string) (bool, error) {
	var exist int

	err := r.db.Get(exist, "SELECT count(*) from posts where id = $1", postId)

	if err != nil {
		return false, err
	}

	if exist > 0 {
		return true, nil
	}

	return false, nil

}
func (r *postRepo) CheckFriendPost(postId string, userId string) (bool, error) {
	var creator string
	var exist int

	err := r.db.Get(creator, "SELECT user_id from posts where id = $1", postId)
	if err != nil {
		return false, err
	}

	err = r.db.Get(creator, "SELECT count(*) from posts where (userid1 = $1 and userid2 = $2) or (userid1 = $3 and userid2 = $4)", userId, creator, creator, userId)
	if err != nil {
		return false, err
	}

	if exist > 0 {
		return true, nil
	}

	return false, nil
}

func (r *postRepo) GetPost(filter *entity.PostFilter) ([]entity.PostData, error) {
	where := ""
	query := `
	SELECT
		p.id AS "post_id",
		p.post_in_html AS "post_in_html",
		p.tags AS "posts_tags",
		p.created_at AS "posts_created_at",
		u.id AS "creator_id",
		u.name AS "creator_name",
		u.image_url AS "creator_image_url",
		u.friend_count AS "creator_friend_count",
		u.created_at AS "creator_created_at",
		c.comment AS comment,
		c.created_at AS "comment_created_at",
		cu.id AS "comment_creator_id",
		cu.name AS "comment_creator_name",
		cu.image_url AS "comment_creator_image_url",
		cu.friend_count AS "comment_creator_friend_count"
	FROM
		posts p
	JOIN
		users u ON p.user_id = u.id
	left JOIN
		comments c ON p.id = c.post_id
	left JOIN
		users cu ON c.user_id = cu.id
	`

	if filter.Search != "" {
		if where != "" {
			where += " AND p.post_in_html LIKE '%" + filter.Search + "%' "
		} else {
			where += " WHERE p.post_in_html LIKE '%" + filter.Search + "%' "
		}
	}

	if len(filter.SearchTag) > 0 {
		jsonTag, err := json.Marshal([]string(filter.SearchTag))
		if err == nil {
			replacer := strings.NewReplacer("[", "{", "]", "}")
			stringTag := replacer.Replace(string(jsonTag))
			if where != "" {
				where += fmt.Sprintf(" AND tags && '%s'", stringTag)
			} else {
				where += fmt.Sprintf(" WHERE tags && '%s'", stringTag)
			}
		}
	}
	query += where
	query += fmt.Sprintf(" ORDER BY p.created_at DESC, c.created_at desc limit %d offset %d", filter.Limit, filter.Offset)
	fmt.Print(query)

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	postData := []entity.PostData{}
	postMap := make(map[string]entity.PostData)
	postMapComment := make(map[string][]entity.PostComments)

	for rows.Next() {
		var postRawDBData entity.PostRawDBData
		var tempPostComment entity.PostComments
		var tempPostDetail entity.PostDetail
		var tempPostCreator entity.PostUser
		var tempPostCommentCreator entity.PostUser
		var tempPostData entity.PostData

		if err = rows.Scan(&postRawDBData.PostID, &postRawDBData.PostInHTML, &postRawDBData.Tags, &postRawDBData.PostCreatedAt,
			&postRawDBData.CreatorID, &postRawDBData.CreatorName, &postRawDBData.CreatorImageURL, &postRawDBData.CreatorFriendCount, &postRawDBData.CreatorCreatedAt, &postRawDBData.Comment, &postRawDBData.CommentCreatedAt,
			&postRawDBData.CommentCreatorID, &postRawDBData.CommentCreatorName, &postRawDBData.CommentCreatorImageUrl, &postRawDBData.CommentCreatorFriendCount,
		); err != nil {
			return nil, err
		}

		tempPostData.PostId = postRawDBData.PostID
		tempPostDetail.PostInHtml = postRawDBData.PostInHTML
		tempPostDetail.Tags = postRawDBData.Tags
		tempPostDetail.CreatedAt = postRawDBData.PostCreatedAt
		tempPostData.Post = tempPostDetail

		tempPostCreator.UserID = postRawDBData.CreatorID
		tempPostCreator.Name = postRawDBData.CreatorName
		tempPostCreator.ImageUrl = postRawDBData.CreatorImageURL
		tempPostCreator.FriendCount = postRawDBData.CreatorFriendCount
		tempPostCreator.CreatedAt = postRawDBData.CreatorCreatedAt
		tempPostData.Creator = tempPostCreator

		if _, ok := postMap[postRawDBData.PostID]; !ok {
			tempPostData.Comments = []entity.PostComments{}
			postMap[postRawDBData.PostID] = tempPostData
		}

		if postRawDBData.CommentCreatorID != nil {
			tempPostCommentCreator.UserID = *postRawDBData.CommentCreatorID
		} else {
			continue
		}

		if postRawDBData.CommentCreatorName != nil {
			tempPostCommentCreator.Name = *postRawDBData.CommentCreatorName
		}

		if postRawDBData.CommentCreatorImageUrl != nil {
			tempPostCommentCreator.ImageUrl = *postRawDBData.CommentCreatorImageUrl
		}
		if postRawDBData.CommentCreatorFriendCount != nil {
			tempPostCommentCreator.FriendCount = *postRawDBData.CommentCreatorFriendCount
		}

		if postRawDBData.Comment != nil {
			tempPostComment.Comment = *postRawDBData.Comment
		}

		tempPostComment.Creator = tempPostCommentCreator

		if postRawDBData.CommentCreatedAt != nil {
			tempPostComment.CreatedAt = *postRawDBData.CommentCreatedAt
		}

		if _, ok := postMapComment[postRawDBData.PostID]; !ok {
			postMapComment[postRawDBData.PostID] = []entity.PostComments{tempPostComment}
		} else {
			postMapComment[postRawDBData.PostID] = append(postMapComment[postRawDBData.PostID], tempPostComment)
		}
	}

	for key, _ := range postMap {
		temp := postMap[key]
		temp.Comments = postMapComment[key]
		postData = append(postData, temp)
	}

	return postData, nil
}
