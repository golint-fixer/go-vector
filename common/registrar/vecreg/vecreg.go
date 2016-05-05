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

package vecreg

import (
	"math/big"

	"github.com/vector/go-vector/common/registrar"
	"github.com/vector/go-vector/xvec"
)

// implements a versioned Registrar on an archiving full node
type EthReg struct {
	backend  *xvec.XEth
	registry *registrar.Registrar
}

func New(xe *xvec.XEth) (self *EthReg) {
	self = &EthReg{backend: xe}
	self.registry = registrar.New(xe)
	return
}

func (self *EthReg) Registry() *registrar.Registrar {
	return self.registry
}

func (self *EthReg) Resolver(n *big.Int) *registrar.Registrar {
	xe := self.backend
	if n != nil {
		xe = self.backend.AtStateNum(n.Int64())
	}
	return registrar.New(xe)
}
