/*
 * Copyright (c) 2020 Mikhail Knyazhev <markus621@gmail.com>.
 * All rights reserved. Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package app

import "context"

type ForceClose struct {
	C     context.Context
	Close context.CancelFunc
}

func newForceClose() *ForceClose {
	ctx, cncl := context.WithCancel(context.Background())

	return &ForceClose{
		C:     ctx,
		Close: cncl,
	}
}
