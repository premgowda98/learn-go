# HTML Parser

[Resource](https://www.zenrows.com/blog/golang-html-parser#parse-html-using-the-tokenizer-api)

1. Install package from go `go get -u golang.org/x/net`
2. This has 2 main API's *Tokenizer* and *Node* parsing API
3. Node parsing API is more high-level API
    1. `html.Parse` takes any object which adheres to *io.Reader* interface