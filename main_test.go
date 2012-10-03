package main

import (
	"fmt"
	"testing"
)

func TestGreg2Chin(t *testing.T) {
	greg := newDate(1901, 1, 1)
	chin := newDate(1900, 11, 11)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(1902, 3, 9)
	chin = newDate(1902, 1, 30)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v", *tmp, *chin)
	}
	greg = newDate(1906, 6, 21)
	chin = newDate(1906, -4, 30)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(1984, 11, 11)
	chin = newDate(1984, 10, 19)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(1984, 12, 20)
	chin = newDate(1984, -10, 28)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(1985, 6, 6)
	chin = newDate(1985, 4, 18)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(1999, 11, 16)
	chin = newDate(1999, 10, 9)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(2001, 10, 14)
	chin = newDate(2001, 8, 28)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
	greg = newDate(2012, 10, 03)
	chin = newDate(2012, 8, 18)
	if tmp := gregorian2chinese(greg); *tmp != *chin {
		t.Errorf("got: %v; expect: %v\n", *tmp, *chin)
	}
}

func TestChin2Greg(t *testing.T) {
	greg := newDate(1901, 2, 19)
	chin := newDate(1901, 1, 1)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}

	greg = newDate(1907, 2, 13)
	chin = newDate(1907, 1, 1)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1906, 6, 21)
	chin = newDate(1906, -4, 30)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1963, 5, 23)
	chin = newDate(1963, -4, 1)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1907, 2, 12)
	chin = newDate(1906, 12, 30)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1910, 4, 13)
	chin = newDate(1910, 3, 4)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1984, 11, 11)
	chin = newDate(1984, 10, 19)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1984, 12, 11)
	chin = newDate(1984, -10, 19)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(1985, 6, 6)
	chin = newDate(1985, 4, 18)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(2012, 10, 3)
	chin = newDate(2012, 8, 18)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
	greg = newDate(2012, 10, 3)
	chin = newDate(2012, 8, 18)
	if tmp := chinese2gregorian(chin); *tmp != *greg {
		t.Errorf("got: %v; expect: %v\n", *tmp, *greg)
	}
}

func TestShouldFail(t *testing.T) {
	c1 := newDate(2110, 8, 18)
	c2 := newDate(2012, 13, 3)
	c3 := newDate(1984, -1, 1)
	c4 := newDate(1984, 10, 31)
	c5 := newDate(1963, 2, 30)
	for _, c := range []*date{c1, c2, c3, c4, c5} {
		if isValidChin(c) {
			t.Errorf("%v should fail", c)
		}
	}

	g1 := newDate(2100, 2, 29)
	g2 := newDate(2110, 1, 1)
	g3 := newDate(2012, 13, 1)
	g4 := newDate(2012, 6, 31)
	for _, g := range []*date{g1, g2, g3, g4} {
		if isValidGreg(g) {
			t.Errorf("%v should fail", g)
		}
	}
}

func TestNumDaysInYears(t *testing.T) {
	for y := 1901; y < 2101; y++ {
		var days int
		for m := 1; ; {
			days += daysInChinMonth(m, y)
			//fmt.Println("this month:", m, "days:", daysInChinMonth(m, y))
			m = nextChinMonth(m, y)
			if m == 1 {
				break
			}
		}
		// fmt.Printf("%d:%d\n", y, days)
	}
}

func TestDebug(t *testing.T) {
	fmt.Println("")
}
