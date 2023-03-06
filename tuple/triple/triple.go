package ekit

import "fmt"

type Triple struct {
	First  any
	Second any
	Third  any
}

// ToString 将 Triple 转为字符串，格式类似于 <key,value>
func (t Triple) ToString() string {
	return fmt.Sprint("<", t.First, ",", t.Second, ",", t.Third, ">")
}

// ToList 将 Triple 转为数组
func (t Triple) ToList() []any {
	return []any{t.First, t.Second, t.Third}
}

// Copy 传入一个 Triple 来修改对应位置的值
func (t Triple) Copy(toTriple Triple) (newTriple Triple) {
	newTriple = Triple{First: t.First, Second: t.Second}
	if toTriple.First != nil {
		newTriple.First = toTriple.First
	}
	if toTriple.Second != nil {
		newTriple.Second = toTriple.Second
	}
	if toTriple.Third != nil {
		newTriple.Third = toTriple.Third
	}
	return
}
