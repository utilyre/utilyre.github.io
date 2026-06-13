package ingestion

import (
	"io"
	"website/internal/domain"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type contentDTO struct {
	Journal []journalLogDTO `yaml:"journal" validate:"min=1"`
}

type journalLogDTO struct {
	Title   string            `yaml:"title" validate:"required"`
	Year    uint              `yaml:"year"`
	Authors []string          `yaml:"authors" validate:"min=1,dive,required"`
	Links   map[string]string `yaml:"links" validate:"dive,required"`
}

func IngestYAML(r io.Reader) (domain.Content, error) {
	validate := validator.New()

	var dto contentDTO
	err := yaml.NewDecoder(r, yaml.Strict(), yaml.Validator(validate)).Decode(&dto)
	if err != nil {
		return domain.Content{}, err
	}

	var content domain.Content

	for _, logDTO := range dto.Journal {
		log := domain.JournalLog{
			Title:   logDTO.Title,
			Year:    logDTO.Year,
			Authors: append([]string(nil), logDTO.Authors...),
		}

		for label, url := range logDTO.Links {
			log.Links = append(log.Links, domain.LabeledLink{
				Label: label,
				URL:   url,
			})
		}

		content.Journal = append(content.Journal, log)
	}

	return content, nil
}
