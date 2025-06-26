package db

import (
	"context"
	"errors"
	"time"

	"github.com/James-D-Wood/blog-api/internal/model"
)

// TODO: database implementation of BlogService interface

var (
	ErrEntityNotFound = errors.New("entity not found")
)

type BlogService interface {
	FetchBlogPost(ctx context.Context, id string) (model.BlogPost, error)
	FetchPublishedBlogPosts(ctx context.Context) ([]model.BlogPost, error)
	CreateBlogPost(ctx context.Context, userID string, blog *model.BlogPost) error
	UpdateBlogPost(ctx context.Context, newVersion *model.BlogPost, previousVersion *model.BlogPost) error
	DeleteBlogPost(ctx context.Context, id string) error
}

// InMemoryBlogService implements BlogService using an in process data store
type InMemoryBlogService struct {
	m map[string]model.BlogPost
}

func NewInMemoryBlogService() *InMemoryBlogService {
	return &InMemoryBlogService{m: map[string]model.BlogPost{}}
}

func (s *InMemoryBlogService) FetchBlogPost(ctx context.Context, id string) (model.BlogPost, error) {
	if blog, ok := s.m[id]; ok {
		return blog, nil
	}
	return model.BlogPost{}, ErrEntityNotFound
}

func (s *InMemoryBlogService) FetchPublishedBlogPosts(ctx context.Context) ([]model.BlogPost, error) {
	blogs := []model.BlogPost{}
	for _, blog := range s.m {
		if blog.Status == model.PUBLISHED {
			blogs = append(blogs, blog)
		}
	}
	return blogs, nil
}

func (s *InMemoryBlogService) CreateBlogPost(ctx context.Context, userID string, post *model.BlogPost) error {
	// TODO: add required field validation - ie: title, description, contents

	// set generated values
	post.ID = assignUUID()
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

func (s *InMemoryBlogService) UpdateBlogPost(ctx context.Context, newVersion *model.BlogPost, previousVersion *model.BlogPost) error {
	// TODO: improve required field validations
	if newVersion.Title == "" {
		return errors.New("title cannot be empty")
	}

	ts := time.Now().Format(time.RFC3339)
	previousVersion.UpdatedTS = ts

	// only allow certain fields to be updated
	previousVersion.Title = newVersion.Title
	previousVersion.Summary = newVersion.Summary
	previousVersion.Contents = newVersion.Contents
	previousVersion.Status = newVersion.Status

	// automatically add publication date
	if newVersion.Status != previousVersion.Status && newVersion.Status == model.PUBLISHED {
		previousVersion.PublishedTS = ts
	}

	s.m[previousVersion.ID] = *previousVersion

	return nil
}

func (s *InMemoryBlogService) DeleteBlogPost(ctx context.Context, id string) error {
	delete(s.m, id)
	return nil
}
