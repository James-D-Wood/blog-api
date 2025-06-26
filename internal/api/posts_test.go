package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/James-D-Wood/blog-api/internal/constant"
	"github.com/James-D-Wood/blog-api/internal/db"
	"github.com/James-D-Wood/blog-api/internal/model"
)

var TestUserMap = map[string]*model.User{
	"kishiguro": {
		ID:       "0197aaed-4a35-74da-8574-4165524a1111",
		Username: "kishiguro",
		Name:     "Kazuo Ishiguro",
		IsAdmin:  false,
	},
	"dsedaris": {
		ID:       "0197aaed-4a35-74da-8574-4165524a2222",
		Username: "dsedaris",
		Name:     "David Sedaris",
		IsAdmin:  false,
	},
	"admin": {
		ID:       "0197aaed-4a35-74da-8574-4165524a3333",
		Username: "admin",
		Name:     "James Wood",
		IsAdmin:  true,
	},
}

// set up basic test cases
var createBlogPostTestCases = []struct {
	Name         string
	RequestBody  map[string]string
	User         string
	ResponseCode int
}{
	{
		Name: "Happy Path",
		RequestBody: map[string]string{
			"title":    "My NEW riveting blog post",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		User:         "0197aaed-4a35-74da-8574-4165524a1111",
		ResponseCode: 201,
	},
	{
		Name: "Unknown User",
		RequestBody: map[string]string{
			"title":    "My NEW riveting blog post",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		User:         "",
		ResponseCode: 500,
	},
	{
		Name: "Duplicate Post",
		RequestBody: map[string]string{
			"title":    "duplicate title",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		User:         "0197aaed-4a35-74da-8574-4165524a1111",
		ResponseCode: 400,
	},
}

func TestCreateBlogPostHandler(t *testing.T) {
	for _, tt := range createBlogPostTestCases {
		t.Run(tt.Name, func(t *testing.T) {

			// InMemory implementations double as mock test implementations for unit tests
			app := App{
				Logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				BlogService: db.NewInMemoryBlogService(),
				UserService: &db.InMemoryUserService{
					Users: TestUserMap,
				},
			}

			// seed an existing post beforehand for duplicate scenario
			app.BlogService.CreateBlogPost(context.TODO(), "0197aaed-4a35-74da-8574-4165524a1111", &model.BlogPost{
				Title: "duplicate title",
			})

			b, _ := json.Marshal(tt.RequestBody)

			req := httptest.NewRequest("POST", "/api/v1/posts", bytes.NewReader(b))
			ctx := context.WithValue(req.Context(), constant.UserIDKey, tt.User)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			app.CreateBlogPostHandler(rr, req)
			if rr.Result().StatusCode != tt.ResponseCode {
				t.Errorf("got %d, want %d", rr.Result().StatusCode, tt.ResponseCode)
			}
		})
	}
}

var fetchBlogPostTestCases = []struct {
	Name         string
	PostID       string
	User         string
	IsDraft      bool
	ResponseCode int
}{
	{
		Name:         "Published Post - Same Owner",
		User:         "0197aaed-4a35-74da-8574-4165524a1111",
		IsDraft:      false,
		ResponseCode: 200,
	},
	{
		Name:         "Published Post - Different Owner",
		User:         "0197aaed-4a35-74da-8574-4165524a2222",
		IsDraft:      false,
		ResponseCode: 200,
	},
	{
		Name:         "Draft Post - Same Owner",
		User:         "0197aaed-4a35-74da-8574-4165524a1111",
		IsDraft:      true,
		ResponseCode: 200,
	},
	{
		Name:         "Draft Post - Different Owner",
		User:         "0197aaed-4a35-74da-8574-4165524a2222",
		IsDraft:      true,
		ResponseCode: 403,
	},
	{
		Name:         "Post Does Not Exist",
		PostID:       "efbfa286-ca55-4ded-a28e-9881118186c8",
		User:         "0197aaed-4a35-74da-8574-4165524a2222",
		IsDraft:      true,
		ResponseCode: 404,
	},
}

func TestFetchBlogPostHandler(t *testing.T) {
	for _, tt := range fetchBlogPostTestCases {
		t.Run(tt.Name, func(t *testing.T) {

			// InMemory implementations double as mock test implementations for unit tests
			app := App{
				Logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				BlogService: db.NewInMemoryBlogService(),
				UserService: &db.InMemoryUserService{
					Users: TestUserMap,
				},
			}

			// set status
			var status model.BlogPostStatus
			if tt.IsDraft {
				status = model.DRAFT
			} else {
				status = model.PUBLISHED
			}

			// seed an existing post beforehand
			blog := &model.BlogPost{
				Status: status,
			}
			err := app.BlogService.CreateBlogPost(context.TODO(), "0197aaed-4a35-74da-8574-4165524a1111", blog)
			if err != nil {
				t.Error(err)
			}

			postID := blog.ID
			if tt.PostID != "" {
				// override post ID for request
				postID = tt.PostID
			}

			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/posts/%s", postID), nil)
			req.SetPathValue("id", postID)

			// set user identity
			ctx := context.WithValue(req.Context(), constant.UserIDKey, tt.User)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			app.FetchBlogPostHandler(rr, req)
			if rr.Result().StatusCode != tt.ResponseCode {
				t.Errorf("got %d, want %d", rr.Result().StatusCode, tt.ResponseCode)
			}
		})
	}
}

var updateBlogPostTestCases = []struct {
	Name         string
	PostID       string
	User         string
	RequestBody  map[string]string
	ResponseCode int
}{
	{
		Name: "Same Owner",
		User: "0197aaed-4a35-74da-8574-4165524a1111",
		RequestBody: map[string]string{
			"title":    "Some title",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		ResponseCode: 200,
	},
	{
		Name: "Different Owner",
		User: "0197aaed-4a35-74da-8574-4165524a2222",
		RequestBody: map[string]string{
			"title":    "Some title",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		ResponseCode: 403,
	},
	{
		Name:   "Post Does Not Exist",
		PostID: "efbfa286-ca55-4ded-a28e-9881118186c8",
		User:   "0197aaed-4a35-74da-8574-4165524a1111",
		RequestBody: map[string]string{
			"title":    "Some title",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		ResponseCode: 404,
	},
	{
		Name: "User Info Missing",
		RequestBody: map[string]string{
			"title":    "Some title",
			"status":   "DRAFT",
			"summary":  "Some summary under N chars",
			"contents": "Some really long string",
		},
		ResponseCode: 403,
	},
}

func TestUpdateBlogPostHandler(t *testing.T) {
	for _, tt := range updateBlogPostTestCases {
		t.Run(tt.Name, func(t *testing.T) {

			// InMemory implementations double as mock test implementations for unit tests
			app := App{
				Logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				BlogService: db.NewInMemoryBlogService(),
				UserService: &db.InMemoryUserService{
					Users: TestUserMap,
				},
			}

			// seed an existing post beforehand
			blog := &model.BlogPost{}
			err := app.BlogService.CreateBlogPost(context.TODO(), "0197aaed-4a35-74da-8574-4165524a1111", blog)
			if err != nil {
				t.Error(err)
			}

			postID := blog.ID
			if tt.PostID != "" {
				// override post ID for request
				postID = tt.PostID
			}

			b, _ := json.Marshal(tt.RequestBody)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/posts/%s", postID), bytes.NewReader(b))
			req.SetPathValue("id", postID)

			// set user identity
			ctx := context.WithValue(req.Context(), constant.UserIDKey, tt.User)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			app.UpdateBlogPostHandler(rr, req)
			if rr.Result().StatusCode != tt.ResponseCode {
				t.Errorf("got %d, want %d", rr.Result().StatusCode, tt.ResponseCode)
			}
		})
	}
}

var deleteBlogPostTestCases = []struct {
	Name         string
	PostID       string
	User         string
	ResponseCode int
}{
	{
		Name:         "Same Owner",
		User:         "0197aaed-4a35-74da-8574-4165524a1111",
		ResponseCode: 204,
	},
	{
		Name:         "Different Owner",
		User:         "0197aaed-4a35-74da-8574-4165524a2222",
		ResponseCode: 403,
	},
	{
		Name:         "Post Does Not Exist",
		PostID:       "efbfa286-ca55-4ded-a28e-9881118186c8",
		User:         "0197aaed-4a35-74da-8574-4165524a1111",
		ResponseCode: 404,
	},
}

func TestDeleteBlogPostHandler(t *testing.T) {
	for _, tt := range deleteBlogPostTestCases {
		t.Run(tt.Name, func(t *testing.T) {

			// InMemory implementations double as mock test implementations for unit tests
			app := App{
				Logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				BlogService: db.NewInMemoryBlogService(),
				UserService: &db.InMemoryUserService{
					Users: TestUserMap,
				},
			}

			// seed an existing post beforehand
			blog := &model.BlogPost{}
			err := app.BlogService.CreateBlogPost(context.TODO(), "0197aaed-4a35-74da-8574-4165524a1111", blog)
			if err != nil {
				t.Error(err)
			}

			postID := blog.ID
			if tt.PostID != "" {
				// override post ID for request
				postID = tt.PostID
			}

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/posts/%s", postID), nil)
			req.SetPathValue("id", postID)

			// set user identity
			ctx := context.WithValue(req.Context(), constant.UserIDKey, tt.User)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			app.DeleteBlogPostHandler(rr, req)
			if rr.Result().StatusCode != tt.ResponseCode {
				t.Errorf("got %d, want %d", rr.Result().StatusCode, tt.ResponseCode)
			}
		})
	}
}

var adminDeleteBlogPostTestCases = []struct {
	Name         string
	PostID       string
	User         string
	IsAdmin      bool
	ResponseCode int
}{
	{
		Name:         "Works Regardless as Long As Middleware Passes",
		ResponseCode: 204,
	},
	{
		Name:         "Admin - Post Does Not Exist",
		PostID:       "efbfa286-ca55-4ded-a28e-9881118186c8",
		User:         "0197aaed-4a35-74da-8574-4165524a2222",
		IsAdmin:      true,
		ResponseCode: 404,
	},
}

func TestAdminDeleteBlogPostHandler(t *testing.T) {
	for _, tt := range adminDeleteBlogPostTestCases {
		t.Run(tt.Name, func(t *testing.T) {

			// InMemory implementations double as mock test implementations for unit tests
			app := App{
				Logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
				BlogService: db.NewInMemoryBlogService(),
				UserService: &db.InMemoryUserService{
					Users: TestUserMap,
				},
			}

			// seed an existing post beforehand
			blog := &model.BlogPost{}
			err := app.BlogService.CreateBlogPost(context.TODO(), "0197aaed-4a35-74da-8574-4165524a1111", blog)
			if err != nil {
				t.Error(err)
			}

			postID := blog.ID
			if tt.PostID != "" {
				// override post ID for request
				postID = tt.PostID
			}

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/posts/%s", postID), nil)
			req.SetPathValue("id", postID)

			// set user identity
			ctx := context.WithValue(req.Context(), constant.UserIDKey, tt.User)
			ctx = context.WithValue(ctx, constant.AdminKey, tt.IsAdmin)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			app.AdminDeleteBlogPostHandler(rr, req)
			if rr.Result().StatusCode != tt.ResponseCode {
				t.Errorf("got %d, want %d", rr.Result().StatusCode, tt.ResponseCode)
			}
		})
	}
}
