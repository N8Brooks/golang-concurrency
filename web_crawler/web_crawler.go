//go:build challenge

package webcrawler

type Challenge struct{}

func NewChallenge() *Challenge {
	return &Challenge{}
}

func (c *Challenge) Crawl(startURL string, htmlParser HTMLParser) []string {
	panic("unimplemented")
}
