package main

//v1: downloads file into download, using Range Header.
//    seems to work for sub-mebibyte files, ie https://www.gutenberg.org/files/57328/57328-0.txt
//    & stat reports  of a 20MB payload


import (
   "fmt"
   "os"
   "io/ioutil"
   "net/http"
)

func check(e error) {
   if e != nil { panic(e) }
}

func main() {

//resource url to download
   resourceUrl := os.Args[1]
   
//size of chunk, mebibyte. 
   chunkunit := 1024*1024
   
//chunks to download
   //chunks := 4

//rudimentary check for url presence. might want to validate w ~
   if len(resourceUrl) == 0 {
      fmt.Println("specify resource url for download\nie - ./multigetclient url");
      return;
   }

//construct header for range, specify (unit, star, end)
//**(-1 as our range later is 0-indexed)
   headerval := fmt.Sprintf("%s=%d-%d", "bytes", 0, chunkunit-1)

//contruct request
   client := &http.Client{}
   fmt.Println("requesting ", resourceUrl)
   req, _ := http.NewRequest("GET", resourceUrl, nil)
   req.Header.Set("Range", headerval)

   res, _ := client.Do(req)

//buffer download in memory for now, could be much more efficient
   download, _ := ioutil.ReadAll(res.Body)
   fmt.Println(res.Body);
   res.Body.Close()
   err := ioutil.WriteFile("download", download, 0644)
   check(err);

}
