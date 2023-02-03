CREATE TABLE commits(
    ID INT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(512),
    video_id VARCHAR(512),
    commit_text VARCHAR(512),
    commit_time timestamp,
    flag INTEGER
);
INSERT INTO commits(
        user_id,
        video_id,
        commit_text,
        commit_time,
        flag
    )
VALUES('1', '1', 'test', '2020-01-11 09:53:32', '1');