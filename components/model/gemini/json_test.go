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
	"math"
	"reflect"
	"testing"

	"github.com/bytedance/sonic"
)

var sortedSonic = sonic.Config{SortMapKeys: true}.Froze()

// TestMarshalStringMatchesSonic: tool call arguments as Gemini returns them
// (map[string]any) are written with the bytes of sonic's default
// configuration, map keys sorted.
func TestMarshalStringMatchesSonic(t *testing.T) {
	for _, args := range []map[string]any{
		nil,
		{},
		{"city": "Roma <centro> & dintorni", "days": 3.0, "ratio": 0.1, "big": 1e21, "small": 1e-7, "neg": -0.0},
		{"nested": map[string]any{"list": []any{"a", 1.5, nil, true, map[string]any{"k": "\u2028\u2029"}}}, "empty": []any{}},
		{"escapes": "quote\" back\\ slash/ tab\t nl\n cr\r ctl\x01", "unicode": "è 😀", "max": math.MaxFloat64, "int": 9007199254740993.0},
	} {
		got, gotErr := marshalString(args)
		want, wantErr := sortedSonic.MarshalToString(args)
		if (gotErr == nil) != (wantErr == nil) || got != want {
			t.Errorf("%v:\n got %s (%v)\nwant %s (%v)", args, got, gotErr, want, wantErr)
		}
	}
}

// TestUnmarshalStringMatchesSonic: tool results and call arguments read into
// map[string]any give what sonic gave, or an error in both.
func TestUnmarshalStringMatchesSonic(t *testing.T) {
	for _, in := range []string{
		`{"output":"ok","items":[1,2.5,"x",null,true,{"k":"v"}],"n":12345678901234567890}`,
		`{"a":1,"a":2}`,
		`{"html":"<b>&amp;</b>","esc":"\u00e8\ud83d\ude00\n"}`,
		`{"big":1e400}`,
		`plain text result`,
		`["not","an","object"]`,
		`"string"`,
		`null`,
		``,
		`{"a":1} trailing`,
		`{"a":}`,
	} {
		got, want := map[string]any{}, map[string]any{}
		gotErr := unmarshalString(in, &got)
		wantErr := sonic.UnmarshalString(in, &want)
		if (gotErr == nil) != (wantErr == nil) || (gotErr == nil && !reflect.DeepEqual(got, want)) {
			t.Errorf("%q:\n got %#v (%v)\nwant %#v (%v)", in, got, gotErr, want, wantErr)
		}
	}
}
