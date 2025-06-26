package db

import (
	"context"
	"errors"
	"time"

	"github.com/James-D-Wood/blog-api/internal/constant"
	"github.com/James-D-Wood/blog-api/internal/model"
)

// TODO: database implementation of BlogService interface

var (
	ErrEntityNotFound = errors.New("entity not found")
)

type BlogService interface {
	FetchBlog(ctx context.Context, id string) (model.BlogPost, error)
	FetchPublishedBlogs(ctx context.Context) ([]model.BlogPost, error)
	CreateBlog(ctx context.Context, blog *model.BlogPost) error
	UpdateBlog(ctx context.Context, blog model.BlogPost) error
	DeleteBlog(ctx context.Context, id string) (model.BlogPost, error)
}

// InMemoryBlogService implements BlogService using an in process data store
type InMemoryBlogService struct {
	m map[string]model.BlogPost
}

func NewInMemoryBlogService() *InMemoryBlogService {
	return &InMemoryBlogService{m: map[string]model.BlogPost{}}
}

func (s *InMemoryBlogService) FetchBlog(ctx context.Context, id string) (model.BlogPost, error) {
	if blog, ok := s.m[id]; ok {
		return blog, nil
	}
	return model.BlogPost{}, ErrEntityNotFound
}

func (s *InMemoryBlogService) FetchPublishedBlogs(ctx context.Context) ([]model.BlogPost, error) {
	blogs := []model.BlogPost{}
	for _, blog := range s.m {
		// if blog is published
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (s *InMemoryBlogService) CreateBlog(ctx context.Context, post *model.BlogPost) error {
	// TODO: add required field validation - ie: title, description, contents

	// set generated values
	post.ID = assignUUID()
	userID, _ := ctx.Value(constant.UserIDKey).(string)
	post.AuthorID = userID
	ts := time.Now().Format(time.RFC3339)
	post.CreatedTS = ts
	post.UpdatedTS = ts

	// check that blog does not already exist
	for _, p := range s.m {
		if p.AuthorID == post.AuthorID && p.Title == post.Title {
			return errors.New("blog post already exists")
		}
	}

	s.m[post.ID] = *post
	return nil
}

func (s *InMemoryBlogService) UpdateBlog(ctx context.Context, blog model.BlogPost) error {
	return errors.New("not implemented")
}

func (s *InMemoryBlogService) DeleteBlog(ctx context.Context, id string) (model.BlogPost, error) {
	return model.BlogPost{}, errors.New("not implemented")
}
