sockjs-go
=========

Sockjs-go implements server side counterpart for the [SockJS](http://sockjs.org)-client browser library.

Status: **EXPERIMENTAL**, use with caution. Passes 68 of 72 tests.

## Installation

    go get github.com/fzzy/sockjs-go/sockjs


## Documentation and examples

Documentation is available at http://gopkgdoc.appspot.com/pkg/github.com/fzzy/sockjs-go/sockjs.

Alternatively, run godoc:

	godoc -http=:8080

and point your browser to http://localhost:8080/pkg/github.com/fzzy/sockjs-go/sockjs.

Also, look into the `examples` folder for examples.

## HACKING

If you make contributions to the project, please follow the guidelines below:

*  Maximum line-width is 100 characters.
*  Run "gofmt -w -s" for all Go code before pushing your code. 
*  Avoid commenting trivial or otherwise obvious code.
*  Avoid writing fancy ascii-artsy comments. 
*  Write terse code without too much newlines or other non-essential whitespace.


## Copyright and licensing

*Copyright 2013 Juhani Åhman*. 
Unless otherwise noted, the source files are distributed under the
*MIT License* found in the LICENSE file.
