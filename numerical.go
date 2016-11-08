package main

import(
    "regexp"
    "strings"
    "strconv"
    "fmt"
    "errors"
    "log"
)

type NumericalExtract struct {
}

func (t *NumericalExtract) Number (text string) (float64, error) {
    reg := regexp.MustCompile(`[\d|\\.]+`)
    result := reg.FindAllString(text, -1)
    if len(result) == 0 {
        reg = regexp.MustCompile(`([一|二|三|四|五|六|七|八|九|十|零|百|千|万|亿]+)`)
        result = reg.FindAllString(text, -1)
        if len(result) > 0 {
            var sum = "0"
            var unit = 0
            var unitStr = ""
            var isFirst = false
            for i, c := range result[0] {
                w := fmt.Sprintf("%c", c)
                switch w {
                    case "十":
                        unit = 10
                        unitStr = "十"
                        if i == 0 {
                            isFirst = true
                        }
                        break
                    case "百":
                        unit = 100
                        unitStr = "百"
                        break
                    case "千":
                        unit = 1000
                        unitStr = "千"
                        break
                    case "万":
                        unit = 1000
                        unitStr = "万"
                        break
                    case "亿":
                        unit = 100000000
                        unitStr = "亿"
                        break
                    case "零":
                        sum = sum + "0"
                        break
                    case "一":
                        sum = sum + "1"
                        break
                    case "二":
                        sum = sum + "2"
                        break
                    case "三":
                        sum = sum + "3"
                        break
                    case "四":
                        sum = sum + "4"
                        break
                    case "五":
                        sum = sum + "5"
                        break
                    case "六":
                        sum = sum + "6"
                        break
                    case "七":
                        sum = sum + "7"
                        break
                    case "八":
                        sum = sum + "8"
                        break
                    case "九":
                        sum = sum + "9"
                        break
                }
                if  w != unitStr{
                    unit = unit / 10
                }
            }
            f, _ := strconv.ParseFloat(sum, 64)
            if isFirst {
                return f + 10, nil
            }
            return f * float64(unit), nil
        }
    } else {
        f, _ :=strconv.ParseFloat(result[0], 64)
        log.Println(f)
        return f, nil
    }
    return 0.0, errors.New("no found number")
}

func (t *NumericalExtract) checkTimeUint(text string) (string) {
    for _, u := range []string{"天", "时", "分", "秒"} {
        if strings.Contains(text, u) {
            switch u {
                case "天":
                    return "day"
                case "时":
                    return "hour"
                case "分":
                    return "minue"
                case "秒":
                    return "second"
            }
        }
    }
    return ""
}

func (t *NumericalExtract) Time(text string) (float64, string, error) {
    reg := regexp.MustCompile(`([\d|半|\\.])+?(天半|天|小时|小时半|个小时|分钟|分钟半|秒|秒半)`)
    result := reg.FindAllString(text, -1)
    var half = 0.0
    if strings.Contains(text, "半") {
        half = 0.5
    }
    if len(result) == 0 {
        reg := regexp.MustCompile(`([一|二|三|四|五|六|七|八|九|十|零|半])+?(天半|天|小时|小时半|个小时|分钟|分钟半|秒|秒半)`)
        result = reg.FindAllString(text, -1)
        if len(result) != 0 {
            f, _ := t.Number(result[0]) 
            return f + half, t.checkTimeUint(result[0]), nil
        } else {
            return 0.0, "", errors.New("time is empty")
        }
    } else {
        f, _ := t.Number(result[0]) 
        return f + half, t.checkTimeUint(result[0]), nil
    }
}

func (t *NumericalExtract) checkLengthUint(text string) (string) {
    for _, u := range []string{"米", "厘米", "分米", "毫米", "丈", "尺", "寸", "里", "公里"} {
        if strings.Contains(text, u) {
            switch u {
                case "米":
                    return "m"
                case "千米":
                    return "km"
                case "厘米":
                    return "cm"
                case "毫米":
                    return "mm"
                case "丈":
                    return "z"
                case "尺":
                    return "c"
                case "寸":
                    return "c"
                case "里":
                    return "l"
                case "公里":
                    return "g"
            }
        }
    }
    return ""
}

func (t *NumericalExtract) Length(text string) (float64, string, error) {
    reg := regexp.MustCompile(`([\d|半|\\.])+?(米半|米|厘米半|厘米|千米|分米|毫米|丈|尺|寸|公里|里)`)
    result := reg.FindAllString(text, -1)
    var half = 0.0
    if strings.Contains(text, "半") {
        half = 0.5
    }
    if len(result) == 0 {
        reg := regexp.MustCompile(`([一|二|三|四|五|六|七|八|九|十|零|半])+?(米半|米|厘米半|厘米|千米|分米|毫米|丈|尺|寸|公里|里)`)
        result = reg.FindAllString(text, -1)
        if len(result) != 0 {
            f, _ := t.Number(result[0]) 
            return f + half, t.checkLengthUint(result[0]), nil
        } else {
            return 0.0, "", errors.New("time is empty")
        }
    } else {
        f, _ := t.Number(result[0]) 
        return f + half, t.checkLengthUint(result[0]), nil
    }
}
