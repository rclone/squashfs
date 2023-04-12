# squashfs

[![PkgGoDev](https://pkg.go.dev/badge/github.com/CalebQ42/squashfs)](https://pkg.go.dev/github.com/CalebQ42/squashfs) [![Go Report Card](https://goreportcard.com/badge/github.com/CalebQ42/squashfs)](https://goreportcard.com/report/github.com/CalebQ42/squashfs)

A PURE Go library to read squashfs. There is currently no plans to add archive creation support as it will almost always be better to just call `mksquashfs`. I could see some possible use cases, but probably won't spend time on it unless it's requested (open a discussion fi you want this feature).

Currently has support for reading squashfs files and extracting files and folders.

Special thanks to <https://dr-emann.github.io/squashfs/> for some VERY important information in an easy to understand format.
Thanks also to [distri's squashfs library](https://github.com/distr1/distri/tree/master/internal/squashfs) as I referenced it to figure some things out (and double check others).

## [TODO](https://github.com/CalebQ42/squashfs/projects/1?fullscreen=true)

## Limitations

* No Xattr parsing. This is simply because I haven't done any research on it and how to apply these in a pure go way.
* Socket files are not extracted.
  * From my research, it seems like a socket file would be useless if it could be created.
* Fifo files are ignored on `darwin`

## Issues

* Larger, more complex archives have significant issues when doing a full extraction.
  * Seems to be mostly a problem with archives with many deep file trees.
  * It seems to not only take exponentially longer per nested folder, but also will eat all your system's memory as it does so.
  * Observed when tested Arch Linux's live iso's airootfs.sfs.
  * Accessing files / folders without extracting is NOT be effected.
* Significantly slower then `unsquashfs` (about 5 ~ 7 times slower on a ~100MB archive using zstd compression)
  * This seems to be related to above along with the general optimization of `unsquashfs` and it's compression libraries.

## Recommendations on Usage

Due to the above issue and performance consideration, this library should only be used to access files within the archive without extraction, or to mount it via Fuse.

* Neither of these use cases are largely effected by the issues above.
