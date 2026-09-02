package main

import (
	"html"
	"testing"
)

func TestParseVacanciesFromSearchHTML(t *testing.T) {
	body := []byte(`
<a data-qa="serp-item__title" href="https://hh.ru/vacancy/123?query=analyst"><span>Системный аналитик</span></a>
<a data-qa="vacancy-serp__vacancy_response" href="/applicant/vacancy_response?vacancyId=123&amp;employerId=7">Откликнуться</a>
<a data-qa="serp-item__title" href="https://hh.ru/vacancy/456"><span>Бизнес &amp; аналитик</span></a>
<a data-qa="vacancy-serp__vacancy_response" href="/applicant/vacancy_response?vacancyId=456">Откликнуться</a>
`)

	vacancies, err := parseVacanciesFromSearchHTML(body)
	if err != nil {
		t.Fatalf("parse vacancies: %v", err)
	}
	if len(vacancies) != 2 {
		t.Fatalf("vacancy count = %d, want 2", len(vacancies))
	}
	if vacancies[0].ID != 123 || vacancies[0].Name != "Системный аналитик" {
		t.Fatalf("first vacancy = %#v", vacancies[0])
	}
	if vacancies[1].Links["desktop"] != "https://hh.ru/vacancy/456" || vacancies[1].Name != "Бизнес & аналитик" {
		t.Fatalf("second vacancy = %#v", vacancies[1])
	}
}

func TestFetchVacanciesDecodesEscapedEmbeddedJSON(t *testing.T) {
	body := []byte(`prefix,&#34;total&#34;:1,&#34;vacancies&#34;:[{&#34;vacancyId&#34;:123,&#34;name&#34;:&#34;Системный аналитик&#34;,&#34;@responseLetterRequired&#34;:true,&#34;userTestPresent&#34;:true,&#34;links&#34;:{&#34;desktop&#34;:&#34;https://hh.ru/vacancy/123&#34;}}]}suffix`)
	decoded := []byte(html.UnescapeString(string(body)))

	var vacancies []Vacancy
	if err := decodeEmbeddedJSON(decoded, `,"vacancies":`, &vacancies); err != nil {
		t.Fatalf("decode escaped embedded JSON: %v", err)
	}
	if len(vacancies) != 1 || !vacancies[0].ResponseLetterRequired || !vacancies[0].UserTestPresent {
		t.Fatalf("vacancies = %#v", vacancies)
	}
}
