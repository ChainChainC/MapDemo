package model

type Player struct {
	// ---玩家信息
	Name string
	// uuid 或者其它能够标识玩家唯一性数据
	// TODO: string 或 int待定
	Uuid IdentifyType
	// 玩家身份
	PlayerType int8
	// 玩家坐标
	PlayerPos Pos
	// -----房间信息
	// 房间号 和玩家uuid类似, string 或 int
	RoomId IdentifyType
	// 玩家是否在房间内
	InRoom bool
	// TODO，玩家在线状态，下线一段时间后需要从Map中清除玩家
	PlayerOnline bool
	// 玩家token
	PlayerJwt string
}

// NewPlayerReq 请求体
type NewPlayerReq struct {
	// ---玩家信息
	Name string       `json:"nickName"`
	Uuid IdentifyType `json:"openId"`
}

type PlayerUpdatePos struct {
	Name      string       `json:"nickName"`
	Uuid      IdentifyType `json:"openId"`
	PlayerPos Pos
}
