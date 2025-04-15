package bus

import (
	"errors"
	"reflect"
)

func getRegistryKey(t reflect.Type) (string, error) {
	name, pkgPath := t.Name(), t.PkgPath()
	if name == "" {
		return "", errors.New("") // TODO
	}
	if pkgPath == "" {
		return "", errors.New("") // TODO
	}

	return name + "@" + pkgPath, nil
}
