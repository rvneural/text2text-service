package rvparser

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type RVParser struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *RVParser {
	return &RVParser{
		logger: logger,
	}
}

func (p *RVParser) ParseRV(url string) (string, error) {
	return p.prepare(url)
}

func (p *RVParser) prepare(url string) (string, error) {
	var errMgs string = "НЕ УДАЛОСЬ ОБРАБОТАТЬ:" + url
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	p.logger.Debug().Msgf("status code: %d", resp.StatusCode)
	if resp.StatusCode != 200 {
		time.Sleep(time.Second)
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}
		if resp.StatusCode != 200 {
			return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
		}
	}
	if resp.Request.URL.Host != "realnoevremya.ru" && resp.Request.URL.Host != "m.realnoevremya.ru" {
		return "", fmt.Errorf("invalid URL")
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	str := string(bytes)
	i1 := strings.Index(str, "<article")
	i2 := strings.Index(str, "<div class=\"detailAuthors\"")
	if i2 == -1 && i1 != -1 {
		i2 = strings.Index(str, "<noindex>")
	}

	if i1 == -1 || i2 == -1 {
		if i1 == -1 {
			p.logger.Error().Msg("CAN'T FIND ARTICLE START for " + url)
		}
		if i2 == -1 {
			p.logger.Error().Msg("CAN'T FIND ARTICLE END for " + url)
		}

		return "", fmt.Errorf(errMgs)
	}

	str = str[i1:i2]
	str = strings.TrimSpace(strings.ReplaceAll(str, "<p><p>", "<p>"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "<br>", ""))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</p></br>", "</p>"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</p>", "</p>\n\n"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</h2>", "</h2>\n\n"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</h3>", "</h3>\n\n"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</ul>", "</ul>\n\n"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</table>", "</table>\n\n"))
	str = strings.TrimSpace(strings.ReplaceAll(str, "</figure>", "</figure>\n\n"))

	h1Idx1 := strings.Index(str, "<h1>")
	h1Idx2 := strings.Index(str, "</h1>")

	h1 := str[h1Idx1+4 : h1Idx2]

	textIdx1 := strings.Index(str, "<p")
	textIdx2 := strings.LastIndex(str, "</p>")

	text := p.clearHTML(str[textIdx1 : textIdx2+4])
	h1 = strings.TrimSpace(h1)
	h1 = "* " + h1 + " *"
	full := h1 + "\n\n" + text
	return full, nil
}

func (p *RVParser) clearHTMLTags(parts []string) (textParts []string) {

	textParts = make([]string, 0, len(parts))
	const Template = `<.*?>`
	r := regexp.MustCompile(Template)

	for _, part := range parts {
		if len(part) == 0 {
			continue
		} else if strings.HasPrefix(part, "<figure") {
			continue
		} else if strings.TrimSpace(part) == "</div>" {
			continue
		} else if strings.HasPrefix(part, "<div") {
			continue
		} else if strings.HasPrefix(part, "<table") {
			continue
		} else if strings.HasPrefix(part, "<p") {
			part = r.ReplaceAllString(part, "")
			part = strings.TrimSpace(part)
			if len(part) == 0 {
				continue
			}
			textParts = append(textParts, part)
		} else if strings.HasPrefix(part, "<ul") {
			part = strings.ReplaceAll(part, "<li>", "* ")
			part = strings.ReplaceAll(part, "</li>", "\n")
			part = r.ReplaceAllString(part, "")
			part = strings.TrimSpace(part)
			if len(part) == 0 {
				continue
			}
			textParts = append(textParts, part)
		} else {
			part = strings.ReplaceAll(part, "<h2>", "** ")
			part = strings.ReplaceAll(part, "</h2>", " **")
			part = strings.ReplaceAll(part, "<h3>", "*** ")
			part = strings.ReplaceAll(part, "</h3>", " ***")
			part = r.ReplaceAllString(part, "")
			part = strings.TrimSpace(part)
			if len(part) == 0 {
				continue
			}
			textParts = append(textParts, part)
		}
	}
	return textParts
}

func (p *RVParser) clearHTML(text string) string {
	parts := strings.Split(text, "\n")
	wg := sync.WaitGroup{}
	for i := 0; i < len(parts); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			parts[i] = strings.TrimSpace(parts[i])
		}(i)
	}
	wg.Wait()

	return strings.Join(p.clearHTMLTags(parts), "\n\n")
}
