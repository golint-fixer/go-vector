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

// Contains the metrics collected by the fetcher.

package fetcher

import (
	"github.com/vector/go-vector/metrics"
)

var (
	propAnnounceInMeter   = metrics.NewMeter("vec/fetcher/prop/announces/in")
	propAnnounceOutTimer  = metrics.NewTimer("vec/fetcher/prop/announces/out")
	propAnnounceDropMeter = metrics.NewMeter("vec/fetcher/prop/announces/drop")
	propAnnounceDOSMeter  = metrics.NewMeter("vec/fetcher/prop/announces/dos")

	propBroadcastInMeter   = metrics.NewMeter("vec/fetcher/prop/broadcasts/in")
	propBroadcastOutTimer  = metrics.NewTimer("vec/fetcher/prop/broadcasts/out")
	propBroadcastDropMeter = metrics.NewMeter("vec/fetcher/prop/broadcasts/drop")
	propBroadcastDOSMeter  = metrics.NewMeter("vec/fetcher/prop/broadcasts/dos")

	blockFetchMeter  = metrics.NewMeter("vec/fetcher/fetch/blocks")
	headerFetchMeter = metrics.NewMeter("vec/fetcher/fetch/headers")
	bodyFetchMeter   = metrics.NewMeter("vec/fetcher/fetch/bodies")

	blockFilterInMeter   = metrics.NewMeter("vec/fetcher/filter/blocks/in")
	blockFilterOutMeter  = metrics.NewMeter("vec/fetcher/filter/blocks/out")
	headerFilterInMeter  = metrics.NewMeter("vec/fetcher/filter/headers/in")
	headerFilterOutMeter = metrics.NewMeter("vec/fetcher/filter/headers/out")
	bodyFilterInMeter    = metrics.NewMeter("vec/fetcher/filter/bodies/in")
	bodyFilterOutMeter   = metrics.NewMeter("vec/fetcher/filter/bodies/out")
)
