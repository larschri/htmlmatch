# htmlmatch

Library used to match elements and attributes inside HTML. The patterns are
plain HTML with only those elements and attributes that should be compared.
That is, elements and attributes can be omitted if they don't need to match.
The HTML elements must match exactly in both order and structure, but
attributes are not sensitive to order.

## Motivation

Make it easy to write tests that verifies the presence of of essential HTML
elements and attributes. Use it to verify [htmx](https://htmx.org/) attributes
and build confidence that the application keeps working, without comprehensive
browser testing tools such as [chromedp](https://github.com/chromedp/chromedp).

## Text matching

Whitespace are trimmed by default before matching text. It is also possible to
match verbatim by prepending the pattern text by `verbatim:`. Substring
matching is supported by prepending the pattern with `substring:`.

## HTML parsing

https://pkg.go.dev/golang.org/x/net/html can be used to parse HTML, but the
results can be confusing since it follows complicated rules for HTML5. This
library provides a primitive alternative `ParseVerbatim` to make it easy to
parse patterns without surprises.

## Example

```go
func ExampleContainsTree() {
        full := MustParseVerbatim(`<html>
                <body>
                        <div hx-test="/foo" hx-target="#mytarget">Click here</div>
                        <div>lorem ipsum</div>
                        <div>lorem <span id="mytarget">ipsum</span></div>
                </body>
                </html>`)
        pattern := MustParseVerbatim(`
                <div hx-target="#mytarget"></div>
                <span id="mytarget" />`)
        fmt.Println(ContainsTree(full, pattern))
        // Output: true
}
```

