package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/miguel-martins/multicloud-storage-go/internal/routines"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	filename := header.Filename
	//username := r.Context().Value("username").(string)

	// TODO
	// Register file with filename, filesize, owner and uploadDate of now().
	// Create initial attributes for file and send to encrypting routine

	// Create another encryption channel, whose job is to update encryption keys etc for deduplicated chunks

	chunkChannel := make(chan routines.Chunk)
	encryptedChunkChannel := make(chan routines.EncryptedChunk)

	var wg sync.WaitGroup
	wg.Add(3)

	// Start chunking stage
	go routines.ChunkFile(file, filename, chunkChannel, &wg)

	// Start encryption stage
	go routines.EncryptChunks(chunkChannel, encryptedChunkChannel, &wg)

	// Start upload stage
	go routines.UploadChunks(encryptedChunkChannel, &wg)

	wg.Wait()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&UploadResponse{Success: true, Filename: filename})

}

type UploadResponse struct {
	Success  bool
	Filename string
}
