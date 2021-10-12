package main

import (
	"pack.ag/tftp"
	"log"
	"os"
	"io"
	"io/ioutil"
)

func main() {
	s, err := tftp.NewServer(":69", tftp.ServerSinglePort(true))
	if err != nil {
		panic(err)
	}
	readHandler := tftp.ReadHandlerFunc(getFile)
	s.ReadHandler(readHandler)

	writeHandler := tftp.WriteHandlerFunc(putFile)
	s.WriteHandler(writeHandler)
	s. ListenAndServe()
	select{}

}

func getFile(w tftp.ReadRequest) {
	log.Printf("[%s] GET %s\n", w.Addr().IP.String(), w.Name() )
	file, err := os.Open("/tftpboot/" + w.Name()) // For read access.
	if err != nil {
		log.Println(err)
		w.WriteError(tftp.ErrCodeFileNotFound, err.Error())
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		log.Println(err)
	}
}

func putFile(w tftp.WriteRequest) {
	log.Printf("write request for %s", w.Name())

	// Get the file size
	size, err := w.Size()
	log.Printf("file size from client %d", size)

	// check size of uploaded file is smaller than we're willing to store (32MB)
	// An error indicates no size was received.
	//if err != nil || size > 32*1024*1024 {
		// Send a "disk full" error.
		//w.WriteError(tftp.ErrCodeDiskFull, "File too large or no size sent")
		//return
	//}

	// Note: The size value is sent by the client, the client could send more data than
	// it indicated in the size option. To be safe we'd want to allocate a buffer
	// with the size we're expecting and use w.Read(buf) rather than ioutil.ReadAll.

	// Read the data from the client into memory
	data, err := ioutil.ReadAll(w)
	if err != nil {
		log.Println(err)
		return
	}

	err = os.WriteFile("/tftpboot/" + w.Name(), data, 0444)
	if err != nil {
		log.Println(err)
		w.WriteError(tftp.ErrCodeDiskFull, err.Error())
		return
	}

	// Log a message with the details
	log.Printf("[%s] PUT %s %d bytes\n", w.Addr().IP.String(), w.Name(), len(data) )
}
