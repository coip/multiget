# multiget - chunked downloads

Downloads specified URL resource in 4 chunks, and optionally (with -whole) the remainder afterwards.

build with `go build`

try `./multiget -h` for help

usage:

*./multiget [flags] _url_*

ex ` ./multiget -chunks=58 -chunksize=8192 -filename="test.jpg" https://i.imgur.com/VG2UvcY.jpg`

further exploration:

  - *there is likely better means of handling the http.Client & \*os.File regarding scope*
  
  - *investigate WriteAt vs seek&Copy*
