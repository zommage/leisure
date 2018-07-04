package dbbase

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	. "github.com/ruoklive/platformSiteGateway/logs"
)

// platform site gateway users
type LeisureUsers struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT;"` // id
	Username  string    `gorm:"column:username;unique_index"`          // 用户名
	Pwd       string    `gorm:"column:pwd"`                            // 用户名密码
	Role      string    `gorm:"column:role"`                           // 角色, 如: admin, normal
	Status    int       `gorm:"column:status"`                         // status = 1, 代表可用, 2: 不可用
	Online    int       `gorm:"column:online"`                         // 用户是否在线, 1: 在线, 2 离线
	UpdatedAt time.Time `gorm:"column:updated_at"`                     // 更新时间
	CreatedAt time.Time `gorm:"column:created_at"`                     // 创建时间
}

// 用户 token, 允许一个用户名在多个浏览器登录
type LeisureUserToken struct {
	ID        int       `gorm:"column:id;primary_key;AUTO_INCREMENT;"` // id
	Username  string    `gorm:"column:username"`                       // 用户名
	Role      string    `gorm:"column:role"`                           // 角色, 如: admin, normal
	Token     string    `gorm:"column:token"`                          // token
	Expire    time.Time `gorm:"column:expire"`                         // token 的过期时间
	CreatedAt time.Time `gorm:"column:created_at"`                     // 创建时间
}

// 根据用户名进行查找
func QueryUserByUsername(username string, status int) (*LeisureUsers, error) {
	row := &LeisureUsers{}
	err := db.Dbs.Where("username = ? AND status = ?", username, status).First(row).Error
	if err != nil {
		return nil, err
	}

	return row, nil
}

// 根据token进行查找
func QueryByToken(token string) (*LeisureUserToken, error) {
	row := &LeisureUserToken{}
	err := db.Dbs.Where("token = ?", token).First(row).Error
	if err != nil {
		return nil, err
	}

	return row, nil
}

// 删除过期的 token, gorm 的删除是让查询不可见, 需要加上 unscoped 才是真正的删除
func DelExpireToken() error {
	var rows []*LeisureUserToken
	expireTime := time.Now().Add(-60 * time.Minute)

	err := db.Dbs.Where("expire <= ?", expireTime).Delete(rows).Error
	if err != nil {
		Log.Errorf("del batch err: %v", err)
		return err
	}

	// Delete record permanently with Unscoped
	err = db.Dbs.Unscoped().Where("expire <= ?", expireTime).Delete(rows).Error
	if err != nil {
		Log.Errorf("del unscoped err: %v", err)
		return err
	}

	return nil
}
