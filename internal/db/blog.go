package db

import (
	"errors"

	"github.com/James-D-Wood/blog-api/internal/model"
)

// TODO database integration

var (
	ErrEntityNotFound = errors.New("entity not found")
)

type BlogService interface {
	FetchBlog(id string) (*model.Blog, error)
	FetchPublishedBlogs() ([]*model.Blog, error)
	CreateBlog(blog *model.Blog) error
	UpdateBlog(blog *model.Blog) error
	DeleteBlog(id string) (*model.Blog, error)
}

// InMemoryBlogService implements BlogService using an in process data store
type InMemoryBlogService struct {
	m map[string]*model.Blog
}

func (s *InMemoryBlogService) FetchBlog(id string) (*model.Blog, error) {
	if blog, ok := s.m[id]; ok {
		return blog, nil
	}
	return nil, ErrEntityNotFound
}

func (s *InMemoryBlogService) FetchPublishedBlogs() ([]*model.Blog, error) {
	var blogs []*model.Blog
	for _, blog := range s.m {
		// if blog is published
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (s *InMemoryBlogService) CreateBlog(blog *model.Blog) error {
	return errors.New("not implemented")
}

func (s *InMemoryBlogService) UpdateBlog(blog *model.Blog) error {
	return errors.New("not implemented")
}

func (s *InMemoryBlogService) DeleteBlog(id string) (*model.Blog, error) {
	return nil, errors.New("not implemented")
}
