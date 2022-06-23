package colorHelper

import (
	"errors"
	"fmt"
	"github.com/wiidz/goutil/helpers/imgHelper"
	"github.com/wiidz/goutil/helpers/mapHelper"
	"github.com/wiidz/goutil/helpers/mathHelper"
	"github.com/wiidz/goutil/helpers/osHelper"
	"github.com/wiidz/goutil/helpers/sliceHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"math"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ColorData struct {
	rgb_r     uint8
	rgb_g     uint8
	rgb_b     uint8
	color_hex string
	xyz_x     float64
	xyz_y     float64
	xyz_z     float64
	lab_l     float64
	lab_a     float64
	lab_b     float64
	diff      float64
	diff_true float64
	grayscale float64
}

type ColorDiffSlice []map[string]interface{}

func (s ColorDiffSlice) Len() int           { return len(s) }
func (s ColorDiffSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ColorDiffSlice) Less(i, j int) bool { return s[i]["diff"].(float64) < s[j]["diff"].(float64) }

type ColorDiffTrueSlice []map[string]interface{}

func (s ColorDiffTrueSlice) Len() int      { return len(s) }
func (s ColorDiffTrueSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ColorDiffTrueSlice) Less(i, j int) bool {
	return s[i]["diff_true"].(float64) < s[j]["diff_true"].(float64)
}

//16进制转rgb颜色
func Hex2rgb(str string) (r, g, b int64, distance float64) {

	if len(str) != 7 {

		return 0, 0, 0, 0
	}
	r1, _ := strconv.ParseInt(str[1:3], 16, 10)
	g1, _ := strconv.ParseInt(str[3:5], 16, 18)
	b1, _ := strconv.ParseInt(str[5:7], 16, 10)

	distance1 := math.Sqrt(float64(r1)*float64(r1) + float64(g1)*float64(g1) + float64(b1)*float64(b1))

	return r1, g1, b1, distance1
}

/**
* @func: GetAvgColor 获取平均色值
* @author： Wiidz
* @date：  2019-08-28
 */
func GetAvgColor(colorArr []map[string]interface{}, narrowAmount int, illum string, degree int) (map[string]interface{}, error) {
	//【1】收缩数组，确定数组长度
	colorArr = sliceHelper.NarrowSlice(colorArr, narrowAmount)
	length := len(colorArr)
	//【2】一次计算平均色
	avg_x := sliceHelper.Float64SliceSum(sliceHelper.GetFloat64FromMapSlice(colorArr, "diff")) / float64(length)
	//fmt.Println("avg_x",avg_x)
	//【3】计算偏移量
	flag := 0
	for k, v := range colorArr {
		if v["diff"].(float64) >= avg_x {
			flag = k
			break
		}
	}
	//fmt.Println("flag",flag)
	offset := flag - int(math.Ceil(float64(length)/float64(2)))
	//fmt.Println("offset",offset)
	//【4】计算修正值，假设我们最多砍掉一边10个数字，每偏移2位砍掉1个，偏移20的时候达到最大值，另一边砍掉/2的数量
	left, right := 0, 0
	if offset < 0 {
		left = int(math.Abs(float64(offset) / float64(2)))
		if left > 10 {
			left = 10
		}
	} else {
		right = int(math.Abs(float64(offset) / float64(2)))
		if right > 10 {
			right = 10
		}
	}
	//fmt.Println("left",left,right)
	//【5】处理colors数组，进行偏移补正
	length = length - left - right
	colorArr = colorArr[left:length]

	avgColor := map[string]interface{}{"lab_l": sliceHelper.Float64SliceSum(sliceHelper.GetFloat64FromMapSlice(colorArr, "lab_l")) / float64(length), "lab_a": sliceHelper.Float64SliceSum(sliceHelper.GetFloat64FromMapSlice(colorArr, "lab_a")) / float64(length), "lab_b": sliceHelper.Float64SliceSum(sliceHelper.GetFloat64FromMapSlice(colorArr, "lab_b")) / float64(length)}
	avgColor, _ = Lab2Rgb(avgColor, illum, degree)
	avgColor["color_hex"] = Rgb2Hex(avgColor)

	return avgColor, nil
}

/**
* @func: SortByDiff 根据颜色数组中的diff键值从小到大排序
* @author： Wiidz
* @date：  2019-08-28
 */
func SortByDiff(colorArr []map[string]interface{}, targetColor map[string]interface{}, illum string, degree int) []map[string]interface{} {
	var colorDiffSlice ColorDiffSlice
	for k, v := range colorArr {
		colorArr[k], _ = Rgb2Lab(v, illum, degree)
		colorArr[k]["diff"] = GetColorDiffE2000(targetColor, v)
		colorDiffSlice = append(colorDiffSlice, colorArr[k])
	}
	sort.Sort(colorDiffSlice)
	return colorDiffSlice
}

/**
* @func: GetPixelRgb 获取图片中像素点的颜色
* @author： Wiidz
* @date：  2019-08-28
 */
func GetPixelRgb(img_uri string, dots []map[string]float64) ([]map[string]interface{}, error) {
	tmp := strings.Split(path.Base(img_uri), ".")
	format := tmp[len(tmp)-1]
	rgbArr := make([]map[string]interface{}, 0)

	//打开图片需要先下载文件，保存目录在tmp下，以时间戳+随机1000+后缀名组成
	localUri := "/tmp/go_downloads/" + strconv.FormatInt(time.Now().Unix(), 10) + typeHelper.Int2Str(mathHelper.GetRandomInt(1, 1000)) + "." + format
	err := osHelper.DownloadFile(img_uri, localUri, func(length, downLen int64) {})
	defer os.Remove(localUri)

	if err != nil {
		fmt.Println(err)
		return rgbArr, err
	}
	//获取image
	m, err := imgHelper.OpenImageFile(localUri)
	if err != nil {
		fmt.Println(err)
		return rgbArr, err
	}
	//循环获取点的rgba
	for _, v := range dots {
		r, g, b, a := m.At(int(v["x"]), int(v["y"])).RGBA()
		r_uint8, g_uint8, b_uint8, a_uint8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		ta := map[string]interface{}{"rgb_r": float64(r_uint8), "rgb_g": float64(g_uint8), "rgb_b": float64(b_uint8), "rgb_a": float64(a_uint8)}
		rgbArr = append(rgbArr, ta)
		//bounds := m.Bounds()
		//fmt.Println(bounds.Max.X,bounds.Min.X,bounds.Max.Y,bounds.Min.Y)
	}
	return rgbArr, nil
}

/**
 * @func: GetRgb2XyzMatrix 获取rgb转xyz时的矩阵
 * @author： Wiidz
 * @date：  2019-08-28
 */
func GetRgb2XyzMatrix(illum string) (matrix [][]float64, err error) {
	switch illum {
	case "D65":
		return [][]float64{{0.4124564, 0.3575761, 0.1804375}, {0.2126729, 0.7151522, 0.0721750}, {0.0193339, 0.1191920, 0.9503041}}, nil
	default:
		return [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, errors.New("no any illum matrix")
	}
}

/**
 * @func: GetXyz2RgbMatrix 获取xyz转rgb时的矩阵
 * @author： Wiidz
 * @date：  2019-08-28
 */
func GetXyz2RgbMatrix(illum string) (matrix [][]float64, err error) {
	switch illum {
	case "D65":
		return [][]float64{{3.2404542, -1.5371385, -0.4985314}, {-0.9692660, 1.8760108, 0.0415560}, {0.0556434, -0.2040259, 1.0572252}}, nil
	default:
		return [][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, errors.New("no any illum matrix")
	}
}

/**
 * @func: GetCieParams 获取视场参数
 * @author： Wiidz
 * @date：  2019-08-28
 */
func GetCieParams(degree int, illum string) (map[string]float64, error) {
	switch degree {
	case 2:
		switch illum {
		case "A":
			return map[string]float64{"x": 109.85, "y": 100, "z": 35.58}, nil
		case "B":
			return map[string]float64{"x": 99.09, "y": 100, "z": 85.31}, nil
		case "C":
			return map[string]float64{"x": 98.07, "y": 100, "z": 118.22}, nil
		case "D50":
			return map[string]float64{"x": 96.42, "y": 100, "z": 82.52}, nil
		case "D55":
			return map[string]float64{"x": 95.68, "y": 100, "z": 92.15}, nil
		case "D65":
			return map[string]float64{"x": 95.04, "y": 100, "z": 108.89}, nil
		case "D75":
			return map[string]float64{"x": 94.97, "y": 100, "z": 122.64}, nil
		case "E":
			return map[string]float64{"x": 100, "y": 100, "z": 100}, nil
		case "F1":
			return map[string]float64{"x": 92.83, "y": 100, "z": 103.66}, nil
		case "F2":
			return map[string]float64{"x": 99.14, "y": 100, "z": 67.32}, nil
		case "F3":
			return map[string]float64{"x": 103.75, "y": 100, "z": 49.86}, nil
		case "F4":
			return map[string]float64{"x": 109.15, "y": 100, "z": 38.81}, nil
		case "F5":
			return map[string]float64{"x": 90.87, "y": 100, "z": 98.72}, nil
		case "F6":
			return map[string]float64{"x": 97.315, "y": 100, "z": 60.19}, nil
		case "F7":
			return map[string]float64{"x": 95.02, "y": 100, "z": 108.63}, nil
		case "F8":
			return map[string]float64{"x": 96.41, "y": 100, "z": 82.33}, nil
		case "F9":
			return map[string]float64{"x": 100.36, "y": 100, "z": 67.87}, nil
		case "F10":
			return map[string]float64{"x": 96.17, "y": 100, "z": 81.71}, nil
		case "F11":
			return map[string]float64{"x": 100.9, "y": 100, "z": 64.26}, nil
		case "F12":
			return map[string]float64{"x": 108.05, "y": 100, "z": 39.23}, nil
		default:
			return map[string]float64{}, nil
		}
	case 10:
		switch illum {
		case "A":
			return map[string]float64{"x": 111.14, "y": 100, "z": 35.2}, nil
		case "B":
			return map[string]float64{"x": 99.18, "y": 100, "z": 84.35}, nil
		case "C":
			return map[string]float64{"x": 97.29, "y": 100, "z": 116.14}, nil
		case "D50":
			return map[string]float64{"x": 96.72, "y": 100, "z": 81.43}, nil
		case "D55":
			return map[string]float64{"x": 95.8, "y": 100, "z": 90.93}, nil
		case "D65":
			return map[string]float64{"x": 94.81, "y": 100, "z": 107.31}, nil
		case "D75":
			return map[string]float64{"x": 94.42, "y": 100, "z": 120.64}, nil
		case "E":
			return map[string]float64{"x": 100, "y": 100, "z": 100}, nil
		case "F1":
			return map[string]float64{"x": 94.79, "y": 100, "z": 103.19}, nil
		case "F2":
			return map[string]float64{"x": 103.25, "y": 100, "z": 68.99}, nil
		case "F3":
			return map[string]float64{"x": 108.97, "y": 100, "z": 51.96}, nil
		case "F4":
			return map[string]float64{"x": 114.96, "y": 100, "z": 40.96}, nil
		case "F5":
			return map[string]float64{"x": 93.37, "y": 100, "z": 98.64}, nil
		case "F6":
			return map[string]float64{"x": 102.15, "y": 100, "z": 62.07}, nil
		case "F7":
			return map[string]float64{"x": 95.78, "y": 100, "z": 107.62}, nil
		case "F8":
			return map[string]float64{"x": 97.11, "y": 100, "z": 81.13}, nil
		case "F9":
			return map[string]float64{"x": 102.12, "y": 100, "z": 67.83}, nil
		case "F10":
			return map[string]float64{"x": 99, "y": 100, "z": 83.13}, nil
		case "F11":
			return map[string]float64{"x": 103.82, "y": 100, "z": 65.56}, nil
		case "F12":
			return map[string]float64{"x": 111.43, "y": 100, "z": 40.35}, nil
		default:
			return map[string]float64{}, nil
		}
	default:
		return map[string]float64{}, nil
	}
}

/**
 * @func: GetCieParams 获取视场参数
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Rgb2Xyz(arr map[string]interface{}, illum string) (map[string]interface{}, error) {
	gamma := func(data float64) float64 {
		if data > 0.04045 {
			return math.Pow((data+0.055)/1.055, 2.4)
		}
		return data / 12.92
	}
	//rgb := map[string]float64{"r":100*gamma(arr["rgb_r"]/float64(255)),"g":100*gamma(arr["rgb_r"]/float64(255)),"b":100*gamma(arr["rgb_r"]/float64(255))}
	rgb := []float64{100 * gamma(arr["rgb_r"].(float64)/float64(255)), 100 * gamma(arr["rgb_g"].(float64)/float64(255)), 100 * gamma(arr["rgb_b"].(float64)/float64(255))}
	matrix, err := GetRgb2XyzMatrix(illum)
	if err != nil {
		return map[string]interface{}{}, err
	}
	temp := mathHelper.MatrixTransform(rgb, matrix)

	arr["xyz_x"] = temp[0]
	arr["xyz_y"] = temp[1]
	arr["xyz_z"] = temp[2]

	return arr, nil
}

/**
 * @func: Xyz2Lab XYZ转Lab
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Xyz2Lab(arr map[string]interface{}, illum string, degree int) (map[string]interface{}, error) {
	//fmt.Println("【arr】",arr)
	handle := func(data float64) float64 {
		if data > float64(216)/float64(24389) {
			return math.Pow(data, float64(1)/float64(3))
		}
		return (float64(24389)/float64(27)*data + float64(16)) / float64(116)
	}
	cieParams, err := GetCieParams(degree, illum)
	//fmt.Println("【cieParams】",cieParams)
	if err != nil {
		return map[string]interface{}{}, err
	}
	xn, yn, zn := cieParams["x"], cieParams["y"], cieParams["z"]
	//fmt.Println("【xn】",xn,yn,zn)
	fx, fy, fz := handle(arr["xyz_x"].(float64)/xn), handle(arr["xyz_y"].(float64)/yn), handle(arr["xyz_z"].(float64)/zn)
	//fmt.Println("【fx】",fx,fy,fz)
	arr["lab_l"] = float64(116)*fy - float64(16)
	arr["lab_a"] = float64(500) * (fx - fy)
	arr["lab_b"] = float64(200) * (fy - fz)

	//fmt.Println("【Xyz2Lab】",arr)

	return arr, nil
}

/**
 * @func: Lab2Xyz Lab转XYZ
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Lab2Xyz(arr map[string]interface{}, illum string, degree int) (map[string]interface{}, error) {
	k := float64(24389) / float64(27)
	e := float64(216) / float64(24389)
	handle := func(data float64) float64 {
		temp := math.Pow(data, 3)
		if temp > e {
			return temp
		}
		return (data - float64(16)/float64(116)) / 7.787
	}

	//fx,fy,fz
	fy := (arr["lab_l"].(float64) + float64(16)) / float64(116)
	fx := arr["lab_a"].(float64)/float64(500) + fy
	fz := fy - arr["lab_b"].(float64)/float64(200)

	//xr,yr,zr
	xr := handle(fx)
	var yr float64
	if arr["lab_l"].(float64) > k*e {
		yr = math.Pow((arr["lab_l"].(float64)+float64(16))/116, 3)
	} else {
		yr = arr["lab_l"].(float64) / k
	}
	zr := handle(fz)

	//确定光源
	cieParams, err := GetCieParams(degree, illum)
	if err != nil {
		return map[string]interface{}{}, err
	}

	//计算最终的x,y,z
	arr["xyz_x"] = xr * cieParams["x"]
	arr["xyz_y"] = yr * cieParams["y"]
	arr["xyz_z"] = zr * cieParams["z"]
	return arr, nil
}

/**
 * @func: Xyz2Rgb XYZ转RGB
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Xyz2Rgb(arr map[string]interface{}, illum string, degree int) (map[string]interface{}, error) {
	companding := func(data float64) float64 {
		if data > 0.0031308 {
			return 1.055*math.Pow(data, float64(1)/float64(2.4)) - 0.055
		}
		return 12.92 * data
	}
	//xyz:=map[string]float64{"xyz_x":arr["xyz_x"]/float64(100),"xyz_y":arr["xyz_y"]/float64(100),"xyz_z":arr["xyz_z"]/float64(100)}
	xyz := []float64{arr["xyz_x"].(float64) / float64(100), arr["xyz_y"].(float64) / float64(100), arr["xyz_z"].(float64) / float64(100)}

	matrix, err := GetXyz2RgbMatrix(illum)
	if err != nil {
		return map[string]interface{}{}, err
	}
	temp := mathHelper.MatrixTransform(xyz, matrix)
	arr["rgb_r"] = float64(255) * companding(temp[0])
	arr["rgb_g"] = float64(255) * companding(temp[1])
	arr["rgb_b"] = float64(255) * companding(temp[2])
	return arr, nil
}

/**
 * @func: Rgb2Lab RGB转Lab
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Rgb2Lab(arr map[string]interface{}, illum string, degree int) (map[string]interface{}, error) {
	arr, err := Rgb2Xyz(arr, illum)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return Xyz2Lab(arr, illum, degree)
}

/**
 * @func: Lab2Rgb Lab转RGB
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Lab2Rgb(arr map[string]interface{}, illum string, degree int) (map[string]interface{}, error) {
	arr, err := Lab2Xyz(arr, illum, degree)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return Xyz2Rgb(arr, illum, degree)
}

/**
 * @func: Rgb2Hex RGB转Hex
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Rgb2Hex(arr map[string]interface{}) string {

	rgb := map[string]string{
		"r": strconv.FormatInt(int64(math.Round(arr["rgb_r"].(float64))), 16),
		"g": strconv.FormatInt(int64(math.Round(arr["rgb_g"].(float64))), 16),
		"b": strconv.FormatInt(int64(math.Round(arr["rgb_b"].(float64))), 16),
	}

	return "#" + rgb["r"] + rgb["g"] + rgb["b"]
}

/**
 * @func: Hex2Rgb Hex转RGB
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Hex2Rgb(hex string) map[string]interface{} {
	fmt.Println("【color】", hex)
	r, _ := strconv.ParseInt(hex[1:3], 16, 10)
	g, _ := strconv.ParseInt(hex[3:5], 16, 18)
	b, _ := strconv.ParseInt(hex[5:], 16, 10)
	return map[string]interface{}{"rgb_r": float64(r), "rgb_g": float64(g), "rgb_b": float64(b)}
}

/**
 * @func: Hex2Rgb Hex转RGB
 * @author： Wiidz
 * @date：  2019-08-28
 */
func Hex2RgbStr(hex string) map[string]string {
	if len(hex) != 7 {
		return map[string]string{"rgb_r": "0", "rgb_g": "0", "rgb_b": "0"}
	}
	r, _ := strconv.ParseInt(hex[1:3], 16, 10)
	g, _ := strconv.ParseInt(hex[3:5], 16, 18)
	b, _ := strconv.ParseInt(hex[5:], 16, 10)
	return map[string]string{"rgb_r": strconv.Itoa(int(r)), "rgb_g": strconv.Itoa(int(g)), "rgb_b": strconv.Itoa(int(b))}
}

/**
 * @func: GetColorDiffDECMC 使用DECMC(l:c)公式计算色差
 * @author： Wiidz
 * @date：  2019-08-28
 * @tips：
 * 在纺织中，l通常设为2，允许在⊿L*上有相对较大的容忍度，这也就是CMC(2:1)公式
 * ⊿L、⊿C、⊿H通过CIELAB1976空间计算得到
 * L*std、C*std、h*std为标准色的色度参数，就是我们要比较的那个颜色称为标准色
 */
func GetColorDiffDECMC(color1, color2 map[string]interface{}, l, c int) float64 {
	//⊿L、⊿A、⊿B、⊿C、⊿H、C1、C2
	deltaL := color1["lab_l"].(float64) - color2["lab_l"].(float64)
	deltaA := color1["lab_a"].(float64) - color2["lab_a"].(float64)
	deltaB := color1["lab_b"].(float64) - color2["lab_b"].(float64)
	C1 := math.Sqrt(math.Pow(color1["lab_a"].(float64), 2) + math.Pow(color1["lab_b"].(float64), 2))
	C2 := math.Sqrt(math.Pow(color2["lab_a"].(float64), 2) + math.Pow(color2["lab_b"].(float64), 2))
	//fmt.Println(math.Pow(color1["lab_a"],2))
	//fmt.Println(math.Pow(color1["lab_b"],2))
	//fmt.Println(math.Pow(color1["lab_a"],2)+math.Pow(color1["lab_b"],2))
	deltaC := C1 - C2
	deltaH := math.Sqrt(math.Pow(deltaA, 2) + math.Pow(deltaB, 2) - math.Pow(deltaC, 2)) //可能为负数
	//fmt.Println("【1】",deltaL,deltaA,deltaB,C1,C2,deltaC,deltaH)
	//fmt.Println(math.Pow(deltaA,2)+math.Pow(deltaB,2)-math.Pow(deltaC,2))
	//H、H1
	H := math.Atan(color1["lab_b"].(float64) / color1["lab_a"].(float64))
	var H1 float64
	if H >= 0 {
		H1 = H
	} else {
		H1 = H + float64(360)
	}
	//fmt.Println("【2】",H,H1)
	//f、T
	var T float64
	if float64(164) <= H1 && H1 <= float64(345) {
		T = 0.56 + math.Abs(0.2*math.Cos(H1+float64(168)))
	} else {
		T = 0.36 + math.Atan(H1+float64(35))
	}
	f := math.Sqrt(math.Pow(C1, 4) / (math.Pow(C1, 4) + 1900))
	//fmt.Println("【3】",f,T)

	//SC、SH、SL
	SC := (0.0638*C1)/(float64(1)+0.0131*C1) + 0.638 //彩度权重
	SH := SC * (f*T + float64(1) - f)                //色调权重
	var SL float64                                   //明度权重
	if color1["lab_l"].(float64) <= 16 {
		SL = 0.511
	} else {
		SL = (0.040975 * color1["lab_l"].(float64)) / (float64(1) + 0.01765*color1["lab_l"].(float64))
	}
	//fmt.Println("【4】",SC,SH,SL)

	//fmt.Println("【5】",math.Pow(deltaL/float64(l)/SL,2))
	//fmt.Println("【6】",math.Pow(deltaC/float64(c)/SC,2))
	//fmt.Println("【7】",math.Pow(deltaH*SH,2))

	return math.Sqrt(math.Pow(deltaL/float64(l)/SL, 2) + math.Pow(deltaC/float64(c)/SC, 2) + math.Pow(deltaH*SH, 2))

}

/**
 * @func: GetColorDiffE2000 使用Delta-E 2000公式计算色差
 * @author： Wiidz
 * @date：  2019-08-28
 */
func GetColorDiffE2000(color1, color2 map[string]interface{}) float64 {

	//fmt.Println("color1",color1)
	//fmt.Println("color2",color2)

	kl, kc, kh := float64(1), float64(1), float64(1)      //参考实验条件见P88
	rt, g, mean_cab := float64(0), float64(0), float64(0) //旋转函数rt,g表示CIELab 颜色空间a轴的调整因子,是彩度的函数,cab表示两个样品彩度的算术平均值
	pi := math.Pi
	mean_cab = (GetSaturation(color1["lab_a"].(float64), color1["lab_b"].(float64)) + GetSaturation(color2["lab_a"].(float64), color2["lab_b"].(float64))) / float64(2)
	mean_cab_pow7 := math.Pow(mean_cab, 7) //两彩度平均值的7次方
	g = 0.5 * (1 - math.Pow(mean_cab_pow7/(mean_cab_pow7+math.Pow(25, 7)), 0.5))

	ll1, aa1, bb1 := color1["lab_l"].(float64), color1["lab_a"].(float64)*(float64(1)+g), color1["lab_b"].(float64)
	ll2, aa2, bb2 := color2["lab_l"].(float64), color2["lab_a"].(float64)*(float64(1)+g), color2["lab_b"].(float64)

	// 两样本的彩度值
	cc1, cc2 := GetSaturation(aa1, bb1), GetSaturation(aa2, bb2)
	// 两样本的色调角
	hh1, hh2 := GetHueAngle(aa1, bb1), GetHueAngle(aa2, bb2)

	diff_ll := ll1 - ll2
	diff_cc := cc1 - cc2
	diff_hh := GetHueAngle(aa1, bb1) - GetHueAngle(aa2, bb2)
	diff_HH := 2 * math.Sin(pi*diff_hh/float64(360)) * math.Pow(cc1*cc2, 0.5)

	//-------第三步--------------
	//计算公式中的加权函数sl,sc,sh,t
	mean_ll := (ll1 + ll2) / 2
	mean_cc := (cc1 + cc2) / 2
	mean_hh := (hh1 + hh2) / 2

	sl := 1 + 0.015*math.Pow(mean_ll-50, 2)/math.Pow(20+math.Pow(mean_ll-50, 2), 0.5)
	sc := 1 + 0.045*mean_cc
	t := 1 - 0.17*math.Cos((mean_hh-30)*pi/180) + 0.24*math.Cos((2*mean_hh)*pi/180) + 0.32*math.Cos((3*mean_hh+6)*pi/180) - 0.2*math.Cos((4*mean_hh-63)*pi/180)
	sh := 1 + 0.015*mean_cc*t

	//------第四步--------
	//计算公式中的rt
	mean_cc_pow7 := math.Pow(mean_cc, 7)
	rc := float64(2) * math.Pow(mean_cc_pow7/(mean_cc_pow7+math.Pow(25, 7)), 0.5)
	diff_xita := float64(30) * math.Exp(-math.Pow((mean_hh-275)/float64(25), 2)) //△θ 以°为单位
	rt = -math.Sin((float64(2)*diff_xita)*pi/float64(180)) * rc

	//        var l_item, c_item, h_item
	l_item := diff_ll / (kl * sl)
	c_item := diff_cc / (kc * sc)
	h_item := diff_HH / (kh * sh)

	E00 := math.Pow(l_item*l_item+c_item*c_item+h_item*h_item+rt*c_item*h_item, 0.5)

	return E00

}

/**
* @func: GetSaturation 彩调计算
* @author： Wiidz
* @date：  2019-08-28
 */
func GetSaturation(a float64, b float64) float64 {
	return math.Pow(a*a+b*b, 0.5)
}

/**
* @func: GetHueAngle 色调角计算
* @author： Wiidz
* @date：  2019-08-28
 */
func GetHueAngle(a, b float64) float64 {
	if a == 0 {
		return 0
	}
	h, hab := float64(0), float64(0)

	h = float64(180) / math.Pi * math.Atan(b/a) //有正有负

	if a > 0 {
		if b > 0 {
			hab = h
		} else {
			hab = float64(360) + h
		}
	} else {
		hab = float64(180) + h
	}
	return hab
}

/**
* @func: GetGrayscale 计算灰度
* @author： Wiidz
* @date：  2019-08-28
 */
func GetGrayscale(rgb map[string]interface{}) interface{} {
	return (int(rgb["rgb_r"].(float64))*38 + int(rgb["rgb_g"].(float64))*75 + int(rgb["rgb_b"].(float64))*15) >> 7
}

/**
* @func: GetSimilarColors 获取数组中和目标颜色相近的颜色
* @author： Wiidz
* @date：  2019-08-28
 */
func GetSimilarColors(target map[string]interface{}, colorArr []map[string]interface{}, amount int, maxDiff float64) []map[string]interface{} {
	//【1】计算与白色的色差
	if !mapHelper.Exist(target, "white_diff") {
		target["white_diff"] = GetColorDiffE2000(map[string]interface{}{"lab_l": float64(100), "lab_a": float64(0), "lab_b": float64(0)}, target)
	}
	//【2】遍历所有颜色，计算色差
	var colorDiffTrueSlice ColorDiffTrueSlice
	for k, v := range colorArr {
		//舍弃diff大于1的color_diff
		if math.Abs(v["white_diff"].(float64)-target["white_diff"].(float64)) > maxDiff {
			continue
		}
		colorArr[k]["diff_true"] = GetColorDiffE2000(target, v)
		colorDiffTrueSlice = append(colorDiffTrueSlice, colorArr[k])
	}
	//【3】再次排序
	sort.Sort(colorDiffTrueSlice)
	return colorDiffTrueSlice
}
