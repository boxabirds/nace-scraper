package main

import (
	"strings"
	"testing"
)

func TestExtractRelevantText(t *testing.T) {
	htmlContent := `
	<div>
		<p>Some irrelevant text</p>
		<p><span class="szakagazat">Definition / Belongs to this division</span>Division 01 first of all
			distinguishes two basic activities: production of crop products and production of animal
			products. This section comprises the organic farming, the cultivation of genetically modified
			plants and the breeding of genetically modified animals. This section contains the outdoor and
			greenhouse crop production. The Mixed farming subsection (number 01.5) does not follow the
			general principles of the identification of the main activity. It accepts that many agricultural
			holdings have reasonably balanced crop and animal production, and that it would be arbitrary to
			classify them in one category or the other.<br><span class="szakagazatmagyar">also belongs to
				this sector</span> This group includes: <br> The services of the agriculture and game
			farming activities<br><span class="szakagazatnem">It does not belong to this
				sector</span>Agricultural activity excludes any subsequent processing of the agricultural
			products (classified under division 10 and 11: Manufacture of food and beverage products, and
			division 12: Manufacture of tobacco products), beyond that needed to prepare them for the
			primary markets. So that the preparation of the crops for the primary market belongs also here
			including the commodity market. The sector excludes the field-works for the construction of the
			bearing surface (eg construction of terraces in steeply sloping areas, drainage, preparation of
			rice fields), which belong to Section F (Construction - National economy section), nor does it
			include cooperative activities for the purchase and distribution of agricultural products. These
			market, commercial activities belong to Section G (Trade, repair of motor vehicles - National
			economy Section). This sector does not include: - activities related to landscape protection,
			maintenance, landscaping, see 8130.<br><span class="szakagazatnem_magyar">It does not belong to
				this branch - Hungarian supplement</span> This section excludes: management of green-areas,
			see 91.04</p>
	</div>
	`

	expected := `Definition / Belongs to this divisionDivision 01 first of all distinguishes two basic activities: production of crop products and production of animal products. This section comprises the organic farming, the cultivation of genetically modified plants and the breeding of genetically modified animals. This section contains the outdoor and greenhouse crop production. The Mixed farming subsection (number 01.5) does not follow the general principles of the identification of the main activity. It accepts that many agricultural holdings have reasonably balanced crop and animal production, and that it would be arbitrary to classify them in one category or the other.also belongs to this sector This group includes: The services of the agriculture and game farming activitiesAgricultural activity excludes any subsequent processing of the agricultural products (classified under division 10 and 11: Manufacture of food and beverage products, and division 12: Manufacture of tobacco products), beyond that needed to prepare them for the primary markets. So that the preparation of the crops for the primary market belongs also here including the commodity market. The sector excludes the field-works for the construction of the bearing surface (eg construction of terraces in steeply sloping areas, drainage, preparation of rice fields), which belong to Section F (Construction - National economy section), nor does it include cooperative activities for the purchase and distribution of agricultural products. These market, commercial activities belong to Section G (Trade, repair of motor vehicles - National economy Section). This sector does not include: - activities related to landscape protection, maintenance, landscaping, see 8130. This section excludes: management of green-areas, see 91.04`

	result, err := extractRelevantText(htmlContent)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Normalize whitespace for comparison
	expected = strings.Join(strings.Fields(expected), " ")
	result = strings.Join(strings.Fields(result), " ")

	if result != expected {
		t.Errorf("Extracted text does not match expected.\nExpected: %s\nGot: %s", expected, result)
	}
}
