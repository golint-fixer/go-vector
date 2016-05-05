// Copyright 2015 The go-vector Authors
// This file is part of the go-vector library.
//
// The go-vector library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vector library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vector library. If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/vector/go-vector/metrics"
)

var (
	hashInMeter      = metrics.NewMeter("vec/downloader/hashes/in")
	hashReqTimer     = metrics.NewTimer("vec/downloader/hashes/req")
	hashDropMeter    = metrics.NewMeter("vec/downloader/hashes/drop")
	hashTimeoutMeter = metrics.NewMeter("vec/downloader/hashes/timeout")

	blockInMeter      = metrics.NewMeter("vec/downloader/blocks/in")
	blockReqTimer     = metrics.NewTimer("vec/downloader/blocks/req")
	blockDropMeter    = metrics.NewMeter("vec/downloader/blocks/drop")
	blockTimeoutMeter = metrics.NewMeter("vec/downloader/blocks/timeout")

	headerInMeter      = metrics.NewMeter("vec/downloader/headers/in")
	headerReqTimer     = metrics.NewTimer("vec/downloader/headers/req")
	headerDropMeter    = metrics.NewMeter("vec/downloader/headers/drop")
	headerTimeoutMeter = metrics.NewMeter("vec/downloader/headers/timeout")

	bodyInMeter      = metrics.NewMeter("vec/downloader/bodies/in")
	bodyReqTimer     = metrics.NewTimer("vec/downloader/bodies/req")
	bodyDropMeter    = metrics.NewMeter("vec/downloader/bodies/drop")
	bodyTimeoutMeter = metrics.NewMeter("vec/downloader/bodies/timeout")

	receiptInMeter      = metrics.NewMeter("vec/downloader/receipts/in")
	receiptReqTimer     = metrics.NewTimer("vec/downloader/receipts/req")
	receiptDropMeter    = metrics.NewMeter("vec/downloader/receipts/drop")
	receiptTimeoutMeter = metrics.NewMeter("vec/downloader/receipts/timeout")

	stateInMeter      = metrics.NewMeter("vec/downloader/states/in")
	stateReqTimer     = metrics.NewTimer("vec/downloader/states/req")
	stateDropMeter    = metrics.NewMeter("vec/downloader/states/drop")
	stateTimeoutMeter = metrics.NewMeter("vec/downloader/states/timeout")
)
