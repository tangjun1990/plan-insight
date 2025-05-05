package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	sysTimeLocation, _ = time.LoadLocation("Asia/Shanghai")
)

const TimeLayout = "15:04"

// Hour 获取时间的 24小时小时数
func Hour(t time.Time) uint8 {
	h := t.In(sysTimeLocation).Format("15")
	res, _ := strconv.Atoi(h)
	return uint8(res)
}

func HourString(t time.Time) string {
	return t.In(sysTimeLocation).Format("15")
}

func MinuteString(t time.Time) string {
	return t.In(sysTimeLocation).Format("04")
}

func ParseInLocation(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, sysTimeLocation)
}

// NowMonth 当前月份
func NowMonth() string {
	return nowFormat("200601")
}

// NowDate 当前日期
func NowDate() string {
	return nowFormat("20060102")
}

// NowMinute 当前分钟数
func NowMinute() string {
	return nowFormat("200601021504")
}

// nowFormat 当前时间吗
func nowFormat(layout string) string {
	return time.Now().In(sysTimeLocation).Format(layout)
}

func InTimeSpan(start, end, check time.Time) bool {
	if check.After(end) {
		return false
	}
	if end.After(start) {
		return check.After(start)
	}
	return check.Before(start)
}

// WaitToExec 距离下一次可执行段开始时间的等待时间
// times 可执行时间段 ["08:00-09:22", "10:00-11:00"]
// times 格式如上,时间依次递增
// return 第一个参数 0 表示在执行时间段内, 否则表示到下一次执行的时间
// return 第二个参数 表示下次执行是否跨天
func WaitToExec(times []string) (time.Duration, bool, error) {
	if len(times) == 0 {
		return 0, false, fmt.Errorf("时间段格式异常 数据为空")
	}

	prevTime := time.Time{}
	now := time.Now()
	date := now.Format("2006-01-02")
	var firstStart time.Time
	for i, t := range times {
		tArr := strings.SplitN(t, "-", 2)
		if len(tArr) != 2 {
			return 0, false, fmt.Errorf("无效的时间段格式: %s", t)
		}

		// 开始时间
		start, err := time.ParseInLocation("2006-01-02 15:04", date+" "+tArr[0], sysTimeLocation)
		if err != nil {
			return 0, false, fmt.Errorf("无效的开始时间格式: %s", tArr[0])
		}
		if prevTime.After(start) {
			return 0, false, fmt.Errorf("时间段必须依次递增")
		}
		if i == 0 {
			firstStart = start
		}

		// 结束时间
		end, err := time.ParseInLocation("2006-01-02 15:04", date+" "+tArr[1], sysTimeLocation)
		if err != nil {
			return 0, false, fmt.Errorf("无效的结束时间格式: %s", tArr[1])
		}
		if start.After(end) {
			return 0, false, fmt.Errorf("时间段顺序错误: %s", t)
		}

		if now.Before(start) {
			return start.Sub(now), false, nil
		}
		if now.Before(end) {
			return 0, false, nil
		}
		prevTime = end
	}

	return firstStart.AddDate(0, 0, 1).Sub(now), true, nil
}

// NowBetween 当前时间是否在两个时间之间
func NowBetween(startHour, startMin, endHour, endMin int) bool {
	now := time.Now().In(sysTimeLocation)
	y := now.Year()
	m := now.Month()
	d := now.Day()
	start := time.Date(y, m, d, startHour, startMin, 0, 0, sysTimeLocation)
	end := time.Date(y, m, d, endHour, endMin, 0, 0, sysTimeLocation)

	return now.After(start) && now.Before(end)
}
