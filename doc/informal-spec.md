Stele
====

Stele is a structurally-typed programming language. It has a type system based entirely around implicitly satisfied interfaces. This allows the language to provide many of the guarantees of static typing _and_ the flexibility of duck typing.

Keywords
--------

```
func var mut priv
type oneof
return continue break
if switch for
```

Predefined Identifiers
----------------------

```
int uint byte float bigint bigfloat
anyint anyfloat numeric
string array
bool true false
any unit
```

Syntax
------

Basic syntax is Go-like, but there are a few differences. Each file is a single package, and the top-level of a file is a list of declarations and import statements. All statements end with a semicolon, but the scanner automatically inserts semicolons into the token stream whenever it sees a newline _unless_ the last token seen was a comma (`,`), or the next token after the newline is a `.`.

A full example package might look like

```stele
import "io"

func mut main() {
	io.stdout.writeln("This is an example.")
}
```

Comments are the typical C-oid language syntax, meaning that a line comment starts with `//` and a multi-line comment is `/* this is a comment */`. As a special case, if the very first token seen is `#`, this also counts as a comment. This allows a shebang line to be inserted at the beginning of scripts.

Variables
---------

Variables are declared using the `var` keyword, such as

```stele
// Declares a variable named example of type int and puts a 3 in it.
var example int = 3
```

Variables, like all declarations, default to being public. A variable, or other declaration, may be made private by prepending the keyword `priv` to it. Privately declared variables, functions, methods, types, and fields are only available inside the package that they are defined in.

Variables are, by default, immutable. To make a variable mutable, use the `mut` keyword:

```stele
var mut example float = 2.3
// ...
example = 5.2
```

Immutable variables must have a value assigned to them when they are declared. Global variables are allowed, but they must be immutable. If an immutable variable is assigned to a mutable variable or vice versa, it is copied, thus preventing mutation of immutable variables even by references.

A single declaration may define as many variables as it likes, separated by commas. Variable assignments of this kind may also be used to destructure a tuple:

```stele
var a, b = someTuple
```

When defining multiple variables, types for each may optionally be listed after the variable name. If multiple variables in a row have the same type, the type may be omitted on all but the last one:

```stele
// a is an int, while b and c are strings.
var a int, b, c string = // ...
```

Note that the use of any types removes type inference from any variables declared to the left of that type usage. In other words, in the above example, `a`, `b`, and `c` are all manually typed, and thus have no inference applied to them from the assignment. This includes `b`, despite it not having a type used directly after its name is declared.

As a shorthand, a local variable may be declared by omitting the `var` keyword and prepending a `:` to the name of the variable in an assignment. Multiple variables may be declared and assigned this way, and each variable prepended with a colon is considered to be a declaration.

```stele
:newVariable = "a new immutable variable"
```

A mutable variable may be declared in this way by inserting the keyword `mut` between the `:` and the name of the variable:

```stele
:mut newVariable = "a new mutable variable"
```

Finally, multiple variables may declared or assigned by mixing these with several identifiers separated by commas:

```stele
existingVariable, :newImmutableVariable, :mut newMutableVariable = ("the number", "of elements", "must match")
```

Due to potential ambiguities, numeric literals in Stele have no type from the point of view of a combined declaration and assignment. For this reason, all declarations that assign a literal of a numeric type must be manually typed:

```stele
// This is an error.
var example = 3

// This works.
var example int = 3
```

It is illegal to declare a new variable of the same name as another in the same scope as it, but a variable with the same name _may_ be declared in a subscope, shadowing the other variable and making it inaccessible in that scope. A variable declaration occurs _after_ evaluating the right-hand side of the assignment, meaning that a variable may be set to some variant of the value of the variable that it is shadowing.

Functions
---------

Functions are declared using the `func` keyword. For example, the following defines a function named `example` that takes an int and returns it multiplied by two:

```stele
func example(v int) int { v * 2 }
```

If the body of a function is a single expression, as above, the result of that expression is automatically returned. Otherwise a function must use the `return` keyword to specify the value to be returned:

```stele
func example(v int) int {
	var doubled = v * 2
	return doubled
}
```

Function arguments are comma-separated. If multiple arguments in a row are of the same type, the type may be omitted on all but the last one:

```stele
// Both a and b are ints.
func example(a, b int) int { /* ... */ }
```

Function arguments may be either mutable or immutable. Mutable arguments are passed by reference, allowing mutation of fields to affect data outside of the method:

```stele
func example(mut v someType) {
	// The instance of someType passed to the function, if also mutable, reflects the change made here.
	v.someIntField = 3
}
```

Whether immutable or not, assigning to the local variable that the function argument defines is not allowed. In other words, the statement `v = &someType {}` is illegal in the above example function. In this way, mutable function arguments differ from regular mutable variables.

Functions can also have receivers. These are defined, much like in Go, via a parenthesized argument before the name of the function:

```stele
func (v int) example(v2 int) int { (v + v2) * 2 }
```

Functions may not be overloaded or shadowed by other functions in the same package, but two methods with different receiver types may have the same name. It is illegal to define a method for a type not declared in the same package as the method.

Functions can also define type parameters. These are defined before anything else, right after the `func` keyword, and are essentially like types that are scoped to the definition. A type parameter has a constraint on it that defines what types are allowed. This constraint may be any type. The predefined type `any` is a special type that is matched by any type.

```stele
// This doesn't actually do anything useful, but it shows the syntax.
func [T any] (v T) example() T { v }
```

If a method has a generic receiver type, it is attached _only_ to applicable types declared inside the same package as it, and it can be overridden by another method with a non-generic receiver type. It is illegal for two methods with the same name and with generic receivers to apply to the same type.

Functions are, like variables, immutable by default. In the case of functions, what this means is that the function is pure in the sense that it can only perform operations that affect state that is explicitly passed to it. It is an error to perform a mutable operation inside of an immutable function, such as assigning to a variable or calling a mutable function. Assigning to a field of a mutable argument is allowed in an immutable function, but calling a mutable method of an argument is not.

To declare a function to be mutable, simply use the `mut` keyword before the function name, as in a variable declaration:

```stele
func mut example() {
	someMutableVariableFromAParentScope = 5
}
```

Finally, the receiver, like any other function argument, must be declared mutable as well if they are going to be mutated:

```stele
func (mut e example) doubleInternalValue() {
	e.val *= 2
}
```

Types
-----

Types are declared via the `type` keyword:

```stele
type example {
	var val int
	func mut print()
}
```

Types define a set of functionality that is available for a value via a list of fields and methods. The above defines a type named `example` which has a `int` field named `val` and a mutable method named `print`. Given the above, a variable defined as

```stele
var mut e example
```

may store any type which has a field named `val` of type `int` and a method named `print` with no arguments and a `unit` return type and, consequently, may have that field and method accessed via that variable.

A type definition may have type parameters, much like a function. If a type specifies any type parameters, it must start the list with an unconstrained type parameter that is used to indicate its own underlying type. For example, to declare a type which has an add method that takes another value of its own type and returns a new value of that type:

```stele
type [T] example {
	func add(T) T
}
```

Other, normal type parameters may be added after the first one. They behave the same as type parameters do for function definitions:

```stele
type [T, E any] list {
	var prev, next option[list[T, E]]
	var val E
}
```

Types may also be embedded into a type definition by simply listing them. If they are, their own set of functionality is replicated added to the set of functionality expected by the type that they are embedded in:

```stele
// example is satisfied by any type that satisfies both otherType and thirdType.
type example {
	otherType
	thirdType
}
```

If a type definition only has a single entry in it, the `{}` around the body may be omitted:

```stele
// example is satisfied by anything that satisfies int.
type example int
```

#### Method Scope

Methods are scoped based on the type of the expression used as the receiver's value. This means that a given method set can be applied, via a conversion, to any value that satisfies an interface. For example,

```stele
type [T] example {
	int
}

func (e example) mut print() {
	// ...
}

// ...

var v int = 3
example(v).print()
// Note that v doesn't need to be mutable, because while the function is mutable, it performs no mutable operations on its receiver. This is why the double syntax is required.
```

In this way, any type may be made to satisfy another interface if the only difference between the two is the methods that each declares.

#### Oneof Types

Oneof types are types that are satisfied by at most one other type's set of functionality defined by a list created at the time of their declaration. They are the Stele equivalent of sum types in many other languages. This list is inserted via the `oneof` keyword:

```stele
// example is satisfied by anything that satisfies int or anything that satisfies string
type example {
	oneof {
		int
		string
	}
}

// Alternative syntax:
type example oneof {
	int
	string
}
```

If a type definition contains a `oneof` list, it may not contain anything else except for embedding of other oneof types. A type may only include one `oneof` list syntactically, though there may be more than one via the embedding process. All of the `oneof` lists are considered to be the same as one long list.

Because a oneof type must only be satisfied by a single entry in the list, all of the entries in the list must match different sets of types. Any intersection in the functionalities of all of the listed types, however, is available to use via the oneof type:

```stele
type example oneof {
	type { func a(); func b() }
	type { func a(); func c() }
}
// ...
var e example = somethingElse
e.a() // Legal because all members have it.
e.b() // Illegal.
```

#### Tuples

Tuples are essentially ordered groupings of values of varying types. They are, in a way, a type of `struct`, but with indexed fields instead of named ones. In Stele, a tuple type is defined as a comma-separated list of at least two other types in parentheses:

```stele
// e is a tuple containing a string and an int, in that order.
var mut e (string, int)
```

To define a named tuple type, simply embed a tuple type definition into the type's declaration:

```stele
type example {
	(string, int)
}
```

As above, the surrounding `{}` may be omitted if there is only one thing in the definition:

```stele
type example (string, int)
```

Only a single tuple type may be embedded into a type definition in order to avoid ambiguities, and, for similar reasons, an array may not be embedded alongside a tuple. However, a tuple may be embedded in a type with fields and methods:

```stele
type example {
	(string, int)
	val name string
}
```

To access elements of a tuple without destructuring, simply access them like array indices:

```stele
var mut e (string, int)
e[0] // string
e[1] // int
```

#### Anonymous Types

Anywhere where a type name can be used, an anonymous type may be used instead. Anonymous types define functionality without being directly reusable. An anonymous type is used similarly to a type declaration, but is not allowed to have any generic parameters besides the self-referential special parameter:

```stele
// Defines a variable named e of an anonymous type containing a single field, v, of e's own underlying type.
var mut e type [T] {
	val v T
}
```

#### Type Assertions

A type assertion allows the value of an expression to be checked for certain underlying functionality at runtime, thus exposing that functionality:

```stele
// Checks if someVariable has the functionality of example available.
if someVariable.(example) {
	// In here, someVariable may be treated as being of type example.
}
```

For more information, see the Type Assertion subsection of the Flow Control section below.

As Stele doesn't really have underlying types, per se, assertions are based on the chain of types that a value has been assigned to. In other words, each time a value is assigned to a variable that has a type different from its current type, the new type is added to a chain of types that track what types the value is associated with. When an assertion is performed, this chain is searched in order from the most recent entry to the original for any type which can be assigned to the type being asserted to. If multiple types match, the first one found is chosen. The type of the returned value is then essentially what it was at that point in the chain, meaning that another assertion may be performed to rewind further if necessary.

#### Type Conversions

Types may be converted from one type to another by calling the name of the type like a function:

```stele
// Converts someValue to an int and sets e to it.
var e = int(someValue)
```

These are similar to type assertions, but they are performed at compile-time, rather than runtime. As such, they always succeed, but they can only convert if the immediate type of the expression that they are passed can be treated as the type being converted to. They are primarily useful for attaching methods to a value by converting it to a type that has those methods assigned to it.

It is exactly equivalent to convert a type or to assign the type to a variable of the converted type:

```stele
var a A = // ...
var b = B(a)
// The above and the below are exactly equivalent:
var a A = // ...
var b B = a
```

Predefined Types
----------------

There are several built-in types in several different categories. All user-defined types are based on these and on types introduced by code written in another language.

All types have zero values. Zero values are the values that a mutable variable or a struct field defaults to if not specified. These are listed for each type in their own section below.

### Numbers

The predefined number types are

* `int` and `uint` for signed and unsigned integers.
* `float` for floating point numbers.
* `bigint` and `bigfloat` for arbitrary precision integer and floating point numbers.
* `byte` for single, unsigned bytes.

All number types have a zero value of `0`.

### Arrays

`array`s are the only predefined collection type, but there are also `iterator`s as a built-in interface. There is no special array syntax for indicating an array type. Instead, the `array` interface is defined as a generic type, meaning that to specify an array of `int`s, simply use `array[int]` as the type.

Arrays are all variable length, and provide methods for manipulating the amount of their contents.

Array indices are accessed via `[]`, such as `someArray[3]` and `someMutableArray[1] = 2`. Arrays are zero-indexed.

The zero value of an array is an empty array, which is an array of length zero.

### Strings

While the predefined `string` type is obviously a string, `byte` actually also satisfies the `string` interface, allowing them to be used as single-character strings. `string` does _not_ satisfy `array[byte]`, however, as strings are internally immutable, unlike arrays. They do, however, provide methods for converting between each other.

The zero value of a string is a string of length zero.

### Booleans

Booleans may be either `true` or `false`, and nothing else. That's about it.

The zero value of booleans is `false`.

### `any` and `unit`

`any` is a type that matches anything, meaning that anything can be stored in a variable of type `any`. Use of `any` outside of constraints is allowed, but discouraged. As a constraint, it allows any type to be used.

`unit` is pretty much the opposite. It is a type that no matter what only has one valid value, `unit`. It is the implicit return type of functions with no return type declared. It is legal to use `unit` as a constraint, but pointless, as it it would only allow the type `unit` itself to be used.

The zero value of `any` is `unit`, as is, obviously, the zero value of `unit`.

### Functions

Function signatures may be specified as a `->` followed by an optional parenthesized argument type list, followed by a return type:

```stele
// e stores a function that takes two ints and returns an int.
var mut e -> (int, int) int
```

As mutability is part of the type of function values, not just of a variable holding a function, mutability must be specified if a variable might hold a mutable function reference:

```stele
// e stores a mutable function that takes no arguments and returns unit.
var mut e -> mut
```

Note that a function must be mutable in order to call other mutable functions, and therefore any function which takes a mutable function as an argument must also be mutable in order to call that function.

The zero value for a function type is a no-op function that returns the zero value of its return type.

Literals
--------

### Numbers

All numeric symbols are treated as an untyped numeric literal. They are converted automatically to whatever type makes sense. Any situation in which several different types are valid but a single type must be chosen is illegal. For example, `var v int = 3` is valid, but `var v = 3` is not. Similarly,

```stele
func example(v numeric) numeric = v * 2
// ...
example(3)
```

is valid, but

```stele
func [T numeric] example(v T) T = v * 2
// ...
example(3)
```

is not. This is because in the first case, the type can just be chosen arbitrarily, as the user has typed it, while in the second case the type system must actually choose a type for `T`, but the literal `3` does not have enough information to do so.

Character literals are single-quoted, such as `'a'`, or `'あ'`. If the character in the literal is too large to fit into a byte, it is an error to attempt to use it in a place that a byte is required. It is, however, valid to use any character literal with any other numeric type.

### Strings

String literals are double-quoted or delineated with backticks. All string literals may contain newlines as well, meaning that they essentially double as heredocs. Double-quoted literals implement string interpolation via a `$name` and `${expression}` syntax:

```stele
var e = "a simple string"
var e2 = "a $someVariable string with interpolation"
```

Backtick-delineated literals produce a raw string with no interpolation. The number of backticks must be matched on both sides of the string, and, as a result, the string may not be empty:

```stele
var e = ``A string with two backticks.``

var e2 = `
A string with one backtick.
`
```

All string literals produce UTF-8 encoded strings.

### Structs

Though Stele technically doesn't have structs, a type containing only field definitions may act like a struct:

```stele
type example {
	var name string
	var val int
}
// ...
var e = &example{
	name = "An Example"
	val = 3
}
```

Any of the fields may be omitted, and each may only be specified at most once each. If any are omitted, they are set to the zero value of their declared type. For more information, see the section on zero values.

### Tuples

Tuples are instantiated via a set of parentheses with a comma-separated list of values in between:

```stele
type example (string, int)
// ...
var e example = ("An Example", 3)
```

A tuple literal may be typed. This looks like a conversion, but has more than one argument:

```stele
type example (string, int)
// ...
var e = example("An Example", 3)
```

If a type contains an embedded tuple and fields, it may be used by combining the two syntaxes:

```stele
type example {
	(string, int)
	val name string
}
// ...
var e example = &example("something", 3) {
	name = "An Example"
}
```

### Arrays

Array literals are similar to tuple literals, except that they use `[]` instead of `()` and may be any length:

```stele
var e = [3, 2, 5]
```

There is no special syntax for literals of named array types. Instead, simply perform a conversion on a regular array literal.

### `unit`

`unit` is both a type and the only valid value of that type. `unit` is the implicit value of a return with values provided. In other words:

```stele
// The lack of a return type is the same as specifying unit manually.
func example() {
	// Exactly the same as return unit.
	return
}
```

The only types which may have a value of `unit` are `unit` and `any`.

### Functions

A closure may be created using the following syntax:

```stele
:add = -> (a, b) { a + b }
```

The body of a closure is the same as any other function, but the argument list differs in that the types on the arguments are optional. The return type, however, is not optional. If a type is omitted from an argument, it will be inferred from the usage of the of closure in an assignment.

Closures may also be marked as mutable:

```stele
// If a closure has no arguments, the () may be omitted.
-> mut { io.stdout.writeln("As an IO function, this is mutable.") }
```

When a function is called, if its final argument is itself a function type, the closure may be moved outside of the parentheses:

```stele
// Given
func someFunction(a, b int, f -> (int) int) { /* ... */ }
// then
someFunction(1, 2, -> (v) { v + 1 })
// is equivalent to
someFunction(1, 2) -> (v) { v + 1 }
```

If the function being called _only_ takes a single function argument, then the parentheses for the call itself may also be omitted:

```stele
someIterator.forEach -> mut (v) { io.stdout.writeln(v) }
```

Control Flow
------------

### `if`

```stele
if condition {
	// ...
} else if otherCondition {
	// ...
} else {
  // ...
}
```

There are no `()` required around the condition and the body must have `{}` around it.

As in many newer languages, `if`s are expressions, not statements:

```stele
var e int = if condition { 3 } else { 2 }
```

As in functions, the body of each branch of an `if`-`else` chain must contain only a single expression in order to return anything. To help ensure that no mistakes are made in this regard, it is illegal to have a complex branch body that would not return anything and use the `if`-`else` chain as an expression.

The type of the value returned from the `if` is not required to be the same for every branch of the `if`-`else` chain. Instead, the return type is a oneof type of all of the possible return types. If there is no `else` condition and none of the conditions are true, the chain returns `unit`.

In all cases, the condition expression used must satisfy `bool`.

### `switch`

```stele
switch n {
	<= 1 { return n }
	else { return fib(n - 1) + fib(n - 2) }
}
```

Switches, despite being named similarly to C-style switches, function similarly to the `when` expression in Kotlin. A quirk, however, is that they do not match exactly by default. Instead, each case must be prepended with a comparison operator, such as `==`, `<=`, or `!=`.

Similarly to `if` expressions, they are also expressions, meaning that the above example can be rewritten as follows:

```stele
return switch n {
	<= 1 { n }
	else { fib(n - 1) + fib(n - 2) }
}
```

Like in `if` expressions, the bodies of the branches must be a single expression to use this format, and the type of entire expression is a oneof type of all of the possible return types. If none of the cases match and there is no `else` case, the returned value will be `unit`, and the oneof of the expression will include `unit` in its type list.

### Type Assertions

Type assertions are managed via flow control, both via `if` and `switch`.

If the condition of an `if` is of the form `expression.(typeName)`, the condition is considered true if the assertion of `expression` to `typeName` is valid. If `expression` is a single identifier which is an immutable variable, then that identifier may be treated as being of the type that it was asserted to inside of the `if` body.

In a `switch`, the syntax is slightly different, but the idea is the same:

```stele
switch expression {
	.(int) { /* In here, expression, if it is a single immutable identifier, may be treated as an int. */ }
	.(string) { /* And in here it would be a string. */ }
	else { /* expression's type is unaffected in here. */ }
}
```

### `for`

`for` is the only loop keyword in Stele. It functions similarly to in Go in that the format determines the way in which it is used, but the `init; condition; step` format that Go gets from C is not present. In other words:

```stele
for {
	// Loop forever.
}

for condition {
	// Loop while the condition is true.
}
```

There is no built-in equivalent of a `foreach`-style loop. Instead, the standard library provides iterators that offer the functionality via a method, as well as functions for creating iterators that range over values in a way similar to the `C`-style `init; condition; step` format.

Miscellaneous
-------------

### Operators

Operators may not be overloaded. All built-in operators map to predefined methods, however, which allows types to match against operator-like functionality. Operators may themselves be applied to any type for which they can be obviously defined.

All binary operators require that both sides of the operation are of the same type. For example, an `int(3) + int(2)` is valid, but `int(3) + float(2)` is not. The result of a binary math operator is always the input type, while the result of a comparison operator is always a `bool`. The `+` operator, as is standard, also performs string concatenation, but other operators are not defined for strings.

Following the above rules:

```stele
int(3) + int(2) // valid, returns int
int(3) + float(2) // invalid

type example int
example(3) + example(2) // valid, returns example
```
