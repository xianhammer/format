# parse
Parse is a small and simple byte slice based parsing tool library.

## API

### `parse.Integer(b []byte) (i int64, n int)`

<code>Integer()</code> parse from the beginning of the given slice to a non-decimal-digit is met, returning the value found <b>i</b> and the number of digits parsed <b>n</b>.

As no provisions are made to prevent overflow this function keeps reading until a non-digit is met or end-of-slice.

### `parse.Float(b []byte) (f float64, n int)`

<code>Float()</code> parse from the beginning of the given slice to a non-float-conforming byte is met, returning the value found <b>f</b> and the number of bytes parsed <b>n</b>.

Scientific (exponential) notation is supported, thereby matching same expression as regexp `\d*(?:\.\d*)(?:[eE][-+]?\d+)`.

As no provisions are made to prevent overflow this function keeps reading until a non-float byte is met or end-of-slice.


### `parse.NewBuffer(size int) *Buffer`
Creates a preallocated, bounded, byte-based buffer.


#### `Buffer.Push(c byte) (err error)`
<code>Push()</code> a byte on the buffer. Returns ErrOutOfBounds only if buffer was full prior to push.


#### `Buffer.FetchData() (d []byte)`
<code>FetchData()</code> return a slice representing the current (pushed) data, then the push position is set to 0 (start).

#### `Buffer.GetData() (d []byte)`
<code>GetData()</code> return a slice representing the current (pushed) data.

#### `Buffer.Clear()`
<code>Clear()</code> reset the push position to 0 (start).

#### `Buffer.Empty() bool`
<code>Empty()</code> return true if the buffer is empty (push position is 0), false otherwise.

#### `Buffer.Full() bool`
<code>Full()</code> return true if the buffer is full (push position equals inital set size), false otherwise.
