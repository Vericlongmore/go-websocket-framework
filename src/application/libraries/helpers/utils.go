package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"application/libraries/toml"

	"github.com/coreos/etcd/clientv3"
	//"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	//"gopkg.in/mgo.v2"
)

func AnyTypeInt(sint interface{}) (ret int) {
	switch sint.(type) {
	case string:
		ret, _ = strconv.Atoi(sint.(string))
		break
	case int:
		ret = sint.(int)
		break
	case int64:
		ret = int(sint.(int64))
		break
	}
	return
}

func LoadFile(filename string) (filedata string, err error) {
	fileDataByte, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Read failed", err)
		return
	}

	filedata = string(fileDataByte)
	return
}

func MysqlDailName(mysqlconfig toml.DBConfig, dbName string) *xorm.Engine {
	db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysqlconfig.User, mysqlconfig.Password, mysqlconfig.Host, mysqlconfig.Port, dbName))
	if err != nil {
		panic(err)
	}
	return db
}

func EtcdDail(etcdconfig toml.DBConfig) *clientv3.Client {
	cfg := clientv3.Config{
		Endpoints: []string{etcdconfig.Host},
		//Transport:   client.DefaultTransport,
		DialTimeout: 5 * time.Second,
	}
	c, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func Pack(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	return string(data), err
}

func UnPack(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func snakeCasedName(name string) string {
	newstr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

//计算64位int中1的位数
func BitCount(i uint64) int {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	return int((((i + (i >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56)
}



var innerId int32 = 49999999

func GenInnerIncrId() int {
	num := atomic.AddInt32(&innerId, 1)
	return int(num)
}


func GetServerIp() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}

	return ""
}
