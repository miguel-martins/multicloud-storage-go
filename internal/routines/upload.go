package routines

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/fentec-project/gofe/abe"
)

const (
	accountName    = "your_account_name"
	accountKey     = "your_account_key"
	containerName  = "your_container_name"
	blobServiceURL = "https://your_account_name.blob.core.windows.net/"
)

// uploadChunks uploads encrypted chunks to Azure Blob Storage.
func UploadChunks(encryptedChunkChannel <-chan EncryptedChunk, wg *sync.WaitGroup) {
	defer wg.Done()
	for encryptedChunk := range encryptedChunkChannel {
		log.Printf("Uploading encrypted chunk to Azure Blob Storage %v", len(encryptedChunk.Data))
		_, err := decrypt(encryptedChunk.Data)
		if err != nil {
			log.Print("Not able to decrypt")
		} else {
			log.Print("Decrypted Chunk")
		}
	}
}

func decrypt(data []byte) (string, error) {
	deserialized, err := deserializeCipher(data)
	if err != nil {
		return "", err
	}

	decrypted, err := decryptCipher(deserialized)
	if err != nil {
		return "", err
	}

	return decrypted, nil

}

func decryptCipher(ciphertext *abe.FAMECipher) (string, error) {
	attribKeys, _ := fame.GenerateAttribKeys([]string{"attrib1"}, sk)
	decryptedCipherText, err := fame.Decrypt(ciphertext, attribKeys, pk)
	if err != nil {
		return "", err
	}
	return decryptedCipherText, nil
}

func deserializeCipher(serializedCipher []byte) (*abe.FAMECipher, error) {
	// Create a new FAMECipher object
	cipher := &abe.FAMECipher{}
	// Deserialize the JSON data into the cipher object
	err := json.Unmarshal(serializedCipher, cipher)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}

func uploadToAzureBlobStorage(data []byte) error {
	return nil
	// credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	// if err != nil {
	// 	return err
	// }
	// p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	// u, _ := url.Parse(fmt.Sprintf("%s%s/%s", blobServiceURL, containerName, "encrypted_chunk"))
	// blockBlobURL := azblob.NewBlockBlobURL(*u, p)

	// _, err = azblob.UploadStreamToBlockBlob(ctx, bytes.NewReader(data), blockBlobURL,
	// 	azblob.UploadStreamToBlockBlobOptions{
	// 		BufferSize: 4 * 1024 * 1024,
	// 		MaxBuffers: 3,
	// 	})
	// if err != nil {
	// 	return err
	// }
	// return nil
}
