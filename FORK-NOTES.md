# Fork mapo80/eino-ext

Fork of `github.com/cloudwego/eino-ext` for ARC (`arc-workspace`). Module paths are unchanged; ARC uses the sub-modules with
`replace github.com/cloudwego/eino-ext/<module> => github.com/mapo80/eino-ext/<module> <tag>`.

- **License:** unchanged (`LICENSE`, Apache 2.0 as upstream).
- **Go:** `go 1.27` (`encoding/json/v2` and `encoding/json/jsontext` from the standard library).
- **Goal:** one JSON library in ARC's server build, the standard library's `encoding/json/v2`: `bytedance/sonic` out of the modules ARC links.
- **Branches:** one per module, each based on the version ARC used, recorded in `main` with a merge that keeps main's content.

## Releases

| tag | module | base | branch |
|---|---|---|---|
| `libs/acl/openai/v0.1.17-arc.1` | `libs/acl/openai` | tag `libs/acl/openai/v0.1.17` (`6dd3b2b7`) | `arc/json-v2-acl-openai` |
| `components/model/gemini/v0.1.37-arc.1` | `components/model/gemini` | commit `2607f61e` (ARC used `v0.1.37-0.20260916065400-2607f61e807f`) | `arc/json-v2-gemini` |

## Changes

| module | commit | change |
|---|---|---|
| `libs/acl/openai` | `862fbc8` | The reasoning content kept in a delta's extra fields is read with `encoding/json/v2` and `encoding/json`'s semantics instead of sonic; the package no longer imports sonic. |
| `components/model/gemini` | `038fcd9` | Tool results, tool call arguments and Gemini function call arguments are read and written with `encoding/json/v2`, `encoding/json`'s semantics and sonic's escaping (no HTML or U+2028/2029 escapes), map keys sorted (`json.go`). |

Both modules: `go 1.27`; `bytedance/mockey` v1.4.8 (the first that builds the tests on Go 1.27) and `bytedance/sonic` v1.15.4 (needed by the upstream eino these modules require for their own build; in ARC eino is replaced by `mapo80/eino`, without sonic). sonic stays only as a test reference.

## Equivalence

- `libs/acl/openai`: `TestReasoningFromExtraMatchesSonic` — the reasoning read from nine raw values (escapes, HTML characters, empty, null, number, object, array, malformed) is sonic's string, or nothing where sonic failed. Upstream suite green with `-gcflags="all=-N -l"`.
- `components/model/gemini`: `TestMarshalStringMatchesSonic` (argument maps with HTML characters, U+2028/2029, escapes, nested values, float edges; same bytes as sonic's default configuration with sorted keys) and `TestUnmarshalStringMatchesSonic` (11 inputs: duplicates, big numbers, 1e400, non-objects, empty, trailing data, malformed: same maps and error outcomes). Upstream suite green with `-gcflags="all=-N -l"`.

Each call reads or writes a tool payload of a few hundred bytes once per model call; the measures of ARC's turn are in `arc-workspace`.
