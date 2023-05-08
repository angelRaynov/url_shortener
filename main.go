package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"url_shortener/internal/repository"
	"url_shortener/internal/url/delivery/http"
	"url_shortener/internal/url/usecase"
)


func main() {
	repo := repository.NewURLRepository()
	useCase := usecase.NewURLUseCase(repo)
	handler := http.NewURLHandler(useCase)

	router := gin.Default()
	router.POST("/", handler.ShortenURL)
	router.GET("/:shortened", handler.ExpandURL)

	log.Println("listening on port :1234")

	log.Fatal(router.Run("localhost:1234"))

}


//// generateUniqueURLs sends a unique url to a shared channel
//// it will create over a billion unique combinations
//// without the need to check if value exists
//func generateUniqueURLs(c chan string) {
//	// Calculate the number of combinations per chunk
//	// Each chunk starts with a different character from the chars set
//	charsetLen := len(charset)
//	combinationsPerChunk := int(math.Pow(float64(charsetLen), 4))
//	// Generate short URLs for each chunk
//	// Loop through all characters in the chars set to start a new chunk
//	for i := 0; i < charsetLen; i++ {
//		// Calculate the start and end index for this chunk
//		// startIndex is the index where this chunk starts in the combinations list
//		startIndex := i * combinationsPerChunk
//		// endIndex is the index where this chunk ends in the combinations list
//		endIndex := (i + 1) * combinationsPerChunk
//		// Loop through all the combinations in this chunk
//		for j := startIndex; j < endIndex; j++ {
//			// url is represented as a string of characters from the chars set
//			url := ""
//			nextCharIndex := j
//			for k := 0; k < 5; k++ {
//				// The next character is determined by the remainder of nextCharIndex divided by charsetLen
//				// Avoid += as it will append next char in the front
//				url = string(charset[nextCharIndex%charsetLen]) + url
//				// Update nextCharIndex to the next digit of the charset
//				nextCharIndex /= charsetLen
//			}
//			//Send the unique url and block until read
//			c <- fmt.Sprintf("%s", url)
//		}
//	}
//
//}
