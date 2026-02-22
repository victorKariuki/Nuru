# Tools as Primitives Audit

**Principle:** "Tools provide capability, not behavior."

- **PRIMITIVE** = Single capability only (read, write, list, convert, one operation). No orchestration, no "do A then B."
- **WORKFLOW** = Encodes business logic or orchestration (combines multiple capabilities, routing, or prescribed flow).

---

## 1. Scope: What Counts as "Tools" in Nuru

Nuru has **no dedicated agent/MCP tool layer**. The surface that would be exposed to an agent (or used as "tools" in a scripting sense) is:

| Layer | Location | Count |
|-------|----------|--------|
| **Global builtins** | `evaluator/builtins.go` | 10 |
| **Module functions** | `module/*.go` (faili, muda, re, jsoni, hisabati, crypto, njia, mtandao, http, url, mfumo, os) | 154 |
| **REPL** | `repl/repl.go` | No extra tools (eval + exit/toka only) |

**Total tools audited:** 164 (builtins + module functions only; type methods on values are not enumerated as separate "tools" for this audit).

---

## 2. Tool Analysis (representative + all workflow)

| Tool | File | Type | Reasoning |
|------|------|------|-----------|
| **Global builtins** | | | |
| jaza | evaluator/builtins.go | PRIMITIVE | Read one line from stdin; single I/O capability. |
| andika | evaluator/builtins.go | PRIMITIVE | Write to stdout; single I/O capability. |
| _andika | evaluator/builtins.go | PRIMITIVE | Return inspected args as string; no orchestration. |
| aina | evaluator/builtins.go | PRIMITIVE | Type inspection only. |
| mfululizo | evaluator/builtins.go | PRIMITIVE | Generate integer sequence; single capability. |
| badilisha | evaluator/builtins.go | PRIMITIVE | Type conversion; single capability. |
| namba | evaluator/builtins.go | PRIMITIVE | Convert to integer. |
| tungo | evaluator/builtins.go | PRIMITIVE | Convert to string. |
| seta | evaluator/builtins.go | PRIMITIVE | Create set from args/array. |
| jozi | evaluator/builtins.go | PRIMITIVE | Create tuple. |
| **faili (fs)** | | | |
| soma, andika, ongeza, fungua, futa, fanya, ipo, orodha, tengenezaSarafa, futaSarafa, hali, niSarafa, niFaili, ruhusu, mmiliki, badilisha, kiungo, somaKiungo, funga | module/fs.go | PRIMITIVE | Each is a single FS capability (read, write, append, open, remove, create, exists, readdir, mkdir, rmdir, stat, isDir, isFile, chmod, chown, rename, symlink, readlink, close). |
| **muda (time)** | | | |
| hasahivi, lala, tangu, leo, baada_ya, tofauti, ongeza, siku | module/time.go | PRIMITIVE | Single time operations (now, sleep, since, today, after, diff, add, day helper). |
| **re (regex)** | | | |
| linganisha, tafuta, tafutaZote, vikundi, badilisha, gawa, tayari | module/regex.go | PRIMITIVE | Match, find, find all, groups, replace, split, compile; one capability each. |
| **jsoni** | | | |
| dikodi | module/json.go | PRIMITIVE | Decode JSON string to object; single capability. |
| enkodi | module/json.go | PRIMITIVE | Encode object to JSON string. |
| **soma** | module/json.go | **WORKFLOW** | Reads file then decodes (readFile + decode). Combines I/O + parse. |
| **hifadhi** | module/json.go | **WORKFLOW** | Encodes then writes file (encode + writeFile). Combines serialize + I/O. |
| pendeza, msailiaji, msailiaji_bora | module/json.go | PRIMITIVE | Pretty-print, encoder constructors; single capability. |
| **hisabati** | | | |
| PI, e, phi, ln10, ln2, log10e, log2e, sqrt1_2, sqrt2, sqrt3, sqrt5, EPSILON, abs, sign, ceil, floor, sqrt, cbrt, root, hypot, random, factorial, round, max, min, exp, expm1, log, log10, log2, log1p, cos, sin, tan, acos, asin, atan, cosh, sinh, tanh, acosh, asinh, atanh, atan2, namba_kubwa | module/hisabati.go | PRIMITIVE | Single math/constant operations. |
| **crypto** | | | |
| md5, sha1, sha256, sha512, hmac_sha256, hmac_sha512, bahatiNasibu_*, base64_encode/decode, kodeBase64, katuaBase64, hex_encode/decode, *_faili, pbkdf2_sha256 | module/crypto.go | PRIMITIVE | Single hash/encode/decode/random operations. |
| **njia (path)** | | | |
| jina, kigawaji, sarafa, ext, umbiza, niKamili, unganisha, sawazisha, changanua, husika, tatua, kitenga, posix, win32, glob | module/path.go | PRIMITIVE | Path inspection, join, normalize, parse, resolve, relative, separator, platform objects, glob; one capability each. |
| **mtandao (net)** | | | |
| peruzi, tuma | module/net.go | PRIMITIVE | Single HTTP GET/POST request; return body. |
| **http** | | | |
| pata, kichwa, tuma, weka, futa, bandika, chaguzi | module/http.go | PRIMITIVE | Single HTTP method calls (delegate to client). |
| undaMteja, undaOmbi, undaServer | module/http.go | PRIMITIVE | Factory/create; single capability. |
| **sajiliNjia** | module/http.go | **WORKFLOW** | Registers route + handler (method + path → call Nuru function). Encodes routing/orchestration. |
| Server.sikiliza, Server.funga | module/http.go | PRIMITIVE | Bind/listen and close; capability only. |
| **url** | | | |
| changanua, URL, tengeneza, tatua, URLSearchParams, kimbiaNjia, tatuaNjia, kimbiaHoja, tatuaHoja, kamusiKwaHoja, hojaKwaKamusi | module/url.go | PRIMITIVE | Parse, build, resolve, escape/unescape, query ↔ dict; single capabilities. |
| **mfumo** | | | |
| safishaMemori, takwimuMemori, takwimuMemoriKwa, kumbukumbaDhaifu | module/mfumo.go | PRIMITIVE | GC, memory stats, weak ref; single capability each. |
| **os** | | | |
| toka, kimbiza | module/os.go | PRIMITIVE | Exit process; run external command. |

---

## 3. Score

**161 out of 164** tools are proper primitives.

**98.2%** (161/164)

---

## 4. Problematic Tools (Workflow)

| Tool | File | Issue |
|------|------|--------|
| **jsoni.soma** | module/json.go | Combines file read + JSON decode. Behavior: "load JSON from path" instead of separate read + decode. |
| **jsoni.hifadhi** | module/json.go | Combines JSON encode + file write. Behavior: "save object to JSON file" instead of encode + write. |
| **http.sajiliNjia** | module/http.go | Registers route and handler together. Encodes routing/orchestration ("when method+path, run this"). |

---

## 5. Recommendations

1. **Keep primitives as-is**  
   Global builtins and the vast majority of module functions are capability-only and align with "tools as primitives."

2. **Optional decomposition for jsoni**  
   - Keep `soma` and `hifadhi` as conveniences for script authors.  
   - Document that primitive-style usage is: `faili.soma(path)` + `jsoni.dikodi(content)` and `jsoni.enkodi(obj)` + `faili.andika(path, content)`.  
   - If exposing Nuru to an agent tool layer, consider exposing `dikodi`/`enkodi` plus `faili.soma`/`faili.andika` as the primitive tools and treating `soma`/`hifadhi` as derived/convenience.

3. **http.sajiliNjia**  
   - This is inherently workflow (routing + handler registration). For an agent, prefer exposing low-level primitives (e.g. "create server", "register route", "start listening") if they exist, or document that `sajiliNjia` is the single orchestration primitive for HTTP routing.  
   - No change required for the language design; just classify it as the one routing/orchestration tool when documenting "tools as primitives."

4. **Policy / docs**  
   - Add a short note in `docs/POLICY.md` (or a linked doc) that new stdlib APIs should prefer primitive capabilities (single read/write/convert/op) and keep multi-step workflows (read+parse, encode+write, route+handler) as explicit conveniences or document their composition from primitives.

---

*Audit date: 2025-02-22. Nuru codebase; no MCP/agent-specific tool definitions found.*
