package postgres

import (
    "context"
    "database/sql"
    "errors"
    "BAZ/internal/db"
    "BAZ/internal/domain"
)

type CommentRepository struct {
    queries *db.Queries
}

func NewCommentRepository(queries *db.Queries) *CommentRepository {
    return &CommentRepository{queries: queries}
}

// Create erstellt einen neuen Comment
func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
    dbComment, err := r.queries.CreateComment(ctx, db.CreateCommentParams{
        Body:   comment.Body,
        BlogID: comment.BlogID,
        UserID: comment.UserID,
    })
    if err != nil {
        return err
    }
    
    comment.CommentID = dbComment.CommentID
    comment.CreatedAt = dbComment.CreatedAt
    comment.UpdatedAt = dbComment.UpdatedAt
    
    return nil
}

// FindByID findet einen Comment anhand seiner ID
func (r *CommentRepository) FindByID(ctx context.Context, id int64) (*domain.Comment, error) {
    dbComment, err := r.queries.GetCommentByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    
    return r.toDomain(dbComment), nil
}

// FindByBlogID findet alle Comments eines Blogs mit Paginierung
func (r *CommentRepository) FindByBlogID(ctx context.Context, blogID int64, page, limit int) ([]domain.Comment, int64, error) {
    offset := (page - 1) * limit
    
    // Comments abrufen
    dbComments, err := r.queries.ListCommentsByBlogID(ctx, db.ListCommentsByBlogIDParams{
        BlogID: blogID,
        Offset: int32(offset),
        Limit:  int32(limit),
    })
    if err != nil {
        return nil, 0, err
    }
    
    // Gesamtanzahl abrufen
    total, err := r.queries.CountCommentsByBlogID(ctx, blogID)
    if err != nil {
        return nil, 0, err
    }
    
    comments := make([]domain.Comment, 0, len(dbComments))
    for _, dbComment := range dbComments {
        comments = append(comments, *r.toDomain(dbComment))
    }
    
    return comments, total, nil
}

// Update aktualisiert einen bestehenden Comment
func (r *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
    dbComment, err := r.queries.UpdateComment(ctx, db.UpdateCommentParams{
        CommentID: comment.CommentID,
        Body:      comment.Body,
        UserID:    comment.UserID,
    })
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.New("comment not found or you are not the author")
        }
        return err
    }
    
    comment.Body = dbComment.Body
    comment.UpdatedAt = dbComment.UpdatedAt
    
    return nil
}

// Delete löscht einen Comment
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
    // Ähnliches Problem wie beim Blog: Löschen ohne userID
    err := r.queries.DeleteCommentByID(ctx, id)
    if err != nil && errors.Is(err, sql.ErrNoRows) {
        return errors.New("comment not found")
    }
    return err
}

// DeleteByBlogID löscht alle Comments eines Blogs (für Cascade Delete)
func (r *CommentRepository) DeleteByBlogID(ctx context.Context, blogID int64) error {
    return r.queries.DeleteCommentsByBlog(ctx, blogID)
}

// DeleteByUserID löscht alle Comments eines Users (für Cascade Delete)
func (r *CommentRepository) DeleteByUserID(ctx context.Context, userID int64) error {
    return r.queries.DeleteCommentsByUser(ctx, userID)
}

// toDomain konvertiert sqlc Comment zu domain.Comment
func (r *CommentRepository) toDomain(c db.Comment) *domain.Comment {
    return &domain.Comment{
        CommentID: c.CommentID,
        Body:      c.Body,
        CreatedAt: c.CreatedAt,
        UpdatedAt: c.UpdatedAt,
        BlogID:    c.BlogID,
        UserID:    c.UserID,
    }
}