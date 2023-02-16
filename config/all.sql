CREATE TABLE user
(
    UID INT PRIMARY KEY,
    Password VARCHAR(512) NOT NULL
);
CREATE TABLE user_info
(
    UID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(512) NOT NULL,
    follow_count INT,
    follower_count INT
);

CREATE TABLE video_url(
    VID INT,
    play_url VARCHAR(255),
    cover_url VARCHAR(255)
);
CREATE TABLE video_info(
    VID INT PRIMARY KEY AUTO_INCREMENT,
    author_id INT,
    Title VARCHAR(255),
    favorite_count INT,
    comment_count INT,
    publish_date timestamp
);

CREATE TABLE favorites
(
    VID INT NOT NULL,
    UID INT NOT NULL,
    Flag INT NOT NULL
);

CREATE TABLE comments
(
    CID INT PRIMARY KEY AUTO_INCREMENT,
    VID INT NOT NULL,
    UID INT NOT NULL,
    Content VARCHAR(512),
    create_date timestamp
);

CREATE TABLE follows
(
    UID INT NOT NULL,
    follow_id INT NOT NULL,
    Flag INT NOT NULL
);

CREATE TABLE followers
(
    UID INT NOT NULL,
    follower_id INT NOT NULL,
    Flag INT NOT NULL
);

