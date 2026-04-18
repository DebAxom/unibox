package utils

import (
	"strings"
	"unibox/models"
)

var deptKeywords = map[string][]string{
	"hostel": {
		"fan", "light", "electricity", "water", "bathroom", "room",
		"lan", "wifi", "mess", "food", "clean", "hostel", "refund", "vending",
		"commode", "jetspray", "hanger", "block",
	},
	"academic": {
		"class", "lecture", "exam", "attendance", "teacher",
		"faculty", "assignment", "marks", "leave", "admission",
		"jossa", "extra", "cgpa", "sgpa",
	},
	"accounts": {
		"tecnoesis", "incandescence", "incand", "tecno", "prize",
	},
	"sw": {
		"scholarship", "harassment", "complaint", "mental",
		"counselling", "ragging", "refund", "admission", "suicide",
		"classroom", "infrastructure",
	},
}

func classifyDept(title, description string) string {
	text := strings.ToLower(title + " " + description)

	scores := make(map[string]int)

	for dept, keywords := range deptKeywords {
		for _, word := range keywords {
			if strings.Contains(text, word) {
				scores[dept]++
			}
		}
	}

	// Find max score
	bestDept := "none"
	maxScore := 0

	for dept, score := range scores {
		if score > maxScore {
			maxScore = score
			bestDept = dept
		}
	}

	return bestDept
}

var hostels = map[string]string{
	"abh": "hostel-abh",
	"bh8": "hostel-bh8",
	"gh1": "hostel-gh1",
	"gh3": "hostel-gh3",
}

func RouteComplain(user models.User, title string, desc string) string {
	dept := classifyDept(title, desc)

	if dept == "none" {
		dept = "hostel"
	}

	if dept == "hostel" {
		dept = hostels[strings.ToLower(user.Hostel)]
	}

	return dept
}
