go-ristretto
============

Many cryptographic schemes need a group of prime order.  Popular and
efficient elliptic curves like (Edwards25519 of `ed25519` fame) are
rarely of prime order.  There is, however, a convenient method
to construct a prime order group from such curves, using a method
called [Ristretto](https://ristretto.group) proposed by Mike Hamburg.

This is a pure Go implementation of the group operations on the
Ristretto prime-order group built from Edwards25519.
Documentation is on [godoc](https://godoc.org/github.com/bwesterb/go-ristretto).


References
----------

The curve and Ristretto implementation is based on
[Peter Schwabe](https://cryptojedi.org/peter/index.shtml)'s unpublished PandA
library — see `cref/cref.c`.  The field operations borrow
from [Adam Langley](https://www.imperialviolet.org)'s
[ed25519](http://github.com/agl/ed25519).

### other platforms
* [Rust](https://github.com/dalek-cryptography/curve25519-dalek)
