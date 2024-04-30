package organizer

var EventCodeToFullMap = map[string]string{
	"333":    "3x3x3",
	"222":    "2x2x2",
	"444":    "4x4x4",
	"555":    "5x5x5",
	"666":    "6x6x6",
	"777":    "7x7x7",
	"333bf":  "3x3x3 Blindfolded",
	"333fm":  "3x3x3 Fewest Moves",
	"333oh":  "3x3x3 One-Handed",
	"clock":  "Clock",
	"minx":   "Megaminx",
	"pyram":  "Pyraminx",
	"skewb":  "Skewb",
	"sq1":    "Square-1",
	"444bf":  "4x4x4 Blindfolded",
	"555bf":  "5x5x5 Blindfolded",
	"333mbf": "3x3x3 Multiple Blindfolded",
}

var EventFullToCodeMap = map[string]string{
	"3x3x3":                      "333",
	"2x2x2":                      "222",
	"4x4x4":                      "444",
	"5x5x5":                      "555",
	"6x6x6":                      "666",
	"7x7x7":                      "777",
	"3x3x3 Blindfolded":          "333bf",
	"3x3x3 Fewest Moves":         "333fm",
	"3x3x3 One-Handed":           "333oh",
	"Clock":                      "clock",
	"Megaminx":                   "minx",
	"Pyraminx":                   "pyram",
	"Skewb":                      "skewb",
	"Square-1":                   "sq1",
	"4x4x4 Blindfolded":          "444bf",
	"5x5x5 Blindfolded":          "555bf",
	"3x3x3 Multiple Blindfolded": "333mbf",
}
