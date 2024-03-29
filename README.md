Stele
=====

*Disclaimer: Stele is currently _very_ early in development. Many things, including syntax, intended features, and the example in this README are subject to change at basically any time.*

Stele is a statically duck-typed language. It uses structural, implicitly satisfied interfaces to provide the benefits of both a static type system and duck typing.

Features
--------

* Static typing based on patterns of behavior, not data structure.
* Generics.
* Simple syntax.
* Easy integration into projects as a library.

Example
-------

```stele
import "io"
import "iter"

func ascii_rot(c! Int) Int {
	switch {
		(c >= 'a') && (c =< 'z') { c - 'a' + 13 % 26 + 'a' }
		(c >= 'A') && (c =< 'Z') { c - 'A' + 13 % 26 + 'A' }
		else { c }
	}
}

// rot13 wraps a writer, transforming text written to it via ROT13.
func rot13(w io.Writer) io.Writer {
	io.Writer {
		func (_) write(data! io.Bytes) Result[Int] {
			w.write(
				iters.of_array(data)
					|> iters.map(ascii_rot)
					|> io.bytes_from_iter(),
			)
		}
	}
}

// main is the entry point of the standard CLI interpreter.
func main() {
	let encoder! = rot13(io.stdout())
	io.writeln(encoder, "This is an example.")
}
```

Documentation
-------------

* [Informal, Incomplete, and often Inaccurate Spec](https://github.com/stelelang/stele/blob/master/doc/informal-spec.md)
