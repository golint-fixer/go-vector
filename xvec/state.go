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

package xvec

import (
	"github.com/vector/go-vector/common"
	"github.com/vector/go-vector/core/state"
)

type State struct {
	xvec  *XEth
	state *state.StateDB
}

func NewState(xvec *XEth, statedb *state.StateDB) *State {
	return &State{xvec, statedb}
}

func (self *State) State() *state.StateDB {
	return self.state
}

func (self *State) Get(addr string) *Object {
	return &Object{self.state.GetStateObject(common.HexToAddress(addr))}
}

func (self *State) SafeGet(addr string) *Object {
	return &Object{self.safeGet(addr)}
}

func (self *State) safeGet(addr string) *state.StateObject {
	object := self.state.GetStateObject(common.HexToAddress(addr))
	if object == nil {
		object = state.NewStateObject(common.HexToAddress(addr), self.xvec.backend.ChainDb())
	}

	return object
}
