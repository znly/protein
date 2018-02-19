// Copyright © 2016 Zenly <hello@zen.ly>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build darwin

package protoscan

import (
	"debug/macho"
	"strings"

	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

func scanSymbols(binPath string) ([]symbol, error) {
	f, err := macho.Open(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	protoFiles := []symbol{}
	for _, s := range f.Symtab.Syms {
		if strings.HasSuffix(s.Name, "/proto.protoFiles") {
			protoFiles = append(protoFiles, symbol{s.Name, uintptr(s.Value)})
		}
	}

	return protoFiles, nil
}
