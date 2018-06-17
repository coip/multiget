# multiget - chunked downloads

Downloads specified URL resource in 4 chunks, and the remainder afterwards.
try `./multigetclient -h` for help
usage:
*./multigetclient [flags] _url_*
ex
./multigetclient -chunks=6 -chunksize=1024 -filename="test.jpg" -whole https://i.imgur.com/VG2UvcY.jpg 
