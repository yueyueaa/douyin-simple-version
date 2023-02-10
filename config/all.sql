// 用户信息
// 用户表
CREATE TABLE user
(
    UID INT PRIMARY KEY,
    password VARCHAR(512) NOT NULL
);
CREATE TABLE user_info
(
    UID INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(512) NOT NULL,
    follow_count INT,
    follower_count INT
);

// 视频信息
// 视频表
CREATE TABLE video_url(
    VID INT,
    play_url CHAR(255),
    cover_url CHAR(255)
);
CREATE TABLE video_info(
    VID INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255),
    play_num INT,
    like_num INT,
    publish_time timestamp,
    author INT,
    comment_num INT
);

// 视频对用户记录的信息
// 点赞表
CREATE TABLE likes
(
    VID INT NOT NULL,
    UID INT NOT NULL,
);

// 评论表
CREATE TABLE comments
(
    VID INT NOT NULL,
    UID INT NOT NULL,
    comment_text VARCHAR(512),
    comment_time timestamp
);

// 用户对用户记录的信息
// 关注表
CREATE TABLE follows
(
    UID INT NOT NULL,
    FOLLOW_ID INT NOT NULL,
    # 可以设置一个触发器，此表被插入时自动在下方表中插入该项
);

// 粉丝表/被关注表
CREATE TABLE followers
(
    UID INT NOT NULL,
    FOLLOWER_ID INT NOT NULL,
); 