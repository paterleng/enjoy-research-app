package controller

import (
	"github.com/bwmarrin/snowflake"
	"github.com/sony/sonyflake"
	"time"
)

var timeTemplates = []string{
	"2006-01-02 15:04:05", //常规类型
	"2006/01/02 15:04:05",
	"2006-01-02",
	"2006/01/02",
	"15:04:05",
}

/* 时间格式字符串转换 */
func TimeStringToGoTime(tm string) time.Time {
	for i := range timeTemplates {
		t, err := time.ParseInLocation(timeTemplates[i], tm, time.Local)
		if nil == err && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}

//雪花算法生成随机id
var node *snowflake.Node

//初始化一个node
func Init(startTime string, machineID int64) (err error) {
	//自定义开始时间
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}

var (
	sonyFlake     *sonyflake.Sonyflake // 实例
	sonyMachineID uint16
	//机器ID
)

func getMachineID() (uint16, error) { //返回全局定义的机器ID
	return sonyMachineID, nil
}
