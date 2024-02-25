package utils_test

import (
	"testing"

	"github.com/saweima12/gdit/internal/utils"
)

type TestItem struct{}

func TestHelper(t *testing.T) {

	t.Run("The name must be `utils_test.TestItem`", func(ct *testing.T) {
		name := utils.GetType[TestItem]()
		if name != "utils_test.TestItem" {
			t.Fail()
		}
	})

	t.Run("Should be 0, emptyString and nil ", func(ct *testing.T) {
		if utils.Empty[int]() != 0 {
			t.Fail()
		}
		if utils.Empty[string]() != "" {
			t.Fail()
		}
		if utils.Empty[*TestItem]() != nil {
			t.Fail()
		}
	})

	t.Run("Should be `utils_test.TestItem` and false", func(ct *testing.T) {
		val, named := utils.GetProviderKey[TestItem]("")
		if val != "utils_test.TestItem" || named {
			t.Fail()
		}
	})

	t.Run("Should be `TestItem` and true ", func(ct *testing.T) {
		val, named := utils.GetProviderKey[TestItem]("TestItem")
		if val != "TestItem" || !named {
			t.Fail()
		}
	})
}
