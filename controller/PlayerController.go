package controller

import (
	"fmt"
	"mapDemo/common"
	"mapDemo/model"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("a_secret_crect")

func NewPlayer(c *gin.Context) {
	newPlayerRedis(c)
}

func newPlayerRedis(c *gin.Context) {
	baseReq := &model.NewPlayerBaseReq{}
	c.BindJSON(baseReq)
	if baseReq.Jwt != nil {
		// 解析jwt，写入redis
	} else {
		// 用code获取openId
		openId := baseReq.Code
		// 签发jwt
		jwt, err := newJwt(openId)
		if err != nil {
			c.JSON(200, gin.H{"code": 200, "data": *baseReq.Code, "msg": "新玩家接入时jwt签发失败"})
			return
		}
		// TODO：直接更新，如果玩家在线，是不是无法连接回房间内
		err = common.LocalRedisClient.UpdatePlayer(openId, &map[string]interface{}{
			"PlayerType": 0,
			"RoomId":     "",
			// "PlayerOnline": 1,
		})
		if err != nil {
			c.JSON(200, gin.H{"code": 200, "data": *baseReq.Code, "msg": "玩家信息写入redis失败"})
		}
		c.JSON(200, gin.H{"code": 200, "data": *baseReq.Code + " Jwt: " + *jwt, "msg": "新玩家接入"})
	}
}

// 玩家后续隔一段时间向服务器发起更新位置请求
func PlayerUpdatePos(c *gin.Context) {
	req := &model.PlayerUpdatePosBaseReq{}
	c.BindJSON(req)
	uuid, err := verifyJwtUuid(req.Jwt)
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "updatePos verifyJwt失败"})
		return
	}
	// 更新Pos
	err = common.LocalRedisClient.UpdatePos(uuid, req.Pos)
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "updatePos Redis更新pos失败"})
		return
	}
	// 在房间内
	if req.Type != 0 {
		// TODO：获取全部玩家坐标，并判断可见性
	} else {
		c.JSON(200, gin.H{"code": 100, "data": req.Pos, "msg": "更新玩家坐标成功，玩家不在房间内"})
		return
	}
}

// PlayerJoinRoom 玩家申请加入房间
func PlayerJoinRoom(c *gin.Context) {
	req := &model.PlayerJoinRoomBaseReq{}
	c.BindJSON(req)
	uuid, err := verifyJwtUuid(req.Jwt)
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom verifyJwt失败"})
		return
	}
	// 查询redis缓存字段，uuid RoomId，如果redis查不到，需要重新登录
	fileds := &[]string{"PlayerType", "RoomId"}
	vals, err := common.LocalRedisClient.GetPlayerInfoByField(uuid, fileds)
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom 查询用户信息失败"})
		return
	}
	fmt.Print(vals)
	if (*vals)[0] == nil {
		// 玩家不在线
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom 玩家缓存丢失，请重新登录"})
		return
	}
	// 判断RoomId是否为空 or 为无效房间号
	if (*vals)[1] == nil || (*vals)[1] == "" {
		// 无效房间号
	} else {
		// 有效房间号，退出之前房间
		rStr, ok := (*vals)[1].(string)
		if !ok {
			c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom 房间号类型断言失败"})
			return
		}
		if err := common.LocalRedisClient.UpdateRoom(uuid, &rStr, 0); err != nil {
			c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom 房间退出失败"})
			return
		}
	}
	// 更新玩家信息
	err = common.LocalRedisClient.UpdatePlayer(uuid, &map[string]interface{}{
		"RoomId":     req.RoomId,
		"PlayerType": 1,
	})
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom 玩家信息修改失败"})
		return
	}
	// 房间加入玩家uuid
	err = common.LocalRedisClient.UpdateRoom(uuid, req.RoomId, 1)
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerJoinRoom 房间信息更新失败"})
		return
	}
	// TODO：获取房间内玩家坐标返回
}

// PlayerQuitRoom 玩家退出房间
func PlayerQuitRoom(c *gin.Context) {
	req := &model.PlayerQuitRoomBaseReq{}
	c.BindJSON(req)
	uuid, err := verifyJwtUuid(req.Jwt)
	if err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerQuitRoom verifyJwt失败"})
		return
	}
	if err := common.LocalRedisClient.UpdateRoom(uuid, req.RoomId, 0); err != nil {
		c.JSON(200, gin.H{"code": 100, "data": err, "msg": "PlayerQuitRoom 房间退出失败"})
		return
	}
	err = common.LocalRedisClient.UpdatePlayer(uuid, &map[string]interface{}{
		"RoomId":     "",
		"PlayerType": 0,
	})
}
