CREATE TABLE video(
    ID INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(512),
    play_num INTEGER,
    like_num INTEGER,
    publish_time timestamp,
    author VARCHAR(512),
    commit_num INTEGER,
    play_url VARCHAR(512),
    cover_url VARCHAR(512)
);
INSERT INTO video(
        title,
        play_num,
        like_num,
        publish_time,
        author,
        commit_num,
        play_url,
        cover_url
    )
VALUES(
        'test',
        0,
        0,
        '2020-01-11 09:53:32',
        'sjx',
        0,
        'somewhere',
        'somewhere'
    );