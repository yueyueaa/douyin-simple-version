CREATE TABLE follows(
    ID INT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(512),
    follow_id VARCHAR(512),
    flag INTEGER
);
INSERT INTO follows(user_id, follow_id, flag)
VALUES('1', '2', 1);