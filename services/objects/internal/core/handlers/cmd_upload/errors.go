package cmd_upload

import "errors"

var ErrObjectExists = errors.New("objects.core.upload-command:object exists")
var ErrObfuscatedResult = errors.New("objects.core.upload-command:could not fetch existing object")
