package helper

import (
	"fmt"
	"golang.org/x/net/idna"
	"net/url"
	"strings"
)

func replaceUnicodeSymbolsWithASCII(input string) string {
	//check if the url contains non ascii characters and try to convert the long url to ascii
	parsedURL, err := url.Parse(input)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return input
	}
	if parsedURL.Host == "" {
		fmt.Println("unable to parse host")
		return input
	}
	// Check if the host contains non-ASCII characters
	isNonASCIIDomain := hasNonASCIIDomain(parsedURL.Host)

	if isNonASCIIDomain {
		// Encode the domain using Punycode
		encodedHost, err := idna.ToASCII(parsedURL.Host)
		if err != nil {
			fmt.Println("Error encoding Punycode:", err)
			return input
		}

		// Rebuild the URL with the encoded domain
		parsedURL.Host = encodedHost
	}

	// Get the converted URL string
	fmt.Println("converted", parsedURL.String())
	return parsedURL.String()

}

func prependProtocolScheme(input string) string {
	if !strings.HasPrefix(input, "https://") && !strings.HasPrefix(input, "http://") {
		return "http://" + input
	}
	return input
}

func LinkPreconditioning(input string) string {
	input = prependProtocolScheme(input)
	input = replaceUnicodeSymbolsWithASCII(input)
	return input
}
