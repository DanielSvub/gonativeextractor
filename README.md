# NativeExtractor module for Golang
This is official Golang binding for the [NativeExtractor](https://github.com/SpongeData-cz/nativeextractor) project.

<p align="center"><img src="https://raw.githubusercontent.com/SpongeData-cz/nativeextractor/main/logo.svg" width="400" /></p>

# Installation
## Requirements
* Golang,
* [NativeExtractor](https://github.com/SpongeData-cz/nativeextractor) installation.

# Usage
## Creating a Extractor
The following parameters are required to create an Extractor:
* **batch** - number of logical symbols to be analyzed in the stream, if **batch** is negative, the default value is set to 2^16,
* **thread** - the number of threads for miners to run on, if **thread** is negative, the default value is set to maximum threads available,
* **flags** - initial flags, more about [Flags](https://github.com/SpongeData-cz/nativeextractor#flags).

The new Extractor must be deallocated with the *Destroy* method after use.

Returns a pointer to a new instance of Extractor.

More about [Extractor](https://github.com/SpongeData-cz/nativeextractor#extractor).

```go
e := gonativeextractor.NewExtractor(-1, -1, 0)
```

## Adding miner


This function allows you to add miners from a Shared Object (.so library).

The following parameters are required to add a miner:
* **sodir** - a path to the Shared Object,
* **symbol** - Shared Object symbol and
* **params** - i.e. what the miner mines, are optional - may be empty array or nil, but if present, has to be terminated with \x00.

The default path to the .so library is set to */usr/lib/nativeextractor_miners*.

Returns an error if a path to a non-existent file is specified or if a non-existent miner name is specified.

More about [Miners](https://github.com/SpongeData-cz/nativeextractor#miners).

```go
err := e.AddMinerSo(gonativeextractor.DEFAULT_MINERS_PATH+"/glob_entities.so", "match_glob", []byte("world\x00"))
```

## Streams
There are two types of streams. One is the FileStream and the other is the BufferStream. 
* When creating a FileStream, the path to the file is required.
* When creating a BufferStream, a byte array terminated with "\x00" is required.

```go
st, err := gonativeextractor.NewFileStream("./fixtures/create_stream.txt")
```
```go
st, err := gonativeextractor.NewBufferStream([]byte("Hello world byte\x00"))
```


The stream needs to be attached to the Extractor.

```go
err = e.SetStream(st)
```

At this point you have an Extractor created, a miner added, and a stream created and attached.

## Flags
It is also possible to explicitly set and unset flags.

An Extractor may have these flags enabled:
* **E_NO_ENCLOSED_OCCURRENCES**
* **E_SORT_RESULTS**

More about [Flags](https://github.com/SpongeData-cz/nativeextractor#flags).

```go
err = e.SetFlags(gonativeextractor.E_SORT_RESULTS)
```
```go
err = e.UnsetFlags(gonativeextractor.E_SORT_RESULTS)
```

## Occurrences and batches
Now you can iterate through individual batches and their individual occurrences. Cycle within a cycle.

For iterate over the batches, you can use the *Eof* function, which returns bool, to check if the stream attached to the Extractor has ended. At the same time, you need to use the *Next* function, which primarily gives the next batch of found entities and returns the first found occurrence and error, if any.

Here comes the inner, second, iteration.
The *Eof* function can be used again, but from the Occurrence class. This checks whether all occurrences have been read.
Again, the *Next* function (again from the Occurrence class) must be used with it. This just moves the pointer to the next occurrence and returns nothing. If it is EOF, it does nothing.

From one occurrence you can get the following:
* **Str** - Creates a string containing found occurrence,
* **Pos** - Casts position of the found occurrence to Go integer type,
* **Upos** - Same as *Pos* but with UTF position,
* **Len** - Casts length of the found occurrence to Go integer type,
* **Ulen** - Same as *Len* but with UTF length,
* **Label** - Casts label of the found occurrence to Go string,
* **Prob** - Casts probability of the found occurrence to Go float type.

```go
for !e.Eof() {
    r, err := e.Next()
    // Q: Ma byt v readme err != nil ??
    if err != nil {
       return err
    }
    for !r.Eof() {
        fmt.Println(r.Str())   // "world"
        fmt.Println(r.Pos())   // 6
        fmt.Println(r.Upos())  // 6
        fmt.Println(r.Len())   // 5
        fmt.Println(r.Ulen())  // 5
        fmt.Println(r.Label()) // "Glob"
        fmt.Println(r.Prob())  // 1

        r.Next()
    }
}

```

## Meta
There is also a *Meta* function that give you meta information about the Extractor such as:
* **Ldpath** - Path to the .so library,
* **Ldsymb** - Miner function name,
* **Meta** - Meta info about miner functions and labels,
* **Params** - Miner parameters,
* **Ldptr** - Pointer to the loaded .so library.

```go
meta := e.Meta()
fmt.Println(meta[0].Params)  // "world"
fmt.Println(meta[0].Meta)    // 1
fmt.Println(meta[0].Meta[0]) // "match_glob"
```

## Destroy
In the end, it is necessary to destroy the Extractor with *Destroy*. This function will unset the stream if it has not been unset, close it, and then destroy the Extractor itself.

```go
err = e.Destroy() 
```

NAH
If you want to explicitly call the *UnsetStream* function, you must first keep a reference to the stream.


If we explicitly want to call the *UnsetStream* function, the *Close* function must be called before.

```go
err = st.Close()
```
```go
e.UnsetStream()
```

<!-- Stream je potreba zavrit pred unsetnutim -->