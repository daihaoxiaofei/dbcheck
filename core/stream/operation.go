package stream

import (
	"fmt"
	"strconv"
	"time"
)

type StreamStruct struct {
	argss []interface{}
	Data  []map[string]interface{}
}

var Strea *StreamStruct

// map函数
func (stre StreamStruct) Filter(f func(interface{}) bool) StreamStruct {
	res := make([]map[string]interface{}, 0, len(stre.Data))
	for _, item := range stre.Data {
		if f(item) {
			res = append(res, item)
		}
	}
	return StreamStruct{Data: res}
}

func (stre *StreamStruct) GetTimeDayArr(startTime, endTime string) (int64, error) {
	// 转成时间戳
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	startUnix, err := time.ParseInLocation(timeLayout, startTime, loc)
	if err != nil {
		return 0, err
	}
	endUnix, err := time.ParseInLocation(timeLayout, endTime, loc)
	if err != nil {
		return 0, err
	}
	// 求相差天数
	dateDay := (endUnix.Unix() - startUnix.Unix()) / 86400
	return dateDay, nil
}

func (stre *StreamStruct) GetTimeHoursArr(startTime, endTime string) (int64, error) {
	// 转成时间戳
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	startUnix, err := time.ParseInLocation(timeLayout, startTime, loc)
	if err != nil {
		return 0, err
	}
	endUnix, err := time.ParseInLocation(timeLayout, endTime, loc)
	if err != nil {
		return 0, err
	}
	// 求相差天数
	dateHourse := (endUnix.Unix() - startUnix.Unix()) / 3600
	return dateHourse, nil
}
func (stre *StreamStruct) GetTimeMinutesArr(startTime, endTime string) (int64, error) {
	// 转成时间戳
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	startUnix, err := time.ParseInLocation(timeLayout, startTime, loc)
	if err != nil {
		return 0, err
	}
	endUnix, err := time.ParseInLocation(timeLayout, endTime, loc)
	if err != nil {
		return 0, err
	}
	// 求相差天数
	dateMinutes := (endUnix.Unix() - startUnix.Unix()) / 600
	return dateMinutes, nil
}

func (stre *StreamStruct) GetTimeSecondsArr(startTime, endTime string) (int64, error) {
	// 转成时间戳
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	startUnix, err := time.ParseInLocation(timeLayout, startTime, loc)
	if err != nil {
		return 0, err
	}
	endUnix, err := time.ParseInLocation(timeLayout, endTime, loc)
	if err != nil {
		return 0, err
	}
	// 求相差天数
	dateSeconds := endUnix.Unix() - startUnix.Unix()
	return dateSeconds, nil
}

func (stre *StreamStruct) Add(args ...interface{}) (int, error) {
	var intArry = make([]int, len(args))
	for i, v := range args {
		d, err := strconv.Atoi(fmt.Sprintf("%v", v))
		if err != nil {
			return 0, err
		}
		intArry[i] = d
	}
	c := intArry[0] + intArry[1]
	return c, nil
}

// Percentage 返回值计算的百分比%  a *100 / b
func (stre *StreamStruct) Percentage(args ...interface{}) (int, error) {
	var intArry = make([]int, len(args))
	var c int
	for i, v := range args {
		d, err := strconv.Atoi(fmt.Sprintf("%v", v))
		if err != nil {
			return 0, err
		}
		intArry[i] = d
	}
	if intArry[0] == 0 {
		c = 0
	} else {
		c = intArry[0] * 100 / intArry[1]
	}
	return c, nil
}
