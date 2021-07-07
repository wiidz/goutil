package structHelper

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

//集合模拟
type Set struct {
	m map[int]bool
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		m: map[int]bool{},
	}
}

func (s *Set) Add(item int) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item int) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[int]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) List() []int {
	s.RLock()
	defer s.RUnlock()
	list := []int{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
func (s *Set) SortList() []int {
	s.RLock()
	defer s.RUnlock()
	list := []int{}
	for item := range s.m {
		list = append(list, item)
	}
	sort.Ints(list)
	return list
}

/**
 * @func: GetStructAllMethods  打印结构体的所有方法
 * @author Wiidz
 * @date   2019-11-16
 */
func  GetStructAllMethods(target interface{}) {
	value := reflect.ValueOf(target)
	typ := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		fmt.Println(fmt.Sprintf("method [%d]%s and type is %v", i, typ.Method(i).Name, typ.Method(i).Type))
	}
}

/**
 * @func: GetStructAllFields  打印结构体的所有键值
 * @author Wiidz
 * @date   2019-11-16
 */
func  GetStructAllFields(target interface{}) {
	value := reflect.ValueOf(target)
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		fmt.Println(fmt.Sprintf("field [%d]%s and type is %v", i, typ.Field(i).Name, typ.Field(i).Type))
	}
}

/**
 * @func: GetStructAllFields  打印结构体的所有键值对
 * @author Wiidz
 * @date   2019-11-16
 */
func  GetStructAllKV(target interface{}) {
	value := reflect.ValueOf(target)
	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		//fmt.Println(fmt.Sprintf("field [%d]%s and type is %v", i, typ.Field(i).Name, typ.Field(i).Type))
		fmt.Println("["+strconv.Itoa(i)+"]", typ.Field(i).Name, typ.Field(i).Type, typ.Field(i))
	}
}


/**
 * @func: GetStructAllFields  打印结构体的所有键值对
 * @author Wiidz
 * @date   2019-11-16
 */
func  GetStructPointAllKV(target interface{}) {

	value := reflect.ValueOf(target)
	typ := value.Elem().Type()

	log.Println("value",value)
	log.Println("typ",typ)

	for i := 0; i < value.Elem().NumField(); i++ {
		//fmt.Println(fmt.Sprintf("field [%d]%s and type is %v", i, typ.Field(i).Name, typ.Field(i).Type))
		fmt.Println("["+strconv.Itoa(i)+"]", typ.Field(i).Name, typ.Field(i).Type, typ.Field(i))
	}
}
