# tea-go

Welcome to the repository of *tea-go*, the Go implementation of the Tealang runtime. **It is NOT compatible to the Python implementation in some aspects due to further language changes. This project is heavily work-in-progress, do not use in any production environment.**

## How does FizzBuzz look?
```tea
use io;

operator /?(a, b: int): bool {
    return a % b == 0;
}

for (var i = 0; i < 100; i++) {
    var a, b: string;
    if (i /? 3) {
        a = "Fizz";
    } else if (i /? 5) {
        b = "Buzz";
    }
    match (var v = a + b) {
    case "" => {
        io.println(i);
    }
    default => {
        io.println(v);
    }
    }
}
```

## Syntax examples
These examples pitch an idea of the general feel of the language, where its strength and weaknesses are.

### Declaring a constant
#### Shortform
```tea
let x = 42;
```
#### Implicit typecast
```tea
let x: int = 3.4 + 8.6;
```
#### Explicit typecast
```tea
let x = (3.4 `+ 8.6): int;
```
### Declaring a variable
```tea
var x = 12;
```

### Declaring a function
#### Basic form
```tea
func piTimes(x: float): float {
    return math.pi * x;
}
```
#### Shortform
```tea
func piTimes(x) {
    return math.pi * x;
}
```
#### Variable form
```tea
let piTimes = func(x) { return math.pi * x }
```
### Declaring an operator
```tea
operator ?(o: object, f: func) {
    if (o != null) {
        f(o);
    }
}
```
### Declaring an object (user defined type)
#### Long form with defaults and custom constructor
```tea
object Person {
    let firstName = "default first name";
    let lastName = "default last name";
    func print() {
        io.printf("%s %s\n", firstName, lastName);
    }
    new(first, last: string) {
        firstName = first;
        lastName = last;
    }
}
```
#### Short form
```tea
object Person {
    let firstName, lastName;
    func print() {
        io.println(firstName, lastName);
    }
}
```
### Declaring an object instance
#### Using new
```tea
let x = new(Person, "Michael", "Jackson");
```
#### Using the type constructor
```tea
let x: Person("Michael", "Jackson");
```
### Declaring a list
#### Basic list of values
```tea
let x = [1, 2, 3, 4];
```
#### Explicit empty list
```tea
let x = [];
```
#### Implicit empty list
```tea
let x: list;
```
#### Range (closed)
```tea
let x = [1..4];
```
#### Range (closed-open)
```tea
let x = [1..5[;
```
#### Range (open-closed)
```tea
let x = ]0..4];
```
#### Range (open)
```tea
let x = ]0..5[;
```
### Declaring a set
#### Set of values
```tea
let x = {1, 2, 3, 4};
```
#### Implicit set of values
```tea
let x = {1, 2, 3, 4, 4};
```
#### Explicit empty set of values
```tea
let x = {};
```
#### Implicit empty set of values
```tea
let x: set;
```
### Declaring a map
#### Explicit hashtable
```tea
let x = {"michael" => "jackson"};
```
#### Explicit empty hashtable
```tea
let x = {=>};
```
#### Implicit empty hashtable
```tea
let x: map;
```
### Control flow branching
#### Long form
```tea
if (x != y) {
    io.println(math.abs(x - y));
} else {
    io.println(0);
}
```
#### Short form with else-if
```tea
if (x > 0) { io.println(-x) }
else if (x < 0) { io.println(x) }
else { io.println(0) }
```
### Control flow looping
#### C-style iteration looping
```tea
for (var i = 0; i < 5; i++) {
    io.println(i);
}
```
#### Single condition looping
```tea
for (i < 5) {
    io.println(i);
}
```
#### For-each on lists
```tea
for each i in [1..5] {
    io.println(i);
}
```
#### For-each on sets
```tea
// Careful! Unordered set!
for each x in {1, 2, 3, 4, 5} {
    io.println(x);
}
```
#### For-each on maps
```tea
for each key, value in {"michael" => "jackson"} {
    io.println(key, "=>", value);
}
```
### Comments
#### Single line comments
```tea
// a single line comment
3 + 4; // after an expression
```
#### Multi line comments
```tea
/* oh well */
/* why
   not     */
```
