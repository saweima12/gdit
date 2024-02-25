package ext_test

import (
	"fmt"
	"testing"

	"github.com/saweima12/gdit/internal/ext"
)

func TestGenericSyncMap(t *testing.T) {

	m := &ext.GSyncMap[int]{}
	m.Store("age", 5)

	t.Run("The load value should equals 5", func(ct *testing.T) {
		val, ok := m.Load("age")
		if val != 5 || !ok {
			t.Fail()
		}
	})

	t.Run("The load value is not exists should be zero", func(ct *testing.T) {
		val, ok := m.Load("agec")
		if ok || val != 0 {
			t.Fail()
		}
	})

	t.Run("The old value should be 5 and new value is 10", func(ct *testing.T) {
		preVal, _ := m.Swap("age", 10)
		nVal, _ := m.Load("age")

		if preVal != 5 || nVal != 10 {
			t.Fail()
		}
	})

	t.Run("The swap value is not exists should be 0", func(ct *testing.T) {
		preVal, _ := m.Swap("agec", 10)
		if preVal != 0 {
			t.Fail()
		}
	})

	t.Run("The Range method will iterate over all values.", func(t *testing.T) {
		items := []int{}
		m.Range(func(key string, value int) bool {
			items = append(items, value)
			return true
		})
		fmt.Println(items)
		if len(items) < 1 {
			t.Fail()
		}
	})

}
