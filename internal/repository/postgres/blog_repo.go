package postgres

import (
    "context"
    "database/sql"
    "errors"
    "BAZ/internal/db"
    "BAZ/internal/domain"
)

type BlogRepository struct {
    queries *db.Queries
}

func NewBlogRepository(queries *db.Queries) *BlogRepository {
    return &BlogRepository{queries: queries}
}

// Create erstellt einen neuen Blog
func (r *BlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
    dbBlog, err := r.queries.CreateBlog(ctx, db.CreateBlogParams{
        Headline: blog.Headline,
        Body:     blog.Body,
        UserID:   blog.UserID,
    })
    if err != nil {
        return err
    }
    
    blog.BlogID = dbBlog.BlogID
    blog.CreatedAt = dbBlog.CreatedAt
    blog.UpdatedAt = dbBlog.UpdatedAt
    
    return nil
}

// FindByID findet einen Blog anhand seiner ID
func (r *BlogRepository) FindByID(ctx context.Context, id int64) (*domain.Blog, error) {
    dbBlog, err := r.queries.GetBlogByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Blog nicht gefunden
        }
        return nil, err
    }
    
    return r.toDomain(dbBlog), nil
}

// FindAll findet alle Blogs mit Paginierung
func (r *BlogRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Blog, int64, error) {
    offset := (page - 1) * limit
    
    // Blogs abrufen
    dbBlogs, err := r.queries.ListBlogs(ctx, db.ListBlogsParams{
        Offset: int32(offset),
        Limit:  int32(limit),
    })
    if err != nil {
        return nil, 0, err
    }
    
    // Gesamtanzahl abrufen
    total, err := r.queries.CountBlogs(ctx)
    if err != nil {
        return nil, 0, err
    }
    
    blogs := make([]domain.Blog, 0, len(dbBlogs))
    for _, dbBlog := range dbBlogs {
        blogs = append(blogs, *r.toDomain(dbBlog))
    }
    
    return blogs, total, nil
}

// Update aktualisiert einen bestehenden Blog
func (r *BlogRepository) Update(ctx context.Context, blog *domain.Blog) error {
    dbBlog, err := r.queries.UpdateBlog(ctx, db.UpdateBlogParams{
        BlogID:   blog.BlogID,
        Headline: blog.Headline,
        Body:     blog.Body,
        UserID:   blog.UserID,
    })
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.New("blog not found or you are not the author")
        }
        return err
    }
    
    blog.Headline = dbBlog.Headline
    blog.Body = dbBlog.Body
    blog.UpdatedAt = dbBlog.UpdatedAt
    
    return nil
}

// Delete löscht einen Blog
func (r *BlogRepository) Delete(ctx context.Context, id int64) error {
    err := r.queries.DeleteBlogByID(ctx, id)
    if err != nil && errors.Is(err, sql.ErrNoRows) {
        return errors.New("blog not found")
    }
    return err
}

// DeleteByUserID löscht alle Blogs eines Users (für Cascade Delete)
func (r *BlogRepository) DeleteByUserID(ctx context.Context, userID int64) error {
    return r.queries.DeleteBlogsByUser(ctx, userID)
}

// toDomain konvertiert sqlc Blog zu domain.Blog
func (r *BlogRepository) toDomain(b db.Blog) *domain.Blog {
    return &domain.Blog{
        BlogID:    b.BlogID,
        Headline:  b.Headline,
        Body:      b.Body,
        CreatedAt: b.CreatedAt,
        UpdatedAt: b.UpdatedAt,
        UserID:    b.UserID,
    }
}