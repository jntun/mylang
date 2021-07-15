<h2>Mylang (jlang)</h2>
<p>This is a personal project language I decided to create as a test of skill and enjoyment.
It is a  Javascript/ECMAScript style language with a mix of go flavorings. It doesn't use any external package dependencies and has a completely hand wrote AST. 

**_Try it on [jntun.com/jlang](https://jntun.com/jlang)!_**
</p>

<h2>_Declarations_</h2>
<h3>Variable</h3>

```go
var test = "hello" + "world: " + pi;
```

<h3>Function</h3>

```go
func area(x, y) {
    return x * y;
}
```

<h3>Arrays</h3>

```go
var vector = [1.0, 1.0, 1.0];
print vector[0];
```

<h3>For</h3>

```go
for var i = 0; i < 10; i = ++i {
	print i;
}
```

<p>or</p>

```go
var i = 10;
for i > 0;  i = --i {
    print i;
}
```

<h3>While</h3>

```javascript
var x = 0;
while x < 50 {
    print x;
    x = x + 1;
}
```