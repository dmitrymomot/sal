package looker_test

import (
	"reflect"
	"testing"

	"github.com/go-gad/sal/looker"
	"github.com/go-gad/sal/looker/bookstore"
)

func TestLookAtInterface(t *testing.T) {
	pf := getLogger(t)
	var typ reflect.Type = reflect.TypeOf((*bookstore.StoreClient)(nil)).Elem()
	intf := looker.LookAtInterface(typ)
	pf("Interface %#v", intf)

}

func getLogger(t *testing.T) func(string, ...interface{}) {
	return func(s string, kv ...interface{}) {
		t.Logf(s, kv...)
	}
}
