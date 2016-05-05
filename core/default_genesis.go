// Copyright 2014 The go-vector Authors
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

package core

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"
)

func NewDefaultGenesisReader() (io.Reader, error) {
	return gzip.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(defaultGenesisBlock)))
}

const defaultGenesisBlock = "H4sIAAAJbogA/6yRz04DIRDG34VzD/xbGHqrrdGDGhN9gRlgLQndNbs0WdPsu4u7B2OihybOgQSY7/fxMRf21Hc+si3jE/9RSrINe02nOBY8vS8NjZVa8sbWi2ccYlfucTz+Ir2+KvF2KgMesOACFIIoQiAdlbZ1Be+U9cKDsjxqamJQCgOpQMZ5CmRRQRSx4SQIZIsVeIfjQzqlsvCM2EE9O6S2Tf6cy8fq8m3/mKbjP6bZ96kjHNePBRO15h7AmADSeLJSKieNFOCAtzUbOqe0qrpdzr1n28s1qtq978OXVQW8lH7At7rpzjlv2A1mXOcr9J/Pnef5MwAA//9ygBnhCAIAAA=="