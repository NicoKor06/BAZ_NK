CREATE TABLE IF NOT EXISTS users (
     user_id     BIGSERIAL PRIMARY KEY,
     username    TEXT NOT NULL UNIQUE,
     firstname   TEXT NOT NULL,
    lastname    TEXT NOT NULL,
     email       TEXT NOT NULL UNIQUE,
     password    TEXT NOT NULL,
     birthday    DATE NOT NULL,
     role        TEXT NOT NULL DEFAULT 'user',
     created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_online TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_providers (
      id BIGSERIAL PRIMARY KEY,
      user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_id)
    );


CREATE TABLE IF NOT EXISTS blogs (
     blog_id    BIGSERIAL PRIMARY KEY,
     headline   TEXT NOT NULL,
     body       TEXT NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id    BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    comment_id BIGSERIAL PRIMARY KEY,
    body       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    blog_id    BIGINT NOT NULL REFERENCES blogs(blog_id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_blogs_user_id ON blogs(user_id);
CREATE INDEX idx_comments_blog_id ON comments(blog_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_user_providers_user_id ON user_providers(user_id);
CREATE INDEX idx_user_providers_provider ON user_providers(provider, provider_id);