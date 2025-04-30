CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       content TEXT NOT NULL,
                       avatar_url TEXT,
                       imgs_urls TEXT[],  -- PostgreSQL array type
                       author TEXT NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       last_commented TIMESTAMP,
                       deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                          parent_comment_id INTEGER REFERENCES comments(id) ON DELETE CASCADE,
                          avatar_url TEXT,
                          imgs_urls TEXT[],  -- PostgreSQL array type
                          content TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          author TEXT NOT NULL
);

CREATE TABLE sessions (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          avatar_url TEXT,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          expires_at TIMESTAMP NOT NULL
);