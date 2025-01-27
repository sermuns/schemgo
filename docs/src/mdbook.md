---
hide:
  - toc
---

# `mdBook` preprocessor

The original idea for `schemgo` was conceived as a simple `mdBook` preprocessor!

You can write `schemgo` directly in your book by adding `schemgo` as a preprocessor.

~~~markdown
## Subheading

This just markdown. Talking about the schematic...

```schemgo
battery.right
line.up
dot
```
~~~

## Setup

Add `schemgo` as a preprocessor by adding this line to your configuration:

```toml
[preprocessor.schemgo]
command = "schemgo mdbook"
```
