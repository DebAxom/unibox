package utils

import (
	"strings"
	"unibox/models"
)

var deptKeywords = map[string]map[string]int{
	"hostel": {
		"wifi": 3, "lan": 3, "internet": 3,
		"fan": 2, "light": 2, "electricity": 2,
		"water": 2, "bathroom": 2, "room": 2,
		"mess": 3, "food": 2, "clean": 2,
		"commode": 2, "jetspray": 2, "hanger": 1,
		"block": 1, "hostel": 2, "aryabhatta": 2,
		"gh3": 2, "bh8": 2, "abh": 2,
	},
	"academic": {
		"exam": 3, "cgpa": 3, "sgpa": 3,
		"class": 2, "lecture": 2, "attendance": 2,
		"faculty": 2, "teacher": 2,
		"assignment": 2, "marks": 2,
		"admission": 2, "jossa": 2, "csab": 2,
	},
	"accounts": {
		"fee": 3, "payment": 3, "refund": 3,
		"tecnoesis": 2, "incandescence": 2,
		"prize": 2, "contset": 1, "competition": 1,
	},
	"sw": {
		"scholarship": 3,
		"harassment":  3, "ragging": 3,
		"mental": 2, "counselling": 2,
		"suicide": 3, "murder": 3,
		"sports": 1, "health": 1,
	},
}

// Hostel Mapping
var hostels = map[string]string{
	"abh": "hostel-abh",
	"bh8": "hostel-bh8",
	"gh1": "hostel-gh1",
	"gh3": "hostel-gh3",
}

// Normalize text
func normalize(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "-", " ")
	text = strings.ReplaceAll(text, "_", " ")
	text = strings.ReplaceAll(text, ".", " ")
	return text
}

// Basic typo tolerance (cheap)
func normalizeWord(w string) string {
	switch w {
	case "wi-fi", "wifi", "wfi":
		return "wifi"
	case "ln":
		return "lan"
	}
	return w
}

// Core Classifier
func classifyDept(title, description string) (string, int) {

	// Boost title importance
	text := normalize(title + " " + title + " " + description)

	words := strings.Fields(text)

	scores := make(map[string]int)

	for _, word := range words {
		word = normalizeWord(word)

		for dept, keywords := range deptKeywords {
			if weight, ok := keywords[word]; ok {
				scores[dept] += weight
			}
		}
	}

	bestDept := "none"
	maxScore := 0

	for dept, score := range scores {
		if score > maxScore {
			maxScore = score
			bestDept = dept
		}
	}

	return bestDept, maxScore
}

// Final Routing Function
func RouteComplain(user models.User, title, desc string) string {

	dept, score := classifyDept(title, desc)

	// Confidence threshold
	if score < 2 {
		return "other"
	}

	// Hostel expansion
	if dept == "hostel" {
		dept = hostels[strings.ToLower(user.Hostel)]
	}

	return dept
}
