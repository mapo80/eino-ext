/*
 * Copyright 2026 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gemini

import (
	"encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

// jsonOptions are the semantics of encoding/json with the escaping sonic's
// default configuration had: '<', '>', '&', U+2028 and U+2029 are written
// as they are. Map keys are written sorted.
var jsonOptions = jsonv2.JoinOptions(
	json.DefaultOptionsV1(),
	jsontext.EscapeForHTML(false),
	jsontext.EscapeForJS(false),
)

func marshalString(v any) (string, error) {
	b, err := jsonv2.Marshal(v, jsonOptions)
	return string(b), err
}

func unmarshalString(s string, v any) error {
	return jsonv2.Unmarshal([]byte(s), v, jsonOptions)
}
