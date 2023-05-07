package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"url_shortener/internal/model"
	"url_shortener/internal/url"
)

type urlHandler struct {
	urlUseCase url.ShortExpander
}

func NewURLHandler(uc url.ShortExpander) url.Handler {
	return &urlHandler{urlUseCase: uc}
}

func (uh *urlHandler) Handle(w http.ResponseWriter, r *http.Request)  {
	var urlRequest model.URL

	shortenedURL := r.URL.Path[1:]
	if shortenedURL != "" {
		urlRequest.Short = strings.TrimSpace(shortenedURL)
		//redirect
			long, err := uh.urlUseCase.Expand(urlRequest.Short)
			if err != nil {
				err = fmt.Errorf("expanding url:%w", err)
				log.Printf("%s", err)
				http.Error(w, "Unable to process request", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, long, http.StatusMovedPermanently)
			return
	}

	long, err := uh.getLongUrlFromBody(r)
	if err != nil {
		log.Printf("getting long url from body: %s",err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	short := uh.urlUseCase.Shorten(long)

	_, err = io.WriteString(w, "http://localhost:1234/"+short)
	if err != nil {
		err = fmt.Errorf("shortening url:%w", err)
		log.Printf("%s", err)
		http.Error(w, "Unable to process request", http.StatusInternalServerError)
	}
	return

}

func (uh *urlHandler) getLongUrlFromBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", fmt.Errorf("empty request body")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("reading body:%w", err)
	}

	l := strings.TrimSpace(string(body))
	return l, nil
}