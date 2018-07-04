package users

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zommage/leisure/common"
	"github.com/zommage/leisure/controllers/base"
	. "github.com/zommage/leisure/logs"
	models "github.com/zommage/leisure/models"
)

type LoginResp struct {
	User  string `json:"user"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

// 用户登录
func Login(c *gin.Context) {
	tokenBytes, err := c.GetRawData()
	if err != nil {
		tmpStr := fmt.Sprintf("get token err: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	tokenReq := &common.LoginToken{}
	err = json.Unmarshal(tokenBytes, &tokenReq)
	if err != nil {
		tmpStr := fmt.Sprintf("fmt json token err: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	// 对 token 进行 ras s1 对齐方式 私钥解密
	loginBytes, err := common.RsaS1Decrypt(tokenReq.Token)
	if err != nil {
		tmpStr := fmt.Sprintf("descrypt s1 token err: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	loginReq := common.LoginReq{}
	err = json.Unmarshal(loginBytes, &loginReq)
	if err != nil {
		tmpStr := fmt.Sprintf("fmt login json err: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	err = common.NumLetterLine(6, 12, loginReq.User)
	if err != nil {
		tmpStr := fmt.Sprintf("username invalid: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	row, err := models.QueryUserByUsername(loginReq.User, 1)
	if err != nil {
		tmpStr := fmt.Sprintf("user not exist or user is unuse")
		Log.Infof("query by user name err: %v", err)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	// 对数据库中用户的密码进行解密
	pwdBytes, err := common.RsaS1Decrypt(row.Pwd)
	if err != nil {
		tmpStr := fmt.Sprintf("descrypt user pwd err: %v", err)
		Log.Info(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	if string(pwdBytes) != loginReq.Pwd {
		tmpStr := fmt.Sprintf("user pwd not match")
		Log.Infof(tmpStr)
		base.WebResp(c, 400, 400, nil, tmpStr)
		return
	}

	// 更新用户信息
	nowTime := time.Now()
	row.UpdatedAt = nowTime
	row.Online = 1

	err = models.UpdateDbs(row)
	if err != nil {
		tmpStr := fmt.Sprintf("user update err: %v", err)
		Log.Infof(tmpStr)
		base.WebResp(c, 500, 500, nil, tmpStr)
		return
	}

	// 创建用户 token
	tokenRes := &common.TokenContent{}
	tokenRes.User = loginReq.User
	tokenRes.Role = row.Role
	contentBytes, _ := json.Marshal(tokenRes)
	encryMsg, err := common.RsaEncrypt(contentBytes)
	if err != nil {
		tmpStr := fmt.Sprintf("create user token err: %v", err)
		Log.Infof(tmpStr)
		base.WebResp(c, 500, 500, nil, tmpStr)
		return
	}

	// 插入用户记录表
	historyRow := &models.LeisureUserToken{}
	historyRow.Username = loginReq.User
	historyRow.Role = row.Role
	historyRow.Token = encryMsg
	historyRow.Expire = nowTime.Add(25 * time.Hour)
	historyRow.CreatedAt = nowTime
	err = models.InsertDbs(historyRow)
	if err != nil {
		tmpStr := fmt.Sprintf("insert user token err: %v", err)
		Log.Infof(tmpStr)
		base.WebResp(c, 500, 500, nil, tmpStr)
		return
	}

	// 删除数据库中过期的 token
	go models.DelExpireToken()

	resp := &LoginResp{
		User:  row.Username,
		Role:  row.Role,
		Token: encryMsg,
	}

	base.WebResp(c, 200, 200, resp, common.Success)
	return
}

// 健康检查
func HealthCheck(c *gin.Context) {
	base.WebResp(c, 200, 200, nil, common.Success)
	return
}


