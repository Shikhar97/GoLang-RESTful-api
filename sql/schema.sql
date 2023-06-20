CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    birthday INTEGER NOT NULL,
    avatar TEXT
);
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    post_date INTEGER NOT NULL,
    author_id INTEGER,
    FOREIGN KEY(author_id) REFERENCES users(id)
);
CREATE TABLE likes (
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    like_date INTEGER NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(post_id) REFERENCES posts(id),
    PRIMARY KEY(post_id, user_id)
);
