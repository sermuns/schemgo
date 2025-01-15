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

This tool generates electrical cirtcuit schematic diagram from a minimal language, heavily inspired by [Schemdraw](https://schemdraw.readthedocs.io/en/stable/).

Very similar to [Schemdraw](https://schemdraw.readthedocs.io/en/stable/), but also accomplishes the same thing as [circuitikz](https://github.com/circuitikz/circuitikz).


## Quickstart

```sh
go install github.com/sermuns/schemgo
```

`schemgo` is now available in the shell, provided you have `$GOPATH/bin` (`$GOBIN`) in your `PATH`.

## Example usage

Create file `simple.schemgo` containing:
```python
battery.right
line.up
resistor.right
line.down
```

Run:
```sh
schemgo -i simple.schemgo -o simple.svg
```

The circuit diagram is created as `simple.svg`:

<!-- ![simple circuit](docs/simple.webp) -->
<div align="center">
<a href="docs/simple.svg"><img src="docs/simple.webp" alt="simple circuit" align="center" /></a>
</div>


> [!NOTE]
> Only svg output is supported at the moment.

## To-do
- [ ] More components
  - [ ] Capacitor
  - [ ] Inductors
  - [ ] Diodes
- [ ] Push and pop.
- [ ] Labels
- [ ] Element attributes
  - [ ] ID for symbolic reference
- [ ] Syntax highlighting
- [ ] Exporting to pdf, png, jpg, webp
