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

// Package failure lists all the possible errors that can be returned either by
// `protein` or any of its sub-packages.
//
// Protein uses the pkg/errors[1] package to handle error propagation throughout
// the call stack; please take a look at the related documentation for more
// information on how to properly handle these errors.
//
// [1]: https://github.com/pkg/errors
package failure
