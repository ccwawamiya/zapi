package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/fatih/structs"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type Kuaicache struct {
	gameKey string
}

var c redis.Conn

var _prefix = ""

func init() {
	d, err := redis.Dial("tcp", beego.AppConfig.String("redisServers"))
	if err == nil {
		c = d
	} else {
		panic("no connect redis")
	}
	_prefix = beego.AppConfig.String("redisPrefix")
}

func (this *Kuaicache) GetGameCache(gameKey string) (bool, map[string]interface{}) {
	info, err := redis.String(c.Do("HGET", _prefix+"game_id_", gameKey))
	if err == nil {
		status, gameinfo := this.GetGameCacheById(info)
		if status {
			return true, gameinfo
		}
	}
	om := orm.NewOrm()
	game := Game{Game_key: gameKey}
	err = om.Read(&game, "Game_key")
	if err != orm.ErrNoRows {
		mapp := structs.Map(game)
		return true, mapp
	}
	return false, nil
}

func (this *Kuaicache) GetGameCacheById(gameId string) (bool, map[string]interface{}) {
	sinfo := make(map[string]interface{})
	info, err := redis.StringMap(c.Do("HGETALL", _prefix+"game_info:"+gameId))
	if err == nil && len(info) > 0 {
		for key, val := range info {
			sinfo[strings.Title(key)] = val
		}
		return true, sinfo
	}
	return false, nil
}

func (this *Kuaicache) GetTokenCache(token string) (bool, map[string]interface{}) {
	sinfo := make(map[string]interface{})
	info, err := redis.StringMap(c.Do("HGETALL", _prefix+"oauthtoken:"+token))
	if err == nil && len(info) > 0 {
		for key, val := range info {
			sinfo[strings.Title(key)] = val
		}
		return true, sinfo
	}
	om := orm.NewOrm()
	oauth := Auth_code{Code: token}
	err = om.Read(&oauth, "Code")
	if err != orm.ErrNoRows {
		mapp := structs.Map(oauth)
		return true, mapp
	}
	return false, nil
}

func (this *Kuaicache) GetUserCacheByPuid(puid int) (bool, map[string]interface{}) {
	sinfo := make(map[string]interface{})
	info, err := redis.StringMap(c.Do("HGETALL", _prefix+"thirdpuid:"+string(puid)))
	if err == nil && len(info) > 0 {
		for key, val := range info {
			sinfo[strings.Title(key)] = val
		}
		return true, sinfo
	}

	om := orm.NewOrm()
	users := Users{Puid: puid}
	err = om.Read(&users, "Puid")
	if err != orm.ErrNoRows {
		mapp := structs.Map(users)
		return true, mapp
	}
	return false, nil
}
