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

// Package protein is an encoding/decoding library for Protobuf that comes with
// schema-versioning and runtime-decoding capabilities.
//
// It has diverse use-cases, including but not limited to:
//   - setting up schema registries
//   - decoding Protobuf payloads without the need to know their schema at
//     compile-time
//   - identifying & preventing applicative bugs and data corruption issues
//   - creating custom-made container formats for on-disk storage
//   - ...and more!
package protein
