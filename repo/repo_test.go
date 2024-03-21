package repo

import (
	// "fmt"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ramadhan1445sprint/sprint_segokuning/config"
	"github.com/ramadhan1445sprint/sprint_segokuning/db"
	"github.com/ramadhan1445sprint/sprint_segokuning/entity"
)

func TestCreatePost(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := db.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := NewPostRepo(db)

	testCases := []struct {
		name        string
		input       entity.Post
		errExpected bool
	}{
		{"Test create post", entity.Post{UserID: "a5801127-40e9-4528-9195-54b49a51091c", PostInHtml: "art", Tags: []string{"dwd", "woko"}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bankAccountRepo.CreatePost(&tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}

}
func TestCreateComment(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := db.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	bankAccountRepo := NewCommentRepo(db)

	testCases := []struct {
		name        string
		input       entity.Comment
		errExpected bool
	}{
		{"Test create bank", entity.Comment{UserID: "7fee1a6b-8aa1-485a-9a3a-73b29886fa59", Comment: "apa", PostID: "f3a7e013-00e4-48d3-9492-081dc09f0ebc"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := bankAccountRepo.CreateComment(&tc.input)

			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}

}

func TestGetPosts(t *testing.T) {
	config.LoadConfig("../.env")

	db, err := db.NewDatabase()
	if err != nil {
		t.Fatalf("failed to create a database connection: %v", err)
	}

	postRepo := NewPostRepo(db)

	testCases := []struct {
		name        string
		input       entity.PostFilter
		errExpected bool
	}{
		{"Test get post", entity.PostFilter{Limit: 2, Offset: 0, SearchTag: []string{"woko"}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := postRepo.GetPost(&tc.input)
			if tc.errExpected {
				if err == nil {
					t.Errorf("Expected error, but no error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}

}
