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

package openai

import (
	"encoding/json"
	"testing"

	"github.com/bytedance/sonic"

	"github.com/cloudwego/eino/schema"
)

// TestReasoningFromExtraMatchesSonic: the reasoning content read from a
// delta's extra fields is the string sonic read, for every raw value the
// providers send (strings with escapes, empty, null, numbers, objects,
// malformed), and absent when sonic failed.
func TestReasoningFromExtraMatchesSonic(t *testing.T) {
	for _, raw := range []string{
		`"step one: compare <prices> & dates"`,
		`"è😀\n\t\"quoted\""`,
		`""`,
		`null`,
		`12`,
		`{"text":"x"}`,
		`["x"]`,
		`"unterminated`,
		`"bad \x"`,
	} {
		for _, key := range otherReasoningKeys {
			msg := &schema.Message{}
			populateRCFromExtra(map[string]json.RawMessage{key: json.RawMessage(raw)}, msg)
			var want string
			if err := sonic.Unmarshal([]byte(raw), &want); err != nil {
				want = ""
			}
			if msg.ReasoningContent != want {
				t.Errorf("%s=%s: got %q, sonic %q", key, raw, msg.ReasoningContent, want)
			}
		}
	}
}
