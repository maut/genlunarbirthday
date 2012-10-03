package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const example = `"father's birthday": 1959.01.01
"mother's birthday": 1963.-4.01
`

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:\n    %s birthdays.txt\nThen a import.csv will be created\n",
		filepath.Base(os.Args[0]))
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "Example birthday.txt:\n"+example)
}

type date struct {
	year, month, day int
}

func newDate(year, month, day int) *date {
	return &date{year, month, day}
}

const (
	// 初始日，公历农历对应日期：
	// 公历 1901 年 1 月 1 日，对应农历 4598 年 11 月 11 日
	bgy = 1901 // base gregorian year
	bgm = 1    // base gregorian month
	bgd = 1    // base gregorian day
	bcy = 1900 // base chinese year
	bcm = 11   // base chinese month
	bcd = 11   // base chinese day
)

var gregMonth = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func gregorian2chinese(greg *date) (chin *date) {
	if !isValidGreg(greg) {
		return nil
	}
	// set starting gregorian date
	sy := bgy
	sm := bgm
	sd := bgd
	chin = new(date)
	chin.year = bcy
	chin.month = bcm
	chin.day = bcd
	// second matching date
	// 公历 2000 年 1 月 1 日，对应农历 4697 年 11 月 25 日
	if greg.year >= 2000 {
		sy = 2000
		sm = 1
		sd = 1
		chin.year = 1999
		chin.month = 11
		chin.day = 25
	}

	var diff int // difference in days to the start date
	// diff for years
	for i := sy; i < greg.year; i++ {
		if isGregLeap(i) {
			diff += 366
		} else {
			diff += 365
		}
	}
	// diff for months
	for i := sm; i < greg.month; i++ {
		diff += gregMonth[i-1]
		if i-1 == 2 && isGregLeap(greg.year) {
			diff++
		}
	}
	// diff for days
	diff += greg.day - sd
	// add base chinese day
	chin.day += diff

	// get the days for the base chinese month
	daysThisMon := daysInChinMonth(chin.month, chin.year)
	// get the next chinese month
	nextMon := nextChinMonth(chin.month, chin.year)
	for chin.day > daysThisMon {
		if abs(nextMon) < abs(chin.month) {
			chin.year++
		}
		chin.day -= daysThisMon
		chin.month = nextMon
		daysThisMon = daysInChinMonth(nextMon, chin.year)
		nextMon = nextChinMonth(chin.month, chin.year)
	}
	return
}

func isValidGreg(d *date) bool {
	if 1901 > d.year || d.year > 2100 {
		return false
	}
	if d.month > 12 {
		return false
	}
	if d.month == 2 && isGregLeap(d.year) {
		if d.day > 29 {
			return false
		}
	} else if d.day > gregMonth[d.month-1] {
		return false
	}
	return true
}

func chinese2gregorian(chin *date) (greg *date) {
	if !isValidChin(chin) {
		return nil
	}
	greg = new(date)
	greg.year = bgy
	greg.month = bgm
	greg.day = bgd

	var diff int
	diff = 49 // number of days from lunar 1900 11 11 to 1901 1 1
	// add days in years
	for y := 1901; y < chin.year; y++ {
		var days int
		for m := 1; ; {
			days += daysInChinMonth(m, y)
			m = nextChinMonth(m, y)
			if m == 1 {
				break
			}
		}
		diff += days
	}
	// add up days in previous months of this year
	for m := 1; ; {
		if chin.month == 1 {
			break
		} else if chin.month == nextChinMonth(m, chin.year) {
			diff += daysInChinMonth(m, chin.year)
			break
		}
		diff += daysInChinMonth(m, chin.year)
		m = nextChinMonth(m, chin.year)
	}
	// add up days diff in this month
	diff += chin.day

	// calculate gregorian date
	// deduct days in years
	for {
		if isGregLeap(greg.year) {
			if diff <= 366 {
				break
			}
			diff -= 366
		}
		if !isGregLeap(greg.year) {
			if diff <= 365 {
				break
			}
			diff -= 365
		}
		greg.year++
	}
	// deduct days in months
	for {
		if greg.month == 2 && isGregLeap(greg.year) {
			if diff <= 29 {
				break
			}
			diff -= 29
		} else {
			if diff <= gregMonth[greg.month-1] {
				break
			}
			diff -= gregMonth[greg.month-1]
		}
		greg.month++
	}
	greg.day = diff
	return
}

func isValidChin(d *date) bool {
	if 1901 > d.year || d.year > 2100 {
		return false
	}
	if d.month > 12 {
		return false
	}
	if d.month < 0 {
		if d.month != nextChinMonth(abs(d.month), d.year) {
			return false
		}
	} else {
		if d.day > daysInChinMonth(d.month, d.year) {
			return false
		}
	}
	return true
}

func isGregLeap(year int) (isleap bool) {
	if year%4 == 0 {
		isleap = true
	}
	if year%100 == 0 {
		isleap = false
	}
	if year%400 == 0 {
		isleap = true
	}
	return
}

// 农历月份大小压缩表，两个字节表示一年。两个字节共十六个二进制位数， 
// 前四个位数表示闰月月份，后十二个位数表示十二个农历月份的大小。
// e.g.: 0xad, 0x08: 0000 1000 1010 1101
// 0000 is leap month, 1000 1010 1101 is for number of days in each month, from right to left
var chMon = []byte{
	0x00, 0x04, 0xad, 0x08, 0x5a, 0x01, 0xd5, 0x54, 0xb4, 0x09, 0x64, 0x05, 0x59, 0x45,
	0x95, 0x0a, 0xa6, 0x04, 0x55, 0x24, 0xad, 0x08, 0x5a, 0x62, 0xda, 0x04, 0xb4, 0x05,
	0xb4, 0x55, 0x52, 0x0d, 0x94, 0x0a, 0x4a, 0x2a, 0x56, 0x02, 0x6d, 0x71, 0x6d, 0x01,
	0xda, 0x02, 0xd2, 0x52, 0xa9, 0x05, 0x49, 0x0d, 0x2a, 0x45, 0x2b, 0x09, 0x56, 0x01,
	0xb5, 0x20, 0x6d, 0x01, 0x59, 0x69, 0xd4, 0x0a, 0xa8, 0x05, 0xa9, 0x56, 0xa5, 0x04,
	0x2b, 0x09, 0x9e, 0x38, 0xb6, 0x08, 0xec, 0x74, 0x6c, 0x05, 0xd4, 0x0a, 0xe4, 0x6a,
	0x52, 0x05, 0x95, 0x0a, 0x5a, 0x42, 0x5b, 0x04, 0xb6, 0x04, 0xb4, 0x22, 0x6a, 0x05,
	0x52, 0x75, 0xc9, 0x0a, 0x52, 0x05, 0x35, 0x55, 0x4d, 0x0a, 0x5a, 0x02, 0x5d, 0x31,
	0xb5, 0x02, 0x6a, 0x8a, 0x68, 0x05, 0xa9, 0x0a, 0x8a, 0x6a, 0x2a, 0x05, 0x2d, 0x09,
	0xaa, 0x48, 0x5a, 0x01, 0xb5, 0x09, 0xb0, 0x39, 0x64, 0x05, 0x25, 0x75, 0x95, 0x0a,
	0x96, 0x04, 0x4d, 0x54, 0xad, 0x04, 0xda, 0x04, 0xd4, 0x44, 0xb4, 0x05, 0x54, 0x85,
	0x52, 0x0d, 0x92, 0x0a, 0x56, 0x6a, 0x56, 0x02, 0x6d, 0x02, 0x6a, 0x41, 0xda, 0x02,
	0xb2, 0xa1, 0xa9, 0x05, 0x49, 0x0d, 0x0a, 0x6d, 0x2a, 0x09, 0x56, 0x01, 0xad, 0x50,
	0x6d, 0x01, 0xd9, 0x02, 0xd1, 0x3a, 0xa8, 0x05, 0x29, 0x85, 0xa5, 0x0c, 0x2a, 0x09,
	0x96, 0x54, 0xb6, 0x08, 0x6c, 0x09, 0x64, 0x45, 0xd4, 0x0a, 0xa4, 0x05, 0x51, 0x25,
	0x95, 0x0a, 0x2a, 0x72, 0x5b, 0x04, 0xb6, 0x04, 0xac, 0x52, 0x6a, 0x05, 0xd2, 0x0a,
	0xa2, 0x4a, 0x4a, 0x05, 0x55, 0x94, 0x2d, 0x0a, 0x5a, 0x02, 0x75, 0x61, 0xb5, 0x02,
	0x6a, 0x03, 0x61, 0x45, 0xa9, 0x0a, 0x4a, 0x05, 0x25, 0x25, 0x2d, 0x09, 0x9a, 0x68,
	0xda, 0x08, 0xb4, 0x09, 0xa8, 0x59, 0x54, 0x03, 0xa5, 0x0a, 0x91, 0x3a, 0x96, 0x04,
	0xad, 0xb0, 0xad, 0x04, 0xda, 0x04, 0xf4, 0x62, 0xb4, 0x05, 0x54, 0x0b, 0x44, 0x5d,
	0x52, 0x0a, 0x95, 0x04, 0x55, 0x22, 0x6d, 0x02, 0x5a, 0x71, 0xda, 0x02, 0xaa, 0x05,
	0xb2, 0x55, 0x49, 0x0b, 0x4a, 0x0a, 0x2d, 0x39, 0x36, 0x01, 0x6d, 0x80, 0x6d, 0x01,
	0xd9, 0x02, 0xe9, 0x6a, 0xa8, 0x05, 0x29, 0x0b, 0x9a, 0x4c, 0xaa, 0x08, 0xb6, 0x08,
	0xb4, 0x38, 0x6c, 0x09, 0x54, 0x75, 0xd4, 0x0a, 0xa4, 0x05, 0x45, 0x55, 0x95, 0x0a,
	0x9a, 0x04, 0x55, 0x44, 0xb5, 0x04, 0x6a, 0x82, 0x6a, 0x05, 0xd2, 0x0a, 0x92, 0x6a,
	0x4a, 0x05, 0x55, 0x0a, 0x2a, 0x4a, 0x5a, 0x02, 0xb5, 0x02, 0xb2, 0x31, 0x69, 0x03,
	0x31, 0x73, 0xa9, 0x0a, 0x4a, 0x05, 0x2d, 0x55, 0x2d, 0x09, 0x5a, 0x01, 0xd5, 0x48,
	0xb4, 0x09, 0x68, 0x89, 0x54, 0x0b, 0xa4, 0x0a, 0xa5, 0x6a, 0x95, 0x04, 0xad, 0x08,
	0x6a, 0x44, 0xda, 0x04, 0x74, 0x05, 0xb0, 0x25, 0x54, 0x03,
}

var bigLeapMonthYears = []int{
	6, 14, 19, 25, 33, 36, 38, 41, 44, 52,
	55, 79, 117, 136, 147, 150, 155, 158, 185, 193,
}

func daysInChinMonth(month, year int) (days int) {
	i := year - bcy
	if i > 202 || i < 0 {
		panic("200 years")
	}
	days = 30
	var v, l uint8
	if 1 <= month && month <= 8 {
		v = chMon[i*2]
		l = uint8(month) - 1
		if v>>l&1 == 1 {
			days = 29
		}
	} else if 9 <= month && month <= 12 {
		v = chMon[i*2+1]
		l = uint8(month) - 9
		if v>>l&1 == 1 {
			days = 29
		}
	} else {
		// leap chinese month
		v = chMon[i*2+1] >> 4
		if int(v) != abs(month) {
			days = 0
		} else {
			days = 29
			for _, m := range bigLeapMonthYears {
				if m == i {
					days = 30
					break
				}
			}
		}
	}
	return
}

func nextChinMonth(month, year int) int {
	nm := abs(month) + 1
	if month > 0 {
		i := year - bcy
		v := chMon[2*i+1] >> 4
		if int(v) == month {
			return -month
		}
	}
	if nm == 13 {
		nm = 1
	}
	return nm
}

func abs(v int) int {
	if v < 0 {
		v = -v
	}
	return v
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		printCSV(os.Stdout, example)
	} else {
		fp, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Println(err)
			flag.Usage()
			os.Exit(1)
		}
		defer fp.Close()
		p, err := ioutil.ReadAll(fp)
		if err != nil {
			fmt.Println(err)
			flag.Usage()
			os.Exit(1)
		}
		fn := "import.csv"
		path := filepath.Dir(flag.Arg(0))
		ofn := filepath.Join(path, fn)
		ofp, err := os.Create(ofn)
		if err != nil {
			fmt.Printf("unable to open %s in %s\n", fn, path)
			os.Exit(1)
		}
		defer ofp.Close()
		printCSV(ofp, string(p))
	}
}

func printCSV(w io.Writer, content string) {
	items := strings.Split(content, "\n")
	fmt.Fprintln(w, "Subject, Start Date, Start Time, End Date, End Time, Private, All Day Event, Location")
	for _, item := range items {
		if strings.TrimSpace(item) == "" {
			continue
		}
		if strings.HasPrefix(item, "#") {
			continue
		}
		birthday := strings.Split(item, ":")
		if len(birthday) < 2 {
			fmt.Println("wrong file format")
			os.Exit(1)
		}
		desc := strings.TrimSpace(birthday[0])
		day := strings.TrimSpace(birthday[1])
		date := strings.Split(day, ".")
		if len(date) != 3 {
			fmt.Println("date format wrong, should be yyyy-mm-dd")
			os.Exit(1)
		}
		year, err := strconv.Atoi(date[0])
		if err != nil {
			fmt.Println("number format wrong, should be yyyy-mm-dd")
		}
		month, err := strconv.Atoi(date[1])
		if err != nil {
			fmt.Println("number format wrong, should be yyyy-mm-dd")
		}
		mday, err := strconv.Atoi(date[2])
		if err != nil {
			fmt.Println("number format wrong, should be yyyy-mm-dd")
		}
		for y := year; y <= 2100; y++ {
			cdate := newDate(y, month, mday)
			gdate := chinese2gregorian(cdate)
			if gdate == nil {
				continue
			}
			fdate := fmt.Sprintf("%d-%d-%d", gdate.year, gdate.month, gdate.day)
			fmt.Fprintf(w, "%s, %s, 8:00 AM, %s, , FALSE, TRUE, HOME\n", desc, fdate, fdate)
		}
	}
}
