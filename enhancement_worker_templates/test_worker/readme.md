# Test Worker

This project implements a basic evidence enhancement service. This is a long running service that can be contacted via an http call.

When executed, the project will convert the provided evidence into a hex dump, like the following:

```s
54 68 69 73 20 69 73 20 61 20 74 65 73 74    This.is.a.test
```

Note that each line is limited to 16 bytes

This has been artificially limited to images for testing purposes.

This project uses typescript and express
