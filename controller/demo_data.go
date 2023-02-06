package controller

var DemoVideos = []Video{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
	{
		Id:            2,
		Author:        DemoUser,
		PlayUrl:       "https://prod-streaming-video-msn-com.akamaized.net/a8c412fa-f696-4ff2-9c76-e8ed9cdffe0f/604a87fc-e7bc-463e-8d56-cde7e661d690.mp4",
		CoverUrl:      "https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AAOEcdM.img",
		FavoriteCount: 100,
		CommentCount:  200,
		IsFavorite:    false,
	},
	{
		Id:            3,
		Author:        DemoUser,
		PlayUrl:       "https://prod-streaming-video-msn-com.akamaized.net/0b927d99-e38a-4f51-8d1a-598fd4d6ee97/3493c85c-f35a-488f-9a8f-633e747fb141.mp4",
		CoverUrl:      "https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AAOEhRG.img",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    true,
	},
	{
		Id:            4,
		Author:        DemoUser,
		PlayUrl:       "https://prod-streaming-video-msn-com.akamaized.net/178161a4-26a5-4f84-96d3-6acea1909a06/2213bcd0-7d15-4da0-a619-e32d522572c0.mp4",
		CoverUrl:      "https://img-s-msn-com.akamaized.net/tenant/amp/entityid/AAOE58C.img",
		FavoriteCount: 150,
		CommentCount:  530,
		IsFavorite:    true,
	},
}

var DemoComments = []Comment{
	{
		Id:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoUser = User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}
