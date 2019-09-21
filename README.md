# Sifweb

This is a testing repository for running GoLang in the browser using wasm,
and specifically loading a Singularity image (SIF) header.

**under development**

## Docker

First, build the container. 

```bash
$ docker build -t vanessa/sifweb .
```

It will install [emscripten](https://emscripten.org/docs/getting_started/FAQ.html),
add the source code to the repository, and compile to wasm. You can then
run the container and expose port 80 to see the compiled interface:

```bash
$ docker run -it --rm -p 80:80 vanessa/sifweb 
``` 
