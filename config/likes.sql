CREATE TABLE likes(
    ID INT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(512),
    video_id VARCHAR(512),
    flag INTEGER
);
INSERT INTO likes(user_id, video_id, flag)
VALUES('1', '1', '1');