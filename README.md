`go-ristretto`
--------------

**Work in progress**

Many cryptographic schemes need a group of prime order.  Popular and
efficient elliptic curves like (Edwards25519 of `ed25519` fame) are
rarely of prime order.  There is, however, a convenient method
to construct a prime order group from such curves, using a method
called [Ristretto](https://ristretto.group) proposed by Mike Hamburg.

This is a pure Go implementation of the group operations on the
Ristretto prime-order group built from Edwards25519.
