package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseNACECodes(t *testing.T) {
	htmlContent := `<div class="nacelist">
		<ul class="level1">
			<li class="level1"><a href="/nace-code/01-agriculture-animal-husbandry-game-farming-and-related-service-activities.html" title="NACE 01 - Agriculture, animal husbandry, game farming and related service activities">01 - Agriculture, animal husbandry, game farming and related service activities</a>
				<ul class="level2">
					<li class="level2"><a href="/nace-code/011-cultivation-of-non-perennial-crops.html" title="NACE 011 - Cultivation of non-perennial crops">011 - Cultivation of non-perennial crops</a>
						<ul class="level3">
							<li class="level3"><a href="/nace-code/0111-cultivation-of-cereals-except-rice-leguminous-crops-and-oil-seeds.html" title="NACE 0111 - Cultivation of cereals (except rice), leguminous crops and oil seeds">0111 - Cultivation of cereals (except rice), leguminous crops and oil seeds</a></li>
						</ul>
					</li>
				</ul>
			</li>
		</ul>
	</div>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	naceCode := NACECode{}
	naceCode.Categories = parseCategories(doc, 1, "http://example.com")

	if len(naceCode.Categories) != 1 {
		t.Fatalf("Expected 1 main category, got %d", len(naceCode.Categories))
	}

	mainCategory := naceCode.Categories[0]
	if mainCategory.Code != "01" || mainCategory.Level != 1 {
		t.Errorf("Unexpected main category: %+v", mainCategory)
	}

	if len(mainCategory.SubCategories) != 1 {
		t.Fatalf("Expected 1 subcategory, got %d", len(mainCategory.SubCategories))
	}

	subCategory := mainCategory.SubCategories[0]
	if subCategory.Code != "011" || subCategory.Level != 2 {
		t.Errorf("Unexpected subcategory: %+v", subCategory)
	}

	if len(subCategory.SubCategories) != 1 {
		t.Fatalf("Expected 1 activity, got %d", len(subCategory.SubCategories))
	}

	activity := subCategory.SubCategories[0]
	if activity.Code != "0111" || activity.Level != 3 {
		t.Errorf("Unexpected activity: %+v", activity)
	}
}
