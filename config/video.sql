CREATE TABLE video_url(
    ID INT PRIMARY KEY,
    play_url CHAR(255),
    cover_url CHAR(255)
);

CREATE TABLE video_info(
    ID INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255),
    play_num INTEGER,
    like_num INTEGER,
    publish_time timestamp,
    author VARCHAR(255),
    comment_num INTEGER
);

INSERT INTO video_info(title,play_num,like_num,publish_time,author,comment_num)
    VALUES("testvideo4",100,200,'2020-01-11 09:53:32',1,300);

INSERT INTO video_info(title,play_num,like_num,publish_time,author,comment_num)
    VALUES("testvideo5",1000,2000,'2020-01-12 09:53:32',1,3000);

INSERT INTO video_info(title,play_num,like_num,publish_time,author,comment_num)
    VALUES("testvideo6",10000,20000,'2020-01-13 09:53:32',5,30000);

INSERT INTO video_url(VID,play_url,cover_url)
    VALUES(4,
    'https://prod-streaming-video-msn-com.akamaized.net/a8c412fa-f696-4ff2-9c76-e8ed9cdffe0f/604a87fc-e7bc-463e-8d56-cde7e661d690.mp4',
    'https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AAOEcdM.img'
    );

INSERT INTO video_url(VID,play_url,cover_url)
    VALUES(5,
    'https://prod-streaming-video-msn-com.akamaized.net/0b927d99-e38a-4f51-8d1a-598fd4d6ee97/3493c85c-f35a-488f-9a8f-633e747fb141.mp4',
    'https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AAOEhRG.img'
    );

INSERT INTO video_url(VID,play_url,cover_url)
    VALUES(6,
    'https://prod-streaming-video-msn-com.akamaized.net/178161a4-26a5-4f84-96d3-6acea1909a06/2213bcd0-7d15-4da0-a619-e32d522572c0.mp4',
    'https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AAOE58C.img'
    );

SELECT * FROM video_info ORDER BY ID DESC LIMIT 30;
SELECT MAX(ID) FROM video_info;

SELECT * FROM video_info WHERE ID = ?;
SELECT * FROM video_url WHERE ID = ?;
