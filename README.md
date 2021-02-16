Stele
=====

*Disclaimer: Stele is currently _very_ early in development. Many things, including syntax, intended features, and the example in this README are subject to change at basically any time.*

Stele is a statically duck-typed language. It uses structural, implicitly satisfied interfaces to provide the benefits of both a static type system and a duck-typed language.

Features
--------

* Static typing based on patterns of behavior, not data structure.
* Generics.
* Mutability for both variables and functions.
* Simple syntax.
* Easy integration into projects as a library.
* LL(1) parsable.

Example
-------

```stele
import "io"
import "iter"

// rot13 wraps a writer, transforming text written to it via ROT13.
type rot13 {
	var w io.writer
}

// write is required to satisfy io.writer.
func (mut r rot13) write(data array[byte]) result[int] {
	r.w.write(
		iter.ofArray(data)
			.map -> (c) {
				switch {
					(c >= 'a') && (c =< 'z') { c - 'a' + 13 % 26 + 'a' }
					(c >= 'A') && (c =< 'Z') { c - 'A' + 13 % 26 + 'A' }
					else { c }
				}
			}
			.toArray(),
	)
}

// main is the entry point of the standard CLI interpreter.
func mut main() {
	// Declare a variable called encoder and build a rot13 instance to
	// assign to it, assigning io.stdout to its w field.
	:encoder = &rot13 {
		w = io.stdout
	}

	// Write to the encoder.
	io.writeln(encoder, "This is an example.")
}
```

Documentation
-------------

* [Informal, Incomplete Spec](https://github.com/stelelang/stele/blob/master/doc/informal-spec.md)
* [Grammar Specification Used by pgen](https://github.com/stelelang/stele/blob/master/res/grammar.ebnf)

Prior Art
---------

Much of the design of the internals of the language are based on, and, in some cases, directly copied from with possible minor changes, previous work done by the same original author on a scripting language known as [WDTE][wdte].

[wdte]: https://github.com/DeedleFake/wdte
