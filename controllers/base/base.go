package base

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zommage/leisure/common"
	. "github.com/zommage/leisure/logs"
	models "github.com/zommage/leisure/models"
)

/* 发往前台的公共接口
*  statusCode : http的状态码
*  errCode: 错误码
*  Msg: 信息
 */
func WebResp(c *gin.Context, statusCode, errCode int, data interface{}, Msg string) {
	respMap := map[string]interface{}{"code": errCode, "msg": Msg, "data": data}
	c.JSON(statusCode, respMap)
	return
}

// 针对请求进行鉴权与签名校验
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 判断是否需要 token 校验
		if common.AuthSwitch != false {
			c.Next()
			return
		}

		var err error
		urlArr := strings.Split(c.Request.URL.Path, "?")
		preUrl := urlArr[0]

		// 如果是过滤的路由检测接口直接跳过, 直接跳到 controllers
		_, exist := common.RouterFilterMap[preUrl]
		if exist {
			c.Next()
			return
		}

		// 鉴权
		statusCode, errCode, err := AuthFunc(c)
		if err != nil {
			tmpStr := fmt.Sprintf("auth fail: %v", err)
			Log.Errorf(tmpStr)
			WebResp(c, statusCode, errCode, nil, tmpStr)
			c.Abort()
			return
		}

		c.Next()
		return
	}
}

// 鉴权函数
func AuthFunc(c *gin.Context) (int, int, error) {
	token := c.GetHeader("Token")

	if token == "" {
		tmpStr := fmt.Sprintf("token is nil")
		return 401, 401, errors.New(tmpStr)
	}

	// 对 token 进行解析
	statusCode, errCode, err := ParseToken(token)
	if err != nil {
		return statusCode, errCode, err
	}

	return 0, 0, nil
}

// 对 token 进行校验
func ParseToken(token string) (int, int, error) {
	// 查询该 token 在数据库中是否存在
	tokenRow, err := models.QueryByToken(token)
	if err != nil {
		return 401, 777, common.TokenNotExist
	}

	// 判断 token 是否已经过期
	if tokenRow.Expire.Before(time.Now()) != false {
		return 401, 777, common.TokenExprire
	}

	return 0, 0, nil
}

// 获取签名的公共参数
func ComSigParam(c *gin.Context) (map[string]interface{}, string, error) {
	params := make(map[string]interface{})

	TimeStamp := c.Query("TimeStamp")
	if TimeStamp == "" {
		return nil, "", fmt.Errorf("signature timestamp is nil")
	}

	SignatureNonce := c.Query("SignatureNonce")
	if SignatureNonce == "" {
		return nil, "", fmt.Errorf("signature noce is nil")
	}

	// 签名
	signature := c.Query("Signature")
	if signature == "" {
		return nil, "", fmt.Errorf("signature is nil")
	}

	// 如果签名中带有 + 号, url会将 + 解析成空格
	signature = strings.Replace(signature, " ", "+", -1)

	// 增加公共参数, 时间格式为YYYY-MM-DDThh:mm:ssZ,例如，2014-11-11T12:00:00Z
	params["TimeStamp"] = TimeStamp

	// 随机字符串
	params["SignatureNonce"] = SignatureNonce

	return params, signature, nil
}
