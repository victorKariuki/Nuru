# Getting Started

**Part 0, unit 0.1** — Environment and first program. Install Nuru, run the REPL, run a script file, and print your first line.

## Learning outcomes

- Install Nuru on your system (Linux, macOS, Windows, or Android/Termux).
- Start the REPL and run a single expression.
- Run a Nuru script from a file.
- Use `andika` to print output.

## Prerequisites

- A terminal (or Termux on Android). For building from source: [Go 1.21+](https://go.dev/dl/).

## Examples

### Install and check

Download the latest binary from the [releases](https://github.com/NuruProgramming/Nuru/releases) page for your platform, or build from source (see the main [README](../../README.md#installation)). Then confirm:

```bash
nuru -v
```

You should see the Nuru version.

### Run the REPL

Start the interactive REPL:

```bash
nuru
```

You’ll see a prompt like `>>> `. Type a Nuru expression and press Enter. The REPL evaluates it and prints the result.

### Your first print

In the REPL, type:

```nuru
andika("Habari")
```

**Output:** `Habari`

`andika` is Nuru’s print builtin. It takes zero or more arguments and prints them, separated by spaces. With no arguments it prints a newline.

### Run a script file

Create a file `habari.nr` with one line:

```nuru
andika("Habari Nuru")
```

Run it:

```bash
nuru habari.nr
```

**Output:** `Habari Nuru`

You can run any `.nr` file this way. Use the REPL for quick experiments and file runs for scripts.

### Print multiple values

`andika` can print several values in one call. They are separated by spaces:

```nuru
andika("Jina", "langu", "ni", "Nuru")
```

**Output:** `Jina langu ni Nuru`

Use `andika()` with no arguments to print a blank line.

## Reference

- [Installation](../../README.md#installation) — Full install steps (Linux, macOS, Windows, Android).
- [Built-in functions](../../repl/docs/en/builtins.md) — `andika`, `jaza`, `aina`, and other builtins.
- [Comments](../../repl/docs/en/comments.md), [Identifiers](../../repl/docs/en/identifiers.md), [Keywords](../../repl/docs/en/keywords.md) — Naming and comments (curriculum 0.1).
- [Language documentation](../../repl/docs/en/README.md) — Full syntax and types.

## Next

**Part 1 — Values, types, and expressions:** [Variables and types](UNITS_SUPPORTED_BY_CODEBASE.md#part-1--values-types-and-expressions) (unit 1.1) — `fanya`, `aina`, operators, numbers, strings, booleans, null. See the [curriculum](UNITS_SUPPORTED_BY_CODEBASE.md) for the full progression.

## Keywords

| Nuru   | Meaning / use        |
|--------|----------------------|
| `andika` | Print to the console (builtin). |
| `nuru`   | Command to run the interpreter or a script. |
| REPL     | Read–eval–print loop; interactive prompt. |

## Run this

Create `habari.nr` as in the example above and run `nuru habari.nr`. For more runnable examples, see the [examples](../../examples/) folder (e.g. [example.nr](../../examples/example.nr)).
