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

package vec

import (
	"github.com/vector/go-vector/metrics"
	"github.com/vector/go-vector/p2p"
)

var (
	propTxnInPacketsMeter     = metrics.NewMeter("vec/prop/txns/in/packets")
	propTxnInTrafficMeter     = metrics.NewMeter("vec/prop/txns/in/traffic")
	propTxnOutPacketsMeter    = metrics.NewMeter("vec/prop/txns/out/packets")
	propTxnOutTrafficMeter    = metrics.NewMeter("vec/prop/txns/out/traffic")
	propHashInPacketsMeter    = metrics.NewMeter("vec/prop/hashes/in/packets")
	propHashInTrafficMeter    = metrics.NewMeter("vec/prop/hashes/in/traffic")
	propHashOutPacketsMeter   = metrics.NewMeter("vec/prop/hashes/out/packets")
	propHashOutTrafficMeter   = metrics.NewMeter("vec/prop/hashes/out/traffic")
	propBlockInPacketsMeter   = metrics.NewMeter("vec/prop/blocks/in/packets")
	propBlockInTrafficMeter   = metrics.NewMeter("vec/prop/blocks/in/traffic")
	propBlockOutPacketsMeter  = metrics.NewMeter("vec/prop/blocks/out/packets")
	propBlockOutTrafficMeter  = metrics.NewMeter("vec/prop/blocks/out/traffic")
	reqHashInPacketsMeter     = metrics.NewMeter("vec/req/hashes/in/packets")
	reqHashInTrafficMeter     = metrics.NewMeter("vec/req/hashes/in/traffic")
	reqHashOutPacketsMeter    = metrics.NewMeter("vec/req/hashes/out/packets")
	reqHashOutTrafficMeter    = metrics.NewMeter("vec/req/hashes/out/traffic")
	reqBlockInPacketsMeter    = metrics.NewMeter("vec/req/blocks/in/packets")
	reqBlockInTrafficMeter    = metrics.NewMeter("vec/req/blocks/in/traffic")
	reqBlockOutPacketsMeter   = metrics.NewMeter("vec/req/blocks/out/packets")
	reqBlockOutTrafficMeter   = metrics.NewMeter("vec/req/blocks/out/traffic")
	reqHeaderInPacketsMeter   = metrics.NewMeter("vec/req/headers/in/packets")
	reqHeaderInTrafficMeter   = metrics.NewMeter("vec/req/headers/in/traffic")
	reqHeaderOutPacketsMeter  = metrics.NewMeter("vec/req/headers/out/packets")
	reqHeaderOutTrafficMeter  = metrics.NewMeter("vec/req/headers/out/traffic")
	reqBodyInPacketsMeter     = metrics.NewMeter("vec/req/bodies/in/packets")
	reqBodyInTrafficMeter     = metrics.NewMeter("vec/req/bodies/in/traffic")
	reqBodyOutPacketsMeter    = metrics.NewMeter("vec/req/bodies/out/packets")
	reqBodyOutTrafficMeter    = metrics.NewMeter("vec/req/bodies/out/traffic")
	reqStateInPacketsMeter    = metrics.NewMeter("vec/req/states/in/packets")
	reqStateInTrafficMeter    = metrics.NewMeter("vec/req/states/in/traffic")
	reqStateOutPacketsMeter   = metrics.NewMeter("vec/req/states/out/packets")
	reqStateOutTrafficMeter   = metrics.NewMeter("vec/req/states/out/traffic")
	reqReceiptInPacketsMeter  = metrics.NewMeter("vec/req/receipts/in/packets")
	reqReceiptInTrafficMeter  = metrics.NewMeter("vec/req/receipts/in/traffic")
	reqReceiptOutPacketsMeter = metrics.NewMeter("vec/req/receipts/out/packets")
	reqReceiptOutTrafficMeter = metrics.NewMeter("vec/req/receipts/out/traffic")
	miscInPacketsMeter        = metrics.NewMeter("vec/misc/in/packets")
	miscInTrafficMeter        = metrics.NewMeter("vec/misc/in/traffic")
	miscOutPacketsMeter       = metrics.NewMeter("vec/misc/out/packets")
	miscOutTrafficMeter       = metrics.NewMeter("vec/misc/out/traffic")
)

// meteredMsgReadWriter is a wrapper around a p2p.MsgReadWriter, capable of
// accumulating the above defined metrics based on the data stream contents.
type meteredMsgReadWriter struct {
	p2p.MsgReadWriter     // Wrapped message stream to meter
	version           int // Protocol version to select correct meters
}

// newMeteredMsgWriter wraps a p2p MsgReadWriter with metering support. If the
// metrics system is disabled, this fucntion returns the original object.
func newMeteredMsgWriter(rw p2p.MsgReadWriter) p2p.MsgReadWriter {
	if !metrics.Enabled {
		return rw
	}
	return &meteredMsgReadWriter{MsgReadWriter: rw}
}

// Init sets the protocol version used by the stream to know which meters to
// increment in case of overlapping message ids between protocol versions.
func (rw *meteredMsgReadWriter) Init(version int) {
	rw.version = version
}

func (rw *meteredMsgReadWriter) ReadMsg() (p2p.Msg, error) {
	// Read the message and short circuit in case of an error
	msg, err := rw.MsgReadWriter.ReadMsg()
	if err != nil {
		return msg, err
	}
	// Account for the data traffic
	packets, traffic := miscInPacketsMeter, miscInTrafficMeter
	switch {
	case rw.version < vec62 && msg.Code == BlockHashesMsg:
		packets, traffic = reqHashInPacketsMeter, reqHashInTrafficMeter
	case rw.version < vec62 && msg.Code == BlocksMsg:
		packets, traffic = reqBlockInPacketsMeter, reqBlockInTrafficMeter

	case rw.version >= vec62 && msg.Code == BlockHeadersMsg:
		packets, traffic = reqHeaderInPacketsMeter, reqHeaderInTrafficMeter
	case rw.version >= vec62 && msg.Code == BlockBodiesMsg:
		packets, traffic = reqBodyInPacketsMeter, reqBodyInTrafficMeter

	case rw.version >= vec63 && msg.Code == NodeDataMsg:
		packets, traffic = reqStateInPacketsMeter, reqStateInTrafficMeter
	case rw.version >= vec63 && msg.Code == ReceiptsMsg:
		packets, traffic = reqReceiptInPacketsMeter, reqReceiptInTrafficMeter

	case msg.Code == NewBlockHashesMsg:
		packets, traffic = propHashInPacketsMeter, propHashInTrafficMeter
	case msg.Code == NewBlockMsg:
		packets, traffic = propBlockInPacketsMeter, propBlockInTrafficMeter
	case msg.Code == TxMsg:
		packets, traffic = propTxnInPacketsMeter, propTxnInTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
	// Account for the data traffic
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	switch {
	case rw.version < vec62 && msg.Code == BlockHashesMsg:
		packets, traffic = reqHashOutPacketsMeter, reqHashOutTrafficMeter
	case rw.version < vec62 && msg.Code == BlocksMsg:
		packets, traffic = reqBlockOutPacketsMeter, reqBlockOutTrafficMeter

	case rw.version >= vec62 && msg.Code == BlockHeadersMsg:
		packets, traffic = reqHeaderOutPacketsMeter, reqHeaderOutTrafficMeter
	case rw.version >= vec62 && msg.Code == BlockBodiesMsg:
		packets, traffic = reqBodyOutPacketsMeter, reqBodyOutTrafficMeter

	case rw.version >= vec63 && msg.Code == NodeDataMsg:
		packets, traffic = reqStateOutPacketsMeter, reqStateOutTrafficMeter
	case rw.version >= vec63 && msg.Code == ReceiptsMsg:
		packets, traffic = reqReceiptOutPacketsMeter, reqReceiptOutTrafficMeter

	case msg.Code == NewBlockHashesMsg:
		packets, traffic = propHashOutPacketsMeter, propHashOutTrafficMeter
	case msg.Code == NewBlockMsg:
		packets, traffic = propBlockOutPacketsMeter, propBlockOutTrafficMeter
	case msg.Code == TxMsg:
		packets, traffic = propTxnOutPacketsMeter, propTxnOutTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	// Send the packet to the p2p layer
	return rw.MsgReadWriter.WriteMsg(msg)
}
