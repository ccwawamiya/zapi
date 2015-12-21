package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	_ "encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"github.com/fatih/structs"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"zapi/models"
)

// Operations about object
type VerifyController struct {
	beego.Controller
}

type params struct {
	Token     string `valid:"Required;Length(32)"`
	Openid    string `valid:"Required;Length(32)"`
	Timestamp int    `valid:"Required"`
	Gamekey   string `valid:"Required;Length(32)"`
	Sign      string `valid:"Required;Length(32)"`
}

// @Title create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
func (this *VerifyController) Post() {
	var response = map[string]interface{}{
		"result":      "10001",
		"result_desc": "Miss Params",
	}
	for index := 0; index == 0; index++ {
		valid := validation.Validation{}
		inputtime, err := strconv.Atoi(this.Input().Get("timestamp"))
		var p = params{Token: this.Input().Get("token"), Openid: this.Input().Get("openid"),
			Timestamp: inputtime, Gamekey: this.Input().Get("gamekey"),
			Sign: this.Input().Get("_sign")}
		b, err := valid.Valid(&p)
		if err != nil {
			response["result_desc"] = "sytem error"
			break
		}
		if !b {
			response["result"] = "10002"
			response["result_desc"] = ""
			for _, err := range valid.Errors {
				response["result_desc"] = err.Key + err.Message
			}
			break
		}
		/*
			om := orm.NewOrm()
			game := models.Game{Game_key: this.Input().Get("gamekey")}
			err = om.Read(&game, "Game_key")
			if err == orm.ErrNoRows {
				response["result_desc"] = "Illegal gamekey"
				break
			}
		*/
		kuaicache := new(models.Kuaicache)
		status, game := kuaicache.GetGameCache(this.Input().Get("gamekey"))
		if !status || game == nil || game["Security_key"] == nil {
			response["result_desc"] = "Illegal gamekey"
			break
		}
		if !this.checkSign(&p, game["Security_key"].(string)) {
			response["result_desc"] = "error sign "
			break
		}
		openid, status := this.checkTokenValid(this.Input().Get("token"))
		if status == false {
			response["result_desc"] = openid
			break
		}
		if openid != this.Input().Get("openid") {
			response["result_desc"] = "error  openid"
			break
		}
		response["result_desc"] = game["Name"].(string)
		break
	}

	this.Data["json"] = response
	this.ServeJson()
	//log := logs.NewLogger(10)
	//log.SetLogger("console", "")
	//log.SetLogger("file", `{"filename":"D:/test.log"}`)
	//log.Info("this is a test")
}

func (this *VerifyController) checkSign(p *params, securitykey string) bool {

	if this.getSign(p, securitykey) != p.Sign {
		return false
	}
	return true
}

func (this *VerifyController) getSign(p *params, securitykey string) string {
	//var data = "gamekey=" + p.Gamekey + "&openid=" + p.Openid +
	//	"&timestamp=" + strconv.Itoa(p.Timestamp) + "&token=" + p.Token
	data := ""
	mapp := structs.Map(p)
	delete(mapp, "Sign")
	sortk := this.sortMap(mapp)
	for _, k := range sortk {
		tmp := ""
		switch mapp[k].(type) {
		case string:
			tmp = mapp[k].(string)
		case int:
			tmp = strconv.Itoa(mapp[k].(int))
		}
		data = data + strings.ToLower(k) + "=" + tmp + "&"
	}
	data = data[0 : len(data)-1]
	m := md5.Sum([]byte(data))
	sign := hex.EncodeToString(m[:]) + securitykey
	m = md5.Sum([]byte(sign))
	sign = hex.EncodeToString(m[:])
	return sign
}

func (this *VerifyController) sortMap(m map[string]interface{}) []string {
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)
	return sorted_keys
}

func (this *VerifyController) checkTokenValid(token string) (string, bool) {
	/*
		om := orm.NewOrm()
		oauth := models.Auth_code{Code: token}
		err := om.Read(&oauth, "Code")
		if err == orm.ErrNoRows {
			return "token is not exists", false
		}
	*/

	kuaicache := new(models.Kuaicache)
	status, oauth := kuaicache.GetTokenCache(token)
	if !status || oauth == nil || oauth["Puid"] == nil {
		return "token is not exists", false
	}
	var puid int
	switch oauth["Puid"].(type) {
	case string:
		puid, _ = strconv.Atoi(oauth["Puid"].(string))
		break
	case int:
		puid = oauth["Puid"].(int)
		break
	}
	ustatus, users := kuaicache.GetUserCacheByPuid(puid)
	if !ustatus || users == nil || users["Guid"] == nil {
		return "error token", false
	}

	m := md5.Sum([]byte(users["Platform_code"].(string) + ":" + users["Guid"].(string)))
	openid := hex.EncodeToString(m[:])
	return openid, true
}

func (this *VerifyController) Finish() {
	logarr := map[string]interface{}{
		"ip":       "error",
		"params":   "error",
		"response": "error",
	}
	logarr["ip"] = this.Ctx.Request.RemoteAddr
	logarr["params"] = this.Ctx.Request.PostForm
	logarr["response"] = this.Data["json"]
	log := logs.NewLogger(10000)
	logdir := beego.AppConfig.String("logdir") + "/verify/" + time.Now().Format("200601")
	_, errors := os.Stat(logdir)
	if errors != nil && os.IsNotExist(errors) {
		os.MkdirAll(logdir, 0766)
	}
	logfilename := logdir + "/" + time.Now().Format("20060102") + ".log"
	logsetting := map[string]interface{}{
		"filename": logfilename,
	}
	logset, _ := json.Marshal(logsetting)
	log.SetLogger("file", string(logset))
	content, err := json.Marshal(logarr)
	if err == nil {
		log.Info(string(content))
	}
}
