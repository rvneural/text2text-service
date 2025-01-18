package services

import (
	"strings"
	"sync"
)

type DigestText struct {
	Text string
	URL  string
}

func (s *Service) createDigestTexts(links []string) []DigestText {
	texts := make([]DigestText, len(links))
	wg := sync.WaitGroup{}
	for i, link := range links {
		wg.Add(1)
		go func(link string, index int) {
			defer wg.Done()
			var text string
			var err error
			if strings.Contains(link, "realnoevremya.ru") {
				text, err = s.rvParser.ParseRV(link)
				if text != "" {
					s_link := strings.ReplaceAll(link, "https://realnoevremya.ru", "")
					s_link = strings.ReplaceAll(s_link, "https://m.realnoevremya.ru", "")
					texts[index] = DigestText{Text: text, URL: s_link}
				}
			} else {
				text, err = s.anotherParser.Parse(link)
				if text != "" {
					texts[index] = DigestText{Text: text, URL: ""}
				}
			}
			if err != nil {
				s.logger.Error().Msg("Error while parsing url: " + err.Error())
				texts[index] = DigestText{Text: "НЕ УДАЛОСЬ ОБРАБОТАТЬ: " + link, URL: ""}
			}

		}(link, i)
	}
	wg.Wait()
	return texts
}

func (s *Service) findVerb(text string) (start, end int) {
	start = 0
	end = 0
	for i, s := range text {
		if s == '.' || s == '!' || s == '?' || s == ',' || s == ';' || s == ':' || s == '—' || s == '-' || s == '(' {
			if (i+1 < len(text) && text[i+1] == ' ') || (i+1 == len(text)) {
				end = i
				break
			}
		}
	}
	return start, end
}

func (s *Service) createDigestResult(texts []DigestText, model string) []string {
	const prompt = "{{ short }}"
	wg := sync.WaitGroup{}
	resultTextList := make([]string, len(texts))
	for i, text := range texts {
		wg.Add(1)
		go func(text DigestText, index int) {
			defer wg.Done()
			var resultText string
			var err error
			if !strings.HasPrefix(text.Text, "НЕ УДАЛОСЬ ОБРАБОТАТЬ: ") {
				resultText, err = s.ProcessText(model, prompt, text.Text, "0.1")
			} else {
				resultText = text.Text
			}

			if err != nil {
				s.logger.Error().Msg("Error while processing text: " + err.Error())
				return
			}

			if !strings.HasPrefix(resultText, "НЕ УДАЛОСЬ ОБРАБОТАТЬ: ") && text.URL != "" {
				url := "<a href=\"" + text.URL + "\" target=\"_blank\">"
				start, end := s.findVerb(resultText)
				resultText = resultText[:start] + url + resultText[start:end] + "</a>" + resultText[end:]
				resultText = strings.ReplaceAll(resultText, "\n", " ")
				resultText = strings.ReplaceAll(resultText, "  ", " ")
				resultText = strings.TrimSpace(resultText)
			}
			resultTextList[index] = resultText
		}(text, i)
	}
	wg.Wait()
	return resultTextList
}

func (s *Service) parseDigest(model, text string) string {

	links := make([]string, 0, 30)
	s.logger.Info().Msg("New request for DIGEST: " + text)
	// Split text by '\n' and ' ', and put it into links
	for _, line := range strings.Fields(text) {
		links = append(links, strings.Fields(line)...)
	}

	texts := s.createDigestTexts(links)

	wg := sync.WaitGroup{}
	var resultTexts = ""
	resultTextList := s.createDigestResult(texts, model)

	wg.Wait()
	for _, text := range resultTextList {
		resultTexts += "<p>" + text + "</p>\n\n"
	}
	resultTexts = strings.TrimSpace(resultTexts)
	return strings.TrimSpace(resultTexts)
}
