package routines

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/fentec-project/gofe/abe"
)

// EncryptedChunk represents an encrypted file chunk.
type EncryptedChunk struct {
	Data []byte
}

var (
	fame        *abe.FAME = abe.NewFAME()
	pk, sk, err           = fame.GenerateMasterKeys()
)

// encryptChunks encrypts the chunks received from the chunking stage.
func EncryptChunks(chunkChannel <-chan Chunk, encryptedChunkChannel chan<- EncryptedChunk, wg *sync.WaitGroup) {
	defer close(encryptedChunkChannel)
	defer wg.Done()

	for chunk := range chunkChannel {
		log.Printf("Received chunk: %v", len(chunk.Data))

		// TODO encrypt chunk with a random generated key from TA or Convergent Encryption
		// Encrypt encryption key by sending access attributes and the to-be-encrypted-content to TA
		encryptedData := encrypt(chunk.Data)

		encryptedChunkChannel <- EncryptedChunk{Data: encryptedData}
	}
}

func encrypt(data []byte) []byte {
	msp, _ := abe.BooleanToMSP("attrib1 AND (attrib2 OR attrib3)", true)
	cipher, _ := fame.Encrypt(string(data), msp, pk)
	serializedCipher, _ := serializeCipher(cipher)
	return serializedCipher
}

func serializeCipher(cipher *abe.FAMECipher) ([]byte, error) {
	// Serialize the cipher object to JSON
	serializedCipher, err := json.Marshal(cipher)
	if err != nil {
		return nil, err
	}
	return serializedCipher, nil
}
