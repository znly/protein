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

// Package protoscan provides the necessary tools & APIs to find, extract,
// version and build the dependency trees of all the protobuf schemas that
// have been instanciated by one or more protobuf library (golang/protobuf,
// gogo/protobuf...).
//
// This is a fairly low-level package, used to power the deeper innards of
// `protein`; as such, it should very rarely be of use to the end-user.
package protoscan
