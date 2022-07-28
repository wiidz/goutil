package mathHelper

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/sliceHelper"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// GetRandomFloat64 获取随机浮点
func GetRandomFloat64() (num float64) {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()
}

// GetRandomInt 获取范围内的int随机数
func GetRandomInt(min, max int) (num int) {
	rand.Seed(time.Now().UnixNano())
	for {
		tmp := rand.Intn(max)
		if tmp >= min {
			num = tmp
			break
		}
	}
	return
}

// GetBezierPoints 获取贝塞尔曲线上的所有规定数量的点
func GetBezierPoints(dots []map[string]float64, amount int) []map[string]float64 {
	points := make([]map[string]float64, 0)
	for i := 0; i < amount; i++ {
		point := multiPointBezier(dots, float64(i)/float64(amount))
		points = append(points, point)
	}
	return points
}

// multiPointBezier 获取贝塞尔曲线上的多个点
func multiPointBezier(dots []map[string]float64, t float64) map[string]float64 {
	len := float64(len(dots))
	x := float64(0)
	y := float64(0)
	erxiangshi := func(start float64, end float64) float64 {
		cs := float64(1)
		bcs := float64(1)
		for end > 0 {
			cs *= start
			bcs *= end
			start--
			end--
		}
		return cs / bcs
	}
	i := 0
	for float64(i) < len {
		point := dots[i]
		temp := math.Pow(float64(1)-t, len-float64(1)-float64(i)) * math.Pow(t, float64(i)) * erxiangshi(len-float64(1), float64(i))
		x += point["x"] * temp
		y += point["y"] * temp
		i++
	}
	return map[string]float64{"x": x, "y": y}
}

// MatrixTransform 矩阵变换
func MatrixTransform(data []float64, matrix [][]float64) []float64 {
	res := make([]float64, 0)
	for _, row := range matrix {
		sum := float64(0)
		i := 0
		for _, num := range data {
			sum += num * row[i]
			i++
		}
		res = append(res, sum)
	}
	return res
}

// GetNearestIntXDots 获取一组点中最接近整数x的点
func GetNearestIntXDots(dots []map[string]float64) map[float64]map[string]float64 {
	//fmt.Println("【dots】",dots)
	intDots := make(map[float64]map[string]float64, 0)
	flag := false
	for k, v := range dots {
		dots[k]["x_int"] = math.Round(v["x"])
		dots[k]["diff"] = math.Abs(math.Round(v["x"]) - v["x"])
		//【1】判断是否存在
		flag = false
		for kk, _ := range intDots {
			if kk == dots[int64(k)]["x_int"] {
				flag = true
			}
		}
		//【2】替换最小值
		if flag {
			if dots[int64(k)]["diff"] > v["diff"] {
				intDots[dots[k]["x_int"]] = dots[k]
			}
		} else {
			intDots[dots[k]["x_int"]] = dots[k]
		}
	}
	//arr := make([]map[string]float64,0)
	//for _,v := range intDots{
	//	arr=append(arr,v)
	//}
	return intDots
}

// GetInsideDots 获取两组点构成的线，包围区域中间的点
func GetInsideDots(dots1 map[float64]map[string]float64, dots2 map[float64]map[string]float64) []map[string]float64 {
	insideDots := make([]map[string]float64, 0)

	intArr1 := getKeys(dots1)
	intArr2 := getKeys(dots2)

	intArr := sliceHelper.Intersect(intArr1, intArr2)

	yArr := make([]float64, 0)
	for _, int := range intArr {
		yArr = GetIntergers(dots1[int.(float64)]["y"], dots2[int.(float64)]["y"], float64(1), float64(10))
		for _, y := range yArr {
			insideDots = append(insideDots, map[string]float64{"x": int.(float64), "y": y})
		}
	}
	return insideDots
}

func getKeys(imap map[float64]map[string]float64) []interface{} {
	var tmp []interface{}
	if len(imap) > 0 {
		for k, _ := range imap {
			tmp = append(tmp, k)
		}
	}
	return tmp
}

// GetIntergers 获取两个数之间指定数量的等分点
func GetIntergers(border1 float64, border2 float64, narrow_range float64, amount float64) []float64 {
	//【1】计算边界
	max := math.Max(border1, border2)
	min := math.Min(border1, border2)

	//【2】修正收缩范围 0.8 0.7 ？
	if narrow_range != 1 {
		diff := max - min
		max = max - diff*(float64(1)-narrow_range/float64(2))
		min = min + diff*(float64(1)-narrow_range/float64(2))
	}
	diff := max - min

	//【3】修正步长
	step := float64(1)
	for diff > float64(amount) {
		len := diff / step
		if len < amount {
			step--
			break
		}
		step++
	}
	return sliceHelper.GetRange(math.Ceil(min), math.Floor(max), step)
}

// Keep 保留几位小数
func Keep(number float64, amount int) float64 {
	newNumber, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(amount)+"f", number), 64)
	return newNumber
}
