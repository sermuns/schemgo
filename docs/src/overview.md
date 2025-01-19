# Introduction

Perhaps one of the best ways to demonstrate `schemgo` is to compare it to its peers.

Let us draw a simple circuit diagram with a battery and a resistor:

![](media/simple.svg)

**`schemgo`:**
<!-- Abusing syntax highlighting... this is NOT haskell. -->
```haskell
battery.right
line.up
resistor.left
line.down
```

**`schemdraw:`**
```python
import schemdraw
from schemdraw import elements as e

with schemdraw.Drawing():
    e.Battery().right()
    e.Line().up()
    e.Resistor().left()
    e.Line().down()
```

**`circuitikz:`**
```latex
\begin{circuitikz}
    \draw (0,0) to[battery] (2,0)
                to[resistor] (2,-2)
                to[short] (0,-2)
                to[short] (0,0);
\end{circuitikz}
```
