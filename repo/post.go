package repo

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

type PostRepo interface {
	CreatePost(post *entity.Post) error
	GetPost(filter *entity.PostFilter) ([]entity.PostData, error)
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
	query += fmt.Sprintf(" ORDER BY p.created_at DESC limit %d offset %d", filter.Limit, filter.Offset)

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

		var tagsSlice []string
		for _, tag := range postRawDBData.Tags.Elements {
			if tag.Status != pgtype.Null {
				tagsSlice = append(tagsSlice, string(tag.String))
			}
		}

		tempPostDetail.Tags = tagsSlice
		tempPostDetail.CreatedAt = postRawDBData.PostCreatedAt.Format(time.RFC3339)
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
		if postMapComment[key] == nil {
			temp.Comments = []entity.PostComments{}
		} else {
			temp.Comments = postMapComment[key]
		}
		postData = append(postData, temp)
	}

	slices.SortFunc(postData, func(a, b entity.PostData) int {
		t1, _ := time.Parse(time.RFC3339, a.Post.CreatedAt)
		t2, _ := time.Parse(time.RFC3339, b.Post.CreatedAt)

		if t1.After(t2) {
			return -1
		} else {
			return 1
		}
	})

	return postData, nil
}
