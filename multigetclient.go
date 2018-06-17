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
)

func check(e error) {
   if e != nil { panic(e) }
}

func downloadChunk(resourceUrl string, chunk int, chunksize int, file *os.File) {
   fmt.Println("----------------------------\n--\nrequesting chunk ", chunk, " on ", resourceUrl, "\n-------------------------")

      //construct header for range, specify (unit, star, end)
      headerval := fmt.Sprintf("%s=%d-%d", "bytes", chunk*chunksize, (chunk+1)*chunksize-1)
   fmt.Println("[Range]:[", headerval, "]")

      //construct request
      client := &http.Client{}
      req, _ := http.NewRequest("GET", resourceUrl, nil)
      req.Header.Set("Range", headerval)

      res, _ := client.Do(req)

      file.Seek(int64(chunk*chunksize), 0)
      bytesWritten, _ := io.Copy(file, res.Body)

   fmt.Println("bytes written: ", bytesWritten)
      defer res.Body.Close()

   fmt.Println("response status: ", res.Status);
//      check(err);
}

func main() {

//resource url to download
   var resourceUrl string = os.Args[1]
   
//filename
   var filename string = os.Args[2]

//size of chunk, mebibyte. sticking with a KB for testing w images
   chunksize := 1024//*1024
   
//chunks to download
   chunks := 4

//rudimentary check for url presence. might want to validate w ~
   if len(resourceUrl) == 0 {
      fmt.Println("specify resource url for download\nie - ./multigetclient url");
      return;
   }

   out, _ := os.Create(filename)
   defer out.Close()

//get the chunks
   for chunk := 0; chunk < chunks; chunk++ {
      downloadChunk(resourceUrl, chunk, chunksize, out)
      //go downloadChunk(resourceUrl, chunk, chunksize, out)
   }
   fmt.Println("downloaded ", chunks, " chunks, each ", chunksize, " bytes. finishing...");


   headerval := fmt.Sprintf("%s=%d-", "bytes", (chunks)*chunksize)
   fmt.Println("[Range]:[", headerval, "]")

      //construct request
      client := &http.Client{}
      req, _ := http.NewRequest("GET", resourceUrl, nil)
      req.Header.Set("Range", headerval)

      res, _ := client.Do(req)

      out.Seek(0, 2) 
      //out.Seek(int64(chunks*chunksize),0)
      bytesWritten, _ := io.Copy(out, res.Body)

   fmt.Println("bytes written: ", bytesWritten)
      defer res.Body.Close()

   fmt.Println("response status: ", res.Status);

}
