package googledrive

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// 	"golang.org/x/oauth2/jwt"
// )

// // ServiceAccount : Use Service account
// func ServiceAccount(credentialFile string) *http.Client {
// 	b, err := io.ReadFile(credentialFile)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var c = struct {
// 		Email      string `json:"client_email"`
// 		PrivateKey string `json:"private_key"`
// 	}{}
// 	json.Unmarshal(b, &c)
// 	config := &jwt.Config{
// 		Email:      c.Email,
// 		PrivateKey: []byte(c.PrivateKey),
// 		Scopes: []string{
// 			drive.DriveScope,
// 		},
// 		TokenURL: google.JWTTokenURL,
// 	}
// 	client := config.Client(oauth2.NoContext)
// 	return client
// }

// func UploadToDrive() {
// 	filename := "sample.txt"                    // Filename
// 	baseMimeType := "text/plain"                // MimeType
// 	client := ServiceAccount("credential.json") // Please set the json file of Service account.

// 	srv, err := drive.New(client)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fileInf, err := file.Stat()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer file.Close()
// 	f := &drive.File{Name: filename}
// 	res, err := srv.Files.
// 		Create(f).
// 		ResumableMedia(context.Background(), file, fileInf.Size(), baseMimeType).
// 		ProgressUpdater(func(now, size int64) { fmt.Printf("%d, %d\r", now, size) }).
// 		Do()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Printf("%s\n", res.Id)
// }
