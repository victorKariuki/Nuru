# CRUD Completeness Audit

**Scope:** First-class entities in the Nuru scripting language (object types, builtins, module APIs).  
**Criteria:** For each entity, whether **Create**, **Read**, **Update**, and **Delete** exist in the language or via agent/module tools.

---

## Entity CRUD Analysis

| Entity | Create | Read | Update | Delete | Score | Notes |
|--------|--------|------|--------|--------|-------|--------|
| **ORODHA (Array)** | ✅ Literal `[]`, `mfululizo()` | ✅ Index `[i]`, `idadi`, `yamwisho`, `tafuta`, `kitanzi` | ✅ `arr[i]=v`, `sukuma`, `geuza`, `panga` | ❌ | 3/4 | No "ondoa" by index; `Remove()` exists in Go only |
| **KAMUSI (Dict)** | ✅ Literal `{}` | ✅ `fungua(k)`, `idadi`, `funguo`, `maana`, `vikundi`, `kitanzi` | ✅ `dict[k]=v` (index assign) | ❌ | 3/4 | No key-remove in language; `Delete()` in Go only |
| **SETI (Set)** | ✅ `seta()` builtin | ✅ `idadi`, `ona(x)`, `kitanzi` | ✅ `ongeza` | ✅ `ondoa` | 4/4 | Full CRUD |
| **JOZI (Tuple)** | ✅ `jozi()` builtin | ✅ Index `[i]`, `idadi`, `kitanzi` | N/A | N/A | 3/4 | Immutable by design |
| **FAILI (File)** | ✅ `faili.fungua()`, `faili.fanya()` | ✅ `soma`, `hali`, `isFungwa` | ✅ `andika`, `ongeza`, `tafuta` | ✅ `faili.futa(path)` (fs), `funga` (handle) | 4/4 | Full CRUD |
| **NENO (String)** | ✅ Literal, `tungo()` | ✅ Index, `idadi`, `gawa`, `ina`, etc. | N/A | N/A | 2/4 | Immutable; methods return new strings |
| **NAMBA (Integer)** | ✅ Literal, `namba()` | ✅ Operators, inspect | N/A | N/A | 2/4 | Value type |
| **DESIMALI (Float)** | ✅ Literal, `badilisha(..., "DESIMALI")` | ✅ Operators, inspect | N/A | N/A | 2/4 | Value type |
| **BOOLEAN** | ✅ Literal, `badilisha(..., "BOOLEAN")` | ✅ Operators | N/A | N/A | 2/4 | Value type |
| **MUDA (Time)** | ✅ `muda.sasa()`, `muda.baadaye()`, etc. | ✅ `panga`, `tangu`, `ongeza` | N/A | N/A | 2/4 | Value type |
| **TAREHE (Date)** | ✅ `muda.tarehe()`, `muda.tareheSasa()` | ✅ `panga` | N/A | N/A | 2/4 | Value type |
| **RE_ILIYOTAYARISHWA (CompiledRegex)** | ✅ `re.tayari(pattern)` | ✅ `linganisha`, `tafuta`, `tafutaZote`, `vikundi`, `badilisha`, `gawa` | N/A | N/A | 2/4 | Immutable |
| **BASE64** | ✅ `crypto.kodeBase64()`, `crypto.katuaBase64()`, `.fromData()` | ✅ `kukata`, `data` | N/A | N/A | 2/4 | Value type |
| **NAMBA_KUBWA (BigInteger)** | ✅ `badilisha(..., "NAMBA_KUBWA")` | ✅ Arithmetic (infix), Inspect | N/A | N/A | 2/4 | Value type |
| **BYTE** | ✅ `crypto.katuaBase64(..., byte)` | ✅ Inspect (no methods in evaluator) | N/A | N/A | 2/4 | Value type, minimal API |
| **MODULE** | ✅ `tumia` import | ✅ Call module functions | N/A | N/A | 2/4 | Not mutable |
| **PAKEJI (Package)** | ✅ `pakeji` syntax | ✅ Method call, `andaa` | N/A | N/A | 2/4 | Not mutable |
| **INSTANCE** | ✅ Call Package as constructor | ✅ Method call | Depends on package | N/A | 2/4 | Package-defined |
| **KITANZI (Iterator)** | ✅ `.kitanzi` on Array/Dict/String/Set/Tuple | ✅ `for-in` (Next/Reset) | N/A | N/A | 2/4 | Not mutable |
| **TUPU (Null)** | ✅ Literal / absence | ✅ Inspect | N/A | N/A | 2/4 | Value type |
| **UNDO (Function)** | ✅ Function literal | ✅ Call | N/A | N/A | 2/4 | Not mutable |

---

## Overall Score

- **Entities where full CRUD is expected (mutable collections + File):** 4 → **Array, Dict, Set, File**
- **Entities with full CRUD:** **2** → **Set, File**
- **Overall (mutable-entity) score:** **2/4 = 50%** of such entities have full CRUD.

If counting *all* user-facing first-class types (including value types where Update/Delete are N/A): **20+ entities**; only **Set** and **File** have all four operations; the rest either are immutable/value types (no U/D) or are missing Delete (Array, Dict).

---

## Incomplete Entities

1. **ORODHA (Array)**  
   - **Missing:** Delete (remove element by index).  
   - `object.Array` has `Remove(index int)` in Go but no language-level method (e.g. `ondoa` or `futa` by index).

2. **KAMUSI (Dict)**  
   - **Missing:** Delete (remove key).  
   - `object.Dict` has `Delete(key Object) bool` in Go but no language-level method (e.g. `kamusi.ondoa("key")` or `faili.futa`-style API).

---

## Recommendations

1. **Array – expose Delete**
   - Add a method (e.g. `ondoa` or `ondoaKiashiria`) that takes an index and removes the element at that index, calling the existing `Remove(index int)` and returning the array (or the removed element).  
   - Alternatively add `ondoaThamani(x)` to remove first occurrence of a value.

2. **Dict – expose Delete**
   - Add a method (e.g. `ondoa(ufunguo)` or `futa(ufunguo)`) that removes the key from the dict, using the existing `Delete(key Object)` and returning the dict or a boolean.  
   - Keeps parity with Set (`ondoa`) and fs (`futa`).

3. **Optional: Byte**
   - Byte has no methods in `evaluator/method.go`. If Byte is first-class, consider at least Read (e.g. indexing or conversion to string/array) and possibly creation helpers beyond `crypto.katuaBase64(..., byte)`.

4. **Document value types**
   - In language/docs, state that NENO, JOZI, MUDA, TAREHE, BASE64, NAMBA_KUBWA, etc. are immutable/value types, so "Update" and "Delete" are intentionally N/A (create/read only).

---

*Audit based on: `object/*.go`, `evaluator/builtins.go`, `evaluator/method.go`, `evaluator/evaluator.go`, `evaluator/index.go`, `module/fs.go`, `module/regex.go`, `module/time.go`, `module/crypto.go`.*
