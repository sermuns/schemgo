# Installation

There are many ways to install `schemgo`. Pick one you feel comfortable with!

## `go`
```shell
go install github.com/sermuns/schemgo
```

## `ubi`
[`ubi`](https://github.com/houseabsolute/ubi) is a lightweight tool that pulls binary releases from GitHub, and is perfect for installing `schemgo`.

```shell
ubi -p sermuns/schemgo -i ~/.local/bin/
```

*With the `-i` flag, you can change install path to wherever you like.*


## `mise`
[`mise`](https://mise.jdx.dev/getting-started.html) is (among other things) a fantastic tool for managing all your tools. It can use `ubi` as its backend:

```shell
mise install ubi:sermuns/schemgo
```

*This will install `schemgo` only in the current directory. To install globally, use `install -g`.*


## `docker`
```shell
docker pull ghcr.io/sermuns/schemgo
```

You can now run `schemgo` using:

```shell
docker run -v .:/app -u $(id -u):$(id -g) ghcr.io/sermuns/schemgo <arguments to schemgo>
```

!!! tip "`schemgo serve` from container"
    To use `schemgo serve` with a docker installation, you need to make sure the container is serving to the host machine. Easiest is to bind to all interfaces. Do this with the serve flag `--bind 0.0.0.0:3000`.

    You can of course choose any port you want, but remember to also open that port when running the image with e.g. `-p 3000:3000`

    The full command for something like this would be:
    ```shell
    docker run -p 3000:3000 -v .:/app -u $(id -u):$(id -g) ghcr.io/sermuns/schemgo serve <input file> --bind 0.0.0.0:3000
    ```

    
