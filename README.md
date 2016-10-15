# Deduper

This is a fork from github.com/mauidude/deduper, that allows you to find near duplicate or similar documents given another document.

The HTTP server, go-raft (as a cluster with other nodes and provide high-availability) etc have been removed, leaving only the duplicate finding
part, which now can be used as a generic document dedupe library.

## Testing

```
$ deduper
Usage of deduper:
  -bands int
        Number of bands (default 100)
  -debug
        Enable debug logging
  -hashes int
        Number of hashes to use (default 2)
  -leader string
        The HTTP host and port of the leader
  -shingles int
        Number of shingles (default 2)
  -threshold float
        Threshold (default 0.5)
Data string argument required

$ deduper "hello world foo baz"
[{"id":"p1","similarity":1}]

$ deduper "world foo baz"
[]

$ deduper "unrelated"
[]

$ deduper "entirely unrelated"
[{"id":"p3","similarity":1}]

$ deduper "entire unrelate"
[]

$ deduper -threshold 0.00001 "unrelated"
+1.000000e-005[]

$ deduper "foo qux bar zomg"
[]

$ deduper "foo qux bar zomg" -threshold 0.0000000001
[]

$ deduper -threshold 0.0000000001 "foo qux bar zomg world"
+1.000000e-010[]

$ deduper "foo qux bar zomg world goodbye"
[]

$ deduper "goodbye world foo qux bar zomg"
[{"id":"p2","similarity":1}]

```

NB, I believe there is something wrong with above test results, refer to
https://github.com/mauidude/deduper/issues/2 for details.
