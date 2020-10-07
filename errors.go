/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import (
	"github.com/pkg/errors"
)

func WrapErrors(err1, err2 error, message string) error {
	if err2 == nil {
		return err1
	}
	if err1 == nil {
		return errors.Wrap(err2, message)
	}
	return errors.Wrap(err1, errors.Wrap(err2, message).Error())
}
