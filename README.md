# tea-go

Welcome to the repository of the Golang implementation of the Tealang interpreter runtime. It is NOT compatible to the Python implementation in some aspects. It is heavily work-in-progress, do not use in any production environment.

## How does FizzBuzz look?
```
use io;

operator /?(a, b: int): bool {
    return a % b == 0;
}

for each _, v in map([1..100], func(i) {
    if (i /? 15) {
        return "FizzBuzz";
    } else if (i /? 3) {
        return "Fizz";
    } else if (i /? 5) {
        return "Buzz";
    } else {
        return i: string;
    }
} {
    io.println(v);
}
```

## Syntax examples
These examples pitch an idea of the general feel of the language, where its strength and weaknesses are.

### Declaring a constant
#### Shortform
```
let x = 42;
```
#### Implicit typecast
```
let x: int = 3.4 + 8.6;
```
#### Explicit typecast
```
let x = (3.4 `+ 8.6): int;
```
### Declaring a variable
```
var x = 12;
```

### Declaring a function
#### Basic form
```
func piTimes(x: float): float {
    return math.pi * x;
}
```
#### Shortform
```
func piTimes(x) {
    return math.pi * x;
}
```
#### Variable form
```
let piTimes = func(x) { return math.pi * x }
```
### Declaring an operator
```
operator ?(o: object, f: func) {
    if (o != null) {
        f(o);
    }
}
```
### Declaring an object (user defined type)
#### Long form with defaults and custom constructor
```
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
```
object Person {
    let firstName, lastName;
    func print() {
        io.println(firstName, lastName);
    }
}
```
### Declaring an object instance
#### Using new
```
let x = new(Person, "Michael", "Jackson");
```
#### Using the type constructor
```
let x: Person("Michael", "Jackson");
```
### Declaring a list
#### Basic list of values
```
let x = [1, 2, 3, 4];
```
#### Explicit empty list
```
let x = [];
```
#### Implicit empty list
```
let x: list;
```
#### Range (closed)
```
let x = [1..4];
```
#### Range (closed-open)
```
let x = [1..5[;
```
#### Range (open-closed)
```
let x = ]0..4];
```
#### Range (open)
```
let x = ]0..5[;
```
### Declaring a set
#### Set of values
```
let x = {1, 2, 3, 4};
```
#### Implicit set of values
```
let x = {1, 2, 3, 4, 4};
```
#### Explicit empty set of values
```
let x = {};
```
#### Implicit empty set of values
```
let x: set;
```
### Declaring a map
#### Explicit hashtable
```
let x = {"michael" => "jackson"};
```
#### Explicit empty hashtable
```
let x = {=>};
```
#### Implicit empty hashtable
```
let x: map;
```
### Control flow branching
#### Long form
```
if (x != y) {
    io.println(math.abs(x - y));
} else {
    io.println(0);
}
```
#### Short form with else-if
```
if (x > 0) { io.println(-x) }
else if (x < 0) { io.println(x) }
else { io.println(0) }
```
### Control flow looping
#### C-style iteration looping
```
for (var i = 0; i < 5; i++) {
    io.println(i);
}
```
#### Single condition looping
```
for (i < 5) {
    io.println(i);
}
```
#### For-each on lists
```
for each i in [1..5] {
    io.println(i);
}
```
#### For-each on sets
```
// Careful! Unordered set!
for each x in {1, 2, 3, 4, 5} {
    io.println(x);
}
```
#### For-each on maps
```
for each key, value in {"michael" => "jackson"} {
    io.println(key, "=>", value);
}
```
### Comments
#### Single line comments
```
// a single line comment
3 + 4; // after an expression
```
#### Multi line comments
```
/* oh well */
/* why
   not     */
```
