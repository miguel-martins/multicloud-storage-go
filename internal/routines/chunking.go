package routines

import (
	"io"
	"log"
	"mime/multipart"
	"sync"

	"github.com/jotfs/fastcdc-go"
)

// Chunk represents a file chunk.
type Chunk struct {
	Data []byte
}

// chunkFile reads and chunks the specified file.
func ChunkFile(file multipart.File, filename string, chunkChannel chan<- Chunk, wg *sync.WaitGroup) {
	defer close(chunkChannel)
	defer wg.Done()

	log.Printf("Chunking file %v\n", filename)

	// Chunk the file using fastcdc-go
	opts := fastcdc.Options{
		MinSize:     256 * 1024,
		AverageSize: 1 * 512 * 1024,
		MaxSize:     4 * 1024 * 1024,
	}

	chunker, _ := fastcdc.NewChunker(file, opts)
	for {
		chunk, err := chunker.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// Register New Chunk with FileId, chunk index, fingerprint

		// Check if chunk.Fingerprint exists in Chunks table
		// If fingerprint doesnt exist Send chunk to encrypting channel
		// If fingerprint exists, send to different encrypting channel, since the decryption key needs to have more attributes

		chunkChannel <- Chunk{Data: chunk.Data}
	}
}
