package bus

import (
	"fmt"
	"reflect"
)

func getRegistryKey(t reflect.Type) (string, error) {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	name, pkgPath := t.Name(), t.PkgPath()
	if name == "" || pkgPath == "" {
		return "", fmt.Errorf("empty type name or package path, name:%q, pkg:%q", name, pkgPath)
	}

	return name + "@" + pkgPath, nil
}
