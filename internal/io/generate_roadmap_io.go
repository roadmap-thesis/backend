package io

import "github.com/roadmap-thesis/backend/internal/domain/object"

type GenerateRoadmapInput struct {
	Topic                  string                        `json:"topic" validate:"required"`
	PersonalizationOptions object.PersonalizationOptions `json:"personalization_options" validate:"required"`
}

type GenerateRoadmapOutput struct {
	Slug string `json:"slug"`
}
