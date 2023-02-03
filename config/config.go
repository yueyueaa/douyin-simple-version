package config

// 数据库用户名和密码
var (
	MySQLUser = "team"
	MySQLPwd  = "team"
	MySQLIP   = "60.255.139.184:20022"
	// 最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	MySQLMaxOpenConns = 20000
	// 闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	MySQLMaxIdleConns = 0
)