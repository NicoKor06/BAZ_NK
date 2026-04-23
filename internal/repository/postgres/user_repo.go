package postgres

import (
    "context"
    "database/sql"
    "errors"
	_"time"
    "BAZ/internal/db"
    "BAZ/internal/domain"
)

type UserRepository struct {
    queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
    return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    dbUser, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Username:   user.Username,
        Firstname:  user.Firstname,
        Lastname:   user.Lastname,
        Email:      user.Email,
        Password:   user.Password,
        Birthday:   user.Birthday, // domain.User hat time.Time, sqlc erwartet time.Time
        Role:       user.Role,
        LastOnline: user.LastOnline,
    })
    if err != nil {
        return err
    }
    
    user.UserID = dbUser.UserID
    user.CreatedAt = dbUser.CreatedAt
    user.UpdatedAt = dbUser.UpdatedAt
    user.LastOnline = dbUser.LastOnline
    
    return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    dbUser, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return r.toDomain(dbUser), nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
    dbUser, err := r.queries.GetUserByUsername(ctx, username)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return r.toDomain(dbUser), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    dbUser, err := r.queries.GetUserByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, err
    }
    return r.toDomain(dbUser), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    dbUser, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
        UserID:    user.UserID,
        Firstname: user.Firstname,
        Lastname:  user.Lastname,
        Email:     user.Email,
        Birthday:  user.Birthday,
    })
    if err != nil {
        return err
    }
    
    // Update des domain.Users mit den zurückgegebenen Werten
    user.Firstname = dbUser.Firstname
    user.Lastname = dbUser.Lastname
    user.Email = dbUser.Email
    user.Birthday = dbUser.Birthday
    user.UpdatedAt = dbUser.UpdatedAt
    
    return nil
}

func (r *UserRepository) UpdateLastOnline(ctx context.Context, id int64) error {
    return r.queries.UpdateLastOnline(ctx, id)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
    return r.queries.DeleteUser(ctx, id)
}

// toDomain konvertiert db.User → domain.User
func (r *UserRepository) toDomain(u db.User) *domain.User {
    return &domain.User{
        UserID:     u.UserID,
        Username:   u.Username,
        Firstname:  u.Firstname,
        Lastname:   u.Lastname,
        Email:      u.Email,
        Password:   u.Password,
        Birthday:   u.Birthday,
        Role:       u.Role,
        CreatedAt:  u.CreatedAt,
        UpdatedAt:  u.UpdatedAt,
        LastOnline: u.LastOnline,
    }
}