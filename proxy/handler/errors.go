// Copyright (c) 2019 ELEME, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"github.com/juju/errors"
	"tedis/proxy/log"
)

var (
	ErrBegionTXN = errors.New("begin transaction error")
	ErrKeySize   = errors.New("invalid key size")
	ErrValueSize = errors.New("invalid value size")
	ErrValueNil  = errors.New("value is null")
	ErrAuthPwd   = errors.New("Invalid password")
)

func errArguments(format string, v ...interface{}) error {
	err := errors.Errorf(format, v...)
	log.Warningf("call store function with invalid arguments - %s", err)
	return errors.Trace(err)
}
