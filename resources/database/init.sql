begin transaction;

CREATE TABLE IF NOT EXISTS features (
    id INT PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS tags (
    id INT PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS user_banners (
    id INT PRIMARY KEY,
    content jsonb,
    is_active bool default true,
    feature_id int,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    FOREIGN KEY(feature_id) REFERENCES features(id)
);

CREATE TABLE IF NOT EXISTS user_banners_tags (
     banner_id INT,
     tag_id INT,
     FOREIGN KEY(banner_id) REFERENCES user_banners(id),
     FOREIGN KEY(tag_id) REFERENCES tags(id)
);


end transaction;