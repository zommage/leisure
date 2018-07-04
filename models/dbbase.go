package dbbase

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zommage/leisure/conf"
	. "github.com/zommage/leisure/logs"
)

type Database struct {
	Dbs *gorm.DB
}

var db *Database

func InitDb() error {
	Log.Info("Init db..............")

	platDb, err := initPlatDb(&conf.Conf.DbConf)
	if err != nil {
		Log.Errorf("Init draw db is err: %v", err)
		return fmt.Errorf("Init draw db is err: %v", err)
	}

	Log.Info("Init platform Db..............")

	db = &Database{
		Dbs: platDb,
	}

	return nil
}

// 初始化数据库
func initPlatDb(dbConf *conf.DbConf) (*gorm.DB, error) {
	if dbConf == nil {
		Log.Error("Db config is nil")
		return nil, fmt.Errorf("Db config is nil")
	}

	dbConStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.DbUser, dbConf.DbPassword, dbConf.DbHost, dbConf.DbPort, dbConf.DbName,
	)

	tmpDb, err := gorm.Open("mysql", dbConStr)
	if err != nil {
		tmpStr := fmt.Sprintf("Connet the db err, dbHost: %v, dbPort: %v, err: %v", dbConf.DbHost, dbConf.DbPort, err)
		Log.Errorf(tmpStr)
		return nil, fmt.Errorf(tmpStr)
	}

	tmpDb.DB().SetMaxOpenConns(dbConf.DbMaxConnect)
	tmpDb.DB().SetMaxIdleConns(dbConf.DbIdleConnect)

	if err = tmpDb.DB().Ping(); err != nil {
		tmpStr := fmt.Sprintf("Ping the db, dbHost: %v, dbPort: %v, err: %v", dbConf.DbHost, dbConf.DbPort, err)
		Log.Error(tmpStr)
		return nil, fmt.Errorf(tmpStr)
	}

	tmpDb.LogMode(dbConf.DbLogEnable)
	tmpDb.SingularTable(true)

	tables := []interface{}{
		&LeisureUsers{},
		&LeisureUserToken{},
	}

	tmpDb = tmpDb.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(tables...)
	for _, v := range tables {
		if !tmpDb.HasTable(v) {
			Log.Errorf("build table %v failed", v)
			return nil, fmt.Errorf("build table %v failed", v)
		}
	}

	return tmpDb, nil
}

func Close() error {
	db.Dbs.Close()
	return nil
}

// 更新 Dbs 数据库
func UpdateDbs(row interface{}) error {
	var err error

	for i := 0; i < 5; i++ {
		err = db.Dbs.Save(row).Error
		if err != nil {
			if i < 4 {
				time.Sleep(time.Duration((i + 1)) * 2 * time.Second)
			}
			continue
		} else {
			break
		}
	}

	if err != nil {
		Log.Error("Update dbs err: ", err)
		return err
	}

	return nil
}

// 插入 dbs 数据库
func InsertDbs(row interface{}) error {
	var err error

	for i := 0; i < 5; i++ {
		err = db.Dbs.Create(row).Error
		if err != nil {
			if i < 4 {
				time.Sleep(time.Duration((i + 1)) * 2 * time.Second)
			}
			continue
		} else {
			break
		}
	}

	if err != nil {
		Log.Error("insert dbs err: ", err)
		return err
	}

	return nil
}
