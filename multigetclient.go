package main

//v2: seemingly works for 4 chunks of MiB on test jar.
//    gutenberg apparently wasnt supporting Range (saw 200, not 206)
//    testing seems positive w/ https://i.imgur.com/VG2UvcY.jpg


//v1: downloads file into download, using Range Header.
//    seems to work for sub-mebibyte files, ie https://www.gutenberg.org/files/57328/57328-0.txt
//    & stat reports  of a 20MB payload


import (
   "fmt"
   "os"
   "io"
   "net/http"
   "flag"
   "sync"
)

var wg sync.WaitGroup
var mu sync.Mutex
var file *os.File



//simple error wrapper from gobyexample
func check(e error) {
   if e != nil { panic(e) }
}

func downloadChunk(resourceUrl string, chunk int, chunksize int) {
      defer wg.Done();
      fmt.Println("--\nrequesting chunk ", chunk, " on ", resourceUrl, "\n--")

      //construct header for range, specify (unit, start, end)
      headerval := fmt.Sprintf("%s=%d-%d", "bytes", chunk*chunksize, (chunk+1)*chunksize-1)
      fmt.Println("[Range]:[", headerval, "]")

      //construct request
      client := &http.Client{}
      req, err := http.NewRequest("GET", resourceUrl, nil)
      check(err)
      req.Header.Set("Range", headerval)

      //perform request
      res, err := client.Do(req)
      check(err)

      mu.Lock()
      //seek file pointer to correct location (prep for parallel)
      file.Seek(int64(chunk*chunksize), 0)
      //io.Copy(file, res.Body)

      bytesWritten, err := io.Copy(file, res.Body)
      mu.Unlock()
      fmt.Println("bytes written: ", bytesWritten)
      check(err)
      defer res.Body.Close()
      fmt.Println("response status: ", res.Status);
      //check(err);
}

func main() {


   //establish defaults for cli args
   filenamePtr    := flag.String("filename", "download", "filename to save download to")
   chunksizePtr   := flag.Int("chunksize", 1024*1024, "chunk size in bytes")
   chunksPtr      := flag.Int("chunks", 4, "how many chunks")
   wholePtr       := flag.Bool("whole", false, "true ? download all : just the chunks specified")

   //parse flags, more importantly: populate flag.Args for url
   flag.Parse()

   //rudimentary check for url presence. might want to validate w ~
   if len(flag.Args()) == 0 {
      fmt.Println("specify resource url for download\nie - ./multigetclient url");
      return;
   }

   //resource url to download, could support list simply enough
   var resourceUrl string = flag.Args()[0]

   //create local file to save downloaded contents into
   file, _ = os.Create(*filenamePtr)
   defer file.Close()
   //check(err)

   //we will have "*chunksPtr" many go-routines to wait for
   wg.Add(*chunksPtr)

   //get the chunks
   for chunk := 0; chunk < *chunksPtr; chunk++ {
      //downloadChunk(resourceUrl, chunk, *chunksizePtr)
      go downloadChunk(resourceUrl, chunk, *chunksizePtr)
   }


   wg.Wait()
   fmt.Println("downloaded ", *chunksPtr, " chunks, each ", *chunksizePtr, " bytes.")

   //for testing/general usability ive added a final request to get remainder of data/payload
   if *wholePtr {
      fmt.Println("finishing...")

      headerval := fmt.Sprintf("%s=%d-", "bytes", (*chunksPtr)*(*chunksizePtr))
      fmt.Println("[Range]:[", headerval, "]")

      //construct request
      client := &http.Client{}
      req, _ := http.NewRequest("GET", resourceUrl, nil)
      req.Header.Set("Range", headerval)

      res, _ := client.Do(req)

      file.Seek(0, 2) 
      //file.Seek(int64(chunks*chunksize),0)
      bytesWritten, _ := io.Copy(file, res.Body)

      fmt.Println("bytes written: ", bytesWritten)
      defer res.Body.Close()
   
      fmt.Println("response status: ", res.Status);
   }
}
