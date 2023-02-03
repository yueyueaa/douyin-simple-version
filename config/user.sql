CREATE TABLE user(
    ID INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(512),
    follow_num INTEGER,
    fans_num INTEGER,
    password VARCHAR(512),
    sex VARCHAR(512),
    token VARCHAR(512),
    other VARCHAR(512)
);
INSERT INTO user(
        name,
        follow_num,
        fans_num,
        password,
        sex,
        token,
        other
    )
VALUES ('sjx', 0, 0, 123456, 'male', 'sjx123456', '');