// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package hub

import "errors"

// ErrMergeFailed indicates hub store merge could not be applied.
var ErrMergeFailed = errors.New("hub merge failed")
