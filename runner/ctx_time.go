package runner

import (
	"errors"
	"fmt"
	"time"
)

// ContextTime 时间处理聚合句柄
// 专门处理框架中各种时间相关的业务场景，统一时间戳存储和转换逻辑
type ContextTime struct {
	ctx *Context
}

// ===== 基准时间常量 =====

const (
	// BaseDate1970 基准日期：1970年1月1日
	// 用于"仅时间"场景，确保时间排序和比较的准确性
	BaseDate1970 = "1970-01-01"

	// BaseTimestamp1970 基准时间戳：1970年1月1日 00:00:00 的毫秒时间戳
	BaseTimestamp1970 = int64(0)

	// MillisecondsPerSecond 每秒毫秒数
	MillisecondsPerSecond = int64(1000)

	// MillisecondsPerMinute 每分钟毫秒数
	MillisecondsPerMinute = int64(60 * 1000)

	// MillisecondsPerHour 每小时毫秒数
	MillisecondsPerHour = int64(60 * 60 * 1000)

	// MillisecondsPerDay 每天毫秒数
	MillisecondsPerDay = int64(24 * 60 * 60 * 1000)
)

// ===== 仅时间处理方法 (kind:time) =====

// TimeOnlyToMillis 将时分秒转换为基于1970-01-01的毫秒时间戳
// 用于 widget:"kind:time" 的时间选择器
// 例如：9:30:00 -> 34200000 (毫秒)
func (ct *ContextTime) TimeOnlyToMillis(hour, minute, second int) int64 {
	return int64(hour)*MillisecondsPerHour +
		int64(minute)*MillisecondsPerMinute +
		int64(second)*MillisecondsPerSecond
}

// TimeOnlyFromString 从时间字符串解析为毫秒时间戳
// 支持格式：HH:mm、HH:mm:ss
// 例如："09:30" -> 34200000, "09:30:15" -> 34215000
func (ct *ContextTime) TimeOnlyFromString(timeStr string) (int64, error) {
	// 尝试解析 HH:mm:ss
	if t, err := time.Parse("15:04:05", timeStr); err == nil {
		return ct.TimeOnlyToMillis(t.Hour(), t.Minute(), t.Second()), nil
	}

	// 尝试解析 HH:mm
	if t, err := time.Parse("15:04", timeStr); err == nil {
		return ct.TimeOnlyToMillis(t.Hour(), t.Minute(), 0), nil
	}

	return 0, errors.New("invalid time format")
}

// MillisToTimeOnly 将毫秒时间戳转换为时分秒
// 返回：hour, minute, second
func (ct *ContextTime) MillisToTimeOnly(millis int64) (int, int, int) {
	hour := int(millis / MillisecondsPerHour)
	remaining := millis % MillisecondsPerHour

	minute := int(remaining / MillisecondsPerMinute)
	remaining = remaining % MillisecondsPerMinute

	second := int(remaining / MillisecondsPerSecond)

	return hour, minute, second
}

// MillisToTimeString 将毫秒时间戳转换为时间字符串
// format: "HH:mm" 或 "HH:mm:ss"
func (ct *ContextTime) MillisToTimeString(millis int64, format string) string {
	hour, minute, second := ct.MillisToTimeOnly(millis)

	switch format {
	case "HH:mm":
		return time.Date(1970, 1, 1, hour, minute, 0, 0, time.UTC).Format("15:04")
	case "HH:mm:ss":
		return time.Date(1970, 1, 1, hour, minute, second, 0, time.UTC).Format("15:04:05")
	default:
		return time.Date(1970, 1, 1, hour, minute, 0, 0, time.UTC).Format("15:04")
	}
}

// ===== 当前时间方法 =====

// NowMillis 获取当前时间的毫秒时间戳
// 框架标准：所有时间存储都使用毫秒时间戳
func (ct *ContextTime) NowMillis() int64 {
	return time.Now().UnixMilli()
}

// TodayTimeOnlyMillis 获取今天指定时间的毫秒时间戳
// 例如：今天的9:00 AM -> 今天日期 + 9:00的时间戳
func (ct *ContextTime) TodayTimeOnlyMillis(hour, minute, second int) int64 {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())
	return today.UnixMilli()
}

// TodayTimeOnlyFromString 从时间字符串获取今天该时间的毫秒时间戳
func (ct *ContextTime) TodayTimeOnlyFromString(timeStr string) (int64, error) {
	// 先解析为仅时间的毫秒数
	timeOnlyMillis, err := ct.TimeOnlyFromString(timeStr)
	if err != nil {
		return 0, err
	}

	// 转换为今天的具体时间
	hour, minute, second := ct.MillisToTimeOnly(timeOnlyMillis)
	return ct.TodayTimeOnlyMillis(hour, minute, second), nil
}

// ===== 时间计算方法 =====

// MillisToDays 毫秒转天数
func (ct *ContextTime) MillisToDays(millis int64) int64 {
	return millis / MillisecondsPerDay
}

// MillisToHours 毫秒转小时数
func (ct *ContextTime) MillisToHours(millis int64) int64 {
	return millis / MillisecondsPerHour
}

// MillisToMinutes 毫秒转分钟数
func (ct *ContextTime) MillisToMinutes(millis int64) int64 {
	return millis / MillisecondsPerMinute
}

// DaysToMillis 天数转毫秒
func (ct *ContextTime) DaysToMillis(days int64) int64 {
	return days * MillisecondsPerDay
}

// HoursToMillis 小时转毫秒
func (ct *ContextTime) HoursToMillis(hours int64) int64 {
	return hours * MillisecondsPerHour
}

// MinutesToMillis 分钟转毫秒
func (ct *ContextTime) MinutesToMillis(minutes int64) int64 {
	return minutes * MillisecondsPerMinute
}

// ===== 时间比较和验证方法 =====

// IsTimeOnlyInRange 检查仅时间是否在指定范围内
// 用于验证上班时间、营业时间等业务逻辑
func (ct *ContextTime) IsTimeOnlyInRange(checkTime, startTime, endTime int64) bool {
	// 将时间戳转换为一天内的毫秒数进行比较
	checkMillisInDay := checkTime % MillisecondsPerDay
	startMillisInDay := startTime % MillisecondsPerDay
	endMillisInDay := endTime % MillisecondsPerDay

	if startMillisInDay <= endMillisInDay {
		// 正常情况：开始时间 <= 结束时间（如 9:00 - 18:00）
		return checkMillisInDay >= startMillisInDay && checkMillisInDay <= endMillisInDay
	} else {
		// 跨天情况：开始时间 > 结束时间（如 22:00 - 6:00）
		return checkMillisInDay >= startMillisInDay || checkMillisInDay <= endMillisInDay
	}
}

// CalculateTimeDifference 计算两个时间戳的差值（分钟）
// 返回正数表示 time1 晚于 time2，负数表示 time1 早于 time2
func (ct *ContextTime) CalculateTimeDifference(time1, time2 int64) int64 {
	diffMillis := time1 - time2
	return diffMillis / MillisecondsPerMinute
}

// ===== 格式化和显示方法 =====

// FormatMillisToDateTime 格式化毫秒时间戳为日期时间字符串
func (ct *ContextTime) FormatMillisToDateTime(millis int64, format string) string {
	t := time.UnixMilli(millis)

	switch format {
	case "yyyy-MM-dd":
		return t.Format("2006-01-02")
	case "yyyy-MM-dd HH:mm":
		return t.Format("2006-01-02 15:04")
	case "yyyy-MM-dd HH:mm:ss":
		return t.Format("2006-01-02 15:04:05")
	case "HH:mm":
		return t.Format("15:04")
	case "HH:mm:ss":
		return t.Format("15:04:05")
	default:
		return t.Format("2006-01-02 15:04:05")
	}
}

// FormatDuration 格式化时长（毫秒）为可读字符串
// 例如：90061000 -> "1小时30分钟1秒"
func (ct *ContextTime) FormatDuration(millis int64) string {
	if millis < 0 {
		millis = -millis
	}

	days := millis / MillisecondsPerDay
	hours := (millis % MillisecondsPerDay) / MillisecondsPerHour
	minutes := (millis % MillisecondsPerHour) / MillisecondsPerMinute
	seconds := (millis % MillisecondsPerMinute) / MillisecondsPerSecond

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%d分钟", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%d秒", seconds)
	}

	if result == "" {
		return "0秒"
	}
	return result
}

// ===== 业务场景快捷方法 =====

// CalculateLateMinutes 计算迟到分钟数
// 如果没有迟到返回0，迟到返回正数分钟
func (ct *ContextTime) CalculateLateMinutes(checkInTime, workStartTime int64) int64 {
	diff := ct.CalculateTimeDifference(checkInTime, workStartTime)
	if diff > 0 {
		return diff
	}
	return 0
}

// CalculateEarlyLeaveMinutes 计算早退分钟数
// 如果没有早退返回0，早退返回正数分钟
func (ct *ContextTime) CalculateEarlyLeaveMinutes(checkOutTime, workEndTime int64) int64 {
	diff := ct.CalculateTimeDifference(workEndTime, checkOutTime)
	if diff > 0 {
		return diff
	}
	return 0
}
