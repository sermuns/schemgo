## `docker`

```bash
docker run -v .:/app -u $(id -u):$(id -g) schemgo build examples/simple.schemgo -o simple.svg
```
