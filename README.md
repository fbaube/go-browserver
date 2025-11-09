**General:**
* This is a fork of [**lhtmx**](https://github.com/elijahmorg/lhtmx)
*(no license)*
* [**Article**](https://elijahm.com/posts/local_first_htmx_part2/)
* [** HN discussion**](https://news.ycombinator.com/item?id=45853536) *(2025.11)*
* [**Video**](https://www.youtube.com/watch?v=O2RB_8ircdE) about
  [the underlying **plumbing**](https://github.com/nlepage/go-wasm-http-server)
  `HTML<=>JS<=>Go` (from 2021, a bit outdated)
* [**Slides**](https://nlepage.github.io/go-wasm-http-talk/) for the video

**elem:**
* [**elem**](https://github.com/chasefleming/elem-go) is used for DOM access & mods
* [**counter demo**](https://github.com/chasefleming/elem-go/examples/htmx-counter)
  of elem should be integrated
* [**this article**](https://dev.to/chasefleming/building-a-go-static-site-generator-using-elem-go-3fhh) uses [**elem-ssg**](https://github.com/chasefleming/elem-ssg) to integrate the Goldmark markdown processor with elem 

# Build instructions

```
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./public/
cd ./cmd/web/
GOOS=js GOARCH=wasm go build -o main.wasm . ; cp main.wasm ../../public
cd ../server
go run main.go &
```

_and then_ [**http://localhost:3000**](http://localhost:3000)

## TinyGo

Unfortunately tinygo can't compile most of the net/http packages

## WASM Exec JS

```
cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js .
```


# Todo List App with `elem-go`, `htmx`, `Go labstack/echo`

Based off this [example](https://github.com/chasefleming/elem-go/tree/main/examples/htmx-fiber-todo) but with modifications. I grabbed this example as it looked like it did most of the todo stuff I wanted.

I did not realize what elem-go was - and so in the future I'll probably swap that out for a template based generation but most of the effort of this was towards doing the local first aspect.

