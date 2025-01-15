<div align="center">
<h1><code>schemgo</code></h1>
<p><em>Dead simple circuit schematic generation.</em></p>
<a href="https://goreportcard.com/report/github.com/sermuns/schemgo"><img alt="goreportcard-badge" src="https://goreportcard.com/badge/github.com/sermuns/schemgo"></a>
<a href="https://www.gnu.org/licenses/gpl-3.0">
<img alt="license-badge" src="https://img.shields.io/badge/License-GPLv3-blue.svg"></a>
</div>

> [!NOTE]
> This tool is under active development and is not currently in a usable state. Stay tuned!

## What is this?

This tool generates electrical circuit schematic diagram from code, such as [Schemdraw](https://schemdraw.readthedocs.io/en/stable/) and [circuitikz](https://github.com/circuitikz/circuitikz).

The language is very minimal and heavily inspired by Schemdraw's.

## Why choose this over Schemdraw/circuitikz?

This ships as a single binary, and is blazingly fast.

**NOT YET IMPLEMENTED:** Is easily included in mdbook as preprocessor.

## Quickstart

```sh
go install github.com/sermuns/schemgo
```

`schemgo` is now available in the shell, provided you have `$GOPATH/bin` (`$GOBIN`) in your `PATH`.

## Example usage

Create a file `simple.schemgo` containing:
```python
battery.right
line.up
resistor.left
line.down
```

Run:
```sh
schemgo -i simple.schemgo -o simple.svg
```

The circuit diagram is created as `simple.svg`:

<div align="center">
<a href="docs/simple.svg"><img src="docs/simple.webp" alt="simple circuit" align="center" /></a>
</div>

> [!NOTE]
> Only svg output is supported at the moment.

## To-do
- [ ] Commenting with `//` and `/* */`
- [ ] `@set` statement to change global defaults (stroke width, padding, color)
- [ ] Components
  - [x] Wire
  - [x] Resistor
  - [x] Battery
  - [ ] Capacitor
  - [ ] Inductors
  - [ ] Diodes
- [ ] Push and pop
- [ ] Labels
- [ ] Element attributes
  - [ ] ID for symbolic reference
- [ ] Syntax highlighting
- [ ] Exporting to pdf, png, jpg, webp
- [ ] mdBook preprocessor
