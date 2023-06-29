package helper

import (
	"fmt"
	"regexp"
)

func IsValidURL(input string) bool {
	regex := `^(?:(?:(?:https?|ftp):)?\/\/)(?:\S+(?::\S*)?@)?(?:(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-zA-Z0-9\p{L}\p{N}][a-zA-Z0-9\p{L}\p{N}_-]{0,62})?[a-zA-Z0-9\p{L}\p{N}]\.)+(?:[a-zA-Z\p{L}]{2,}\.?))(?::\d{2,5})?(?:[/?#]\S*)?$`
	match, err := regexp.MatchString(regex, input)
	if err != nil {
		fmt.Println("Error parsing URL with regex:", err)
		return false
	}

	return match

}
func hasNonASCIIDomain(domain string) bool {
	for _, c := range domain {
		if c > 127 {
			return true
		}
	}
	return false
}

//var validURLs = []string{
//	"http://foo.com/blah_blah",
//	"http://foo.com/blah_blah/",
//	"http://foo.com/blah_blah_(wikipedia)",
//	"http://foo.com/blah_blah_(wikipedia)_(again)",
//	"http://www.example.com/wpstyle/?p=364",
//	"https://www.example.com/foo/?bar=baz&inga=42&quux",
//	"https://www.example.com",
//	"http://✪df.ws/123",
//	"http://userid:password@example.com:8080",
//	"http://userid:password@example.com:8080/",
//	"http://userid@example.com",
//	"http://userid@example.com/",
//	"http://userid@example.com:8080",
//	"http://userid@example.com:8080/",
//	"http://userid:password@example.com",
//	"http://userid:password@example.com/",
//	"http://➡.ws/䨹",
//	"http://⌘.ws",
//	"http://⌘.ws/",
//	"http://foo.com/blah_(wikipedia)#cite-1",
//	"http://foo.com/blah_(wikipedia)_blah#cite-1",
//	"http://foo.com/unicode_(✪)_in_parens",
//	"http://foo.com/(something)?after=parens",
//	"http://☺.damowmow.com/",
//	"http://code.google.com/events/#&product=browser",
//	"http://j.mp",
//	"ftp://foo.bar/baz",
//	"http://foo.bar/?q=Test%20URL-encoded%20stuff",
//	"http://مثال.إختبار",
//	"http://例子.测试",
//	"http://-.~_!$&'()*+,;=:%40:80%2f::::::@example.com",
//	"http://1337.net",
//	"http://a.b-c.de",
//	"http://223.255.255.254",
//}
//
//var invalidURLs = []string{
//	"http://",
//	"http://.",
//	"http://..",
//	"http://../",
//	"http://?",
//	"http://??",
//	"http://??/",
//	"http://#",
//	"http://##",
//	"http://##/",
//	"http://foo.bar?q=Spaces should be encoded",
//	"//",
//	"//a",
//	"///a",
//	"///",
//	"http:///a",
//	"foo.com",
//	"rdar://1234",
//	"h://test",
//	"http:// shouldfail.com",
//	":// should fail",
//	"http://foo.bar/foo(bar)baz quux",
//	"ftps://foo.bar/",
//	"http://-error-.invalid/",
//	"http://-a.b.co",
//	"http://a.b-.co",
//	"http://0.0.0.0",
//	"http://3628126748",
//	"http://.www.foo.bar/",
//	"http://www.foo.bar./",
//	"http://.www.foo.bar./",
//}
