-- Users table
CREATE TABLE
    IF NOT EXISTS users (
        user_id BIGSERIAL PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        firstname TEXT NOT NULL,
        lastname TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        birthday DATE NOT NULL,
        role TEXT NOT NULL DEFAULT 'user',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        last_online TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Blogs table
CREATE TABLE
    IF NOT EXISTS blogs (
        blog_id BIGSERIAL PRIMARY KEY,
        headline TEXT NOT NULL,
        body TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        user_id BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE
    );

-- Comments table
CREATE TABLE
    IF NOT EXISTS comments (
        comment_id BIGSERIAL PRIMARY KEY,
        body TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        blog_id BIGINT NOT NULL REFERENCES blogs (blog_id) ON DELETE CASCADE,
        user_id BIGINT NOT NULL REFERENCES users (user_id) ON DELETE CASCADE
    );

-- Indexes for better performance
CREATE INDEX idx_blogs_user_id ON blogs (user_id);

CREATE INDEX idx_comments_blog_id ON comments (blog_id);

CREATE INDEX idx_comments_user_id ON comments (user_id);