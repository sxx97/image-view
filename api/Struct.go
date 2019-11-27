package api


type Token struct {
	Token string `json:"token"`
}

type User struct {
	Id int `bson:"id", json:"id"`
	Account string `bson:"account",json:"account"`
	Email string `bson:"email",json:"email"`
	Password string `bson:"password",json:"password"`
}

type Image struct {
	ID  interface{} `id,omitempty`         // 简写bson映射口
	Alt string        `json:"alt",bson:"alt"` // bson和json映射
	Src string        `json:"src",bson:"src"`// 属性名 为全小写的key
	FullSrc string  `json:"fullSrc",bson:"fullSrc"`
	UserId int `json:"userId",bson:"userId"`
}

type FeedbackNote struct {
	Note string `json:"note", bson:"note"`
	ID interface{} `_id,omitempty`
	UserId int `json:"userId", bson:"userId"`
}

type ResponseResult struct  {
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}


