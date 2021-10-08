package structs

import (
	"testing"
	"time"
)

type Student struct {
	Info
	Name    string      `json:"name"`
	Age     int         `json:"age"`
	Address Address     `json:"address"`
	Class   interface{} `json:"class"`
}

type Address struct {
	Province string `json:"province"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Detail   string `json:"detail"`
}

type Info struct {
	Id        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func TestToMap(t *testing.T) {
	student := Student{
		Info: Info{
			Id:        1,
			CreatedAt: time.Now(),
		},
		Name: "hello",
		Age:  133,
		Address: Address{
			Province: "zhejiang",
			City:     "hangzhou",
		},
		Class: map[string]interface{}{
			"hello": "world",
			"abc":   123,
		},
	}

	m := ToMap(&student, "json")
	t.Logf("%+v", m)
}

func TestMergeMap(t *testing.T) {
	m1 := map[string]interface{}{
		"abc":   123,
		"hello": "world",
	}

	m2 := map[string]interface{}{
		"aaa": "haha",
		"abc": 234,
	}

	MergeMap(m1, m2)
	t.Log(m1)
}

func TestMergeStruct(t *testing.T) {
	a1 := Address{
		Province: "zhejiang",
		City:     "hangzhou",
		Region:   "gs",
		Detail:   "jjyz",
	}
	a2 := Address{
		Province: "zhejiang2",
		City:     "hangzhou2",
	}
	a3 := Address{
		Province: "",
		City:     "hangzhou3",
	}

	MergeStruct(&a1, &a2, a3)
	t.Logf("%+v", a1)
}
