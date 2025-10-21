DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS blog_posts;

CREATE TABLE blog_posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    blog_post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_blog_post
        FOREIGN KEY (blog_post_id)
        REFERENCES blog_posts(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_comments_post_id_created_at ON comments(blog_post_id, created_at DESC);
