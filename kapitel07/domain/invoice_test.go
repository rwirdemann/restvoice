package domain

import (
	"testing"

	"github.com/rwirdemann/restvoice/kapitel05/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddPosition(t *testing.T) {
	// Setup
	i := domain.Invoice{}

	// Run
	i.AddPosition(1, "Programmierung", 20, 60)
	i.AddPosition(1, "Programmierung", 12, 60)
	i.AddPosition(1, "Qualitätssicherung", 3, 55)
	i.AddPosition(2, "Projektmanagement", 24, 50)
	i.AddPosition(2, "Qualitätssicherung", 8, 55)

	// Assert
	expected := domain.Position{Hours: 32, Price: 1920}
	assert.Equal(t, expected, i.Positions[1]["Programmierung"])
	expected = domain.Position{Hours: 3, Price: 165}
	assert.Equal(t, expected, i.Positions[1]["Qualitätssicherung"])
	expected = domain.Position{Hours: 24, Price: 1200}
	assert.Equal(t, expected, i.Positions[2]["Projektmanagement"])
	expected = domain.Position{Hours: 8, Price: 440}
	assert.Equal(t, expected, i.Positions[2]["Qualitätssicherung"])
}
