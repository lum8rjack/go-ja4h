# JA4H Fingerprint in Go

An implementation of the [JA4H hash algorithm](https://github.com/FoxIO-LLC/ja4) in Go.

## JA4H_b Issues

Note that this is not a perfect implementation of the algorithm. 

The JA4H_b section will not be correct because the fingerprint should be the truncated SHA256 hash of the request headers in the order they appear.

Since Go stores the headers in a map, it does not keep the ordering as they appeared in the request.

This implementation of the JA4H_b section sorts the headers before hashing to make the fingerprint consistent.

# References

- [JA4+ Network Fingerprinting](https://blog.foxio.io/ja4+-network-fingerprinting)
- [ja4 GitHub](https://github.com/FoxIO-LLC/ja4)
