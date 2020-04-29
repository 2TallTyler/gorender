package colour

import (
	"fmt"
	"strings"
	"testing"
)

const exampleJson = "{\"entries\": [[0,0,0],[255,255,255],[255,127,0]], \"ranges\": [{\"start\": 0, \"end\": 1, \"is_primary_company_colour\": true}]}"

func TestGetPaletteFromJson(t *testing.T) {
	expectedRanges := []PaletteRange{{0, 1, true, false}}
	expectedEntries := []PaletteEntry{{R: 0, G: 0, B: 0}, {R: 255, G: 255, B: 255}, {R: 255, G: 127, B: 0}}

	palette, err := GetPaletteFromJson(strings.NewReader(exampleJson))

	if err != nil {
		t.Errorf("encountered error: %v", err)
	}

	if len(palette.Ranges) != len(expectedRanges) {
		t.Fatalf("ranges length %d is too short, expected %d", len(palette.Ranges), len(expectedRanges))
	}

	if len(palette.Entries) != len(expectedEntries) {
		t.Fatalf("entries length %d is too short, expected %d", len(palette.Entries), len(expectedEntries))
	}

	for i, r := range expectedRanges {
		if palette.Ranges[i] != r {
			t.Errorf("palette range %d not loaded correctly: was %v, expected %v", i, palette.Ranges[i], r)
		}
	}

	for i, e := range expectedEntries {
		if palette.Entries[i].R != e.R || palette.Entries[i].G != e.G || palette.Entries[i].B != e.B {
			t.Errorf("palette entry %d not loaded correctly: was %v, expected %v", i, palette.Entries[i], e)
		}
	}
}

func TestPalette_GetRGB(t *testing.T) {
	palette, _ := GetPaletteFromJson(strings.NewReader(exampleJson))

	fmt.Printf("check!\n")
	palette.SetRanges([]PaletteRange{{Start: 2, End: 2, IsPrimaryCompanyColour: true}})

	expected := [][]uint32{{0, 0, 0}, {65535, 65535, 65535}, {38731, 38731, 38731}, {0, 0, 0}}

	for i, e := range expected {
		if r, g, b := palette.GetRGB(byte(i)); r != e[0] || g != e[1] || b != e[2] {
			t.Errorf("entry at %d not returned correctly: was [%d %d %d], expected %v", i, r, g, b, e)
		}
	}
}

func TestPalette_GetFromReader_DetectsDuplicateRanges(t *testing.T) {
	const json = "{\"entries\": [[0,0,0],[255,255,255],[255,127,0]], \"ranges\": [{\"start\": 0, \"end\": 1},{\"start\": 1, \"end\": 2}]}"
	_, err := GetPaletteFromJson(strings.NewReader(json))

	if err == nil || err.Error() != "range 1 overlaps colour 1" {
		t.Errorf("encountered unexpected error: %v", err)
	}
}