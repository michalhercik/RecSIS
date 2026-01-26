package degreeplandetail

import (
	"strconv"

	"github.com/michalhercik/RecSIS/language"
)

type text struct {
	pageTitle             string
	studyField            string
	validity              string
	offcanvasMenu         string
	saveDegreePlan        string
	removeSavedDegreePlan string
	searchDegreePlans     string
	tabDetail             string
	tabRecommendedPlan    string
	tabRecommendedPlanSm  string
	tabReqMap             string
	recommendedPlan       string
	addRecToBPBtn         string
	addRecToBPTitle       string
	addRecToBPText        string
	mergeRecToBP          string
	clearBPaddRec         string
	showPrerequisites     string
	showCorequisites      string
	showIncompatibilities string
	legendTitle           string
	legendEdges           string
	legendNodes           string
	legendDirection       string
	prerequisite          string
	corequisite           string
	incompatibility       string
	legendNodeInPlan      string
	legendNodeOutsidePlan string
	code                  string
	title                 string
	status                string
	completed             string
	inBlueprint           string
	unassigned            string
	notCompleted          string
	credits               string
	creditsShort          string
	needed                string
	winter                string
	summer                string
	both                  string
	guarantors            string
	blueprint             string
	assign                string
	courseIsUnassigned    string
	year                  string
	winterAssign          string
	summerAssign          string
	// language
	language language.Language
	// errors
	errCannotGetUserDP        string
	errCannotGetDP            string
	errCannotSaveDP           string
	errCannotDeleteSavedDP    string
	errMissingMaxYearParam    string
	errInvalidMaxYearParam    string
	errCannotRewriteBlueprint string
	errCannotMergeToBlueprint string
	errDPNotFound             string
	errPageNotFound           string
	// tooltips
	ttUncheckAll       string
	ttAssignedCredits  string
	ttBlueprintCredits string
	ttCompletedCredits string
}

func (t text) yearStr(year int) string {
	switch t.language {
	case language.CS:
		return strconv.Itoa(year) + ". " + t.year
	case language.EN:
		return t.year + " " + strconv.Itoa(year)
	}
	return ""
}

var texts = map[language.Language]text{
	language.CS: {
		pageTitle:             "Studijní plán",
		studyField:            "Studijní obor",
		validity:              "Platnost",
		offcanvasMenu:         "Nabídka",
		saveDegreePlan:        "Uložit studijní plán",
		removeSavedDegreePlan: "Odstranit uložený studijní plán",
		searchDegreePlans:     "Vyhledat studijní plán",
		tabDetail:             "Detail plánu",
		tabRecommendedPlan:    "Doporučený průběh studia",
		tabRecommendedPlanSm:  "Dop. průběh",
		tabReqMap:             "Mapa předmětových závislostí",
		recommendedPlan:       "Doporučený průběh studia",
		addRecToBPBtn:         "Vložit doporučený průběh do blueprintu",
		addRecToBPTitle:       "Vložit doporučený průběh do blueprintu",
		addRecToBPText:        "Chcete přidat doporučený průběh studia do vašeho studijního plánu (blueprintu)? Můžete si vybrat, zda chcete kurzy přidat k již existujícím kurzům v blueprintu, nebo zda chcete nejprve vymazat stávající kurzy v blueprintu a poté vložit pouze kurzy z doporučeného průběhu.",
		mergeRecToBP:          "Sloučit s blueprintem",
		clearBPaddRec:         "Přepsat blueprint",
		showPrerequisites:     "Zobrazit prerekvizity",
		showCorequisites:      "Zobrazit korekvizity",
		showIncompatibilities: "Zobrazit neslučitelnosti",
		legendTitle:           "Legenda",
		legendEdges:           "Hrany",
		legendNodes:           "Uzly",
		legendDirection:       "Hrany vedou od kurzu k jeho prerekvizitě.",
		prerequisite:          "Prerekvizita",
		corequisite:           "Korekvizita",
		incompatibility:       "Neslučitelnost",
		legendNodeInPlan:      "V plánu",
		legendNodeOutsidePlan: "Mimo plán",
		code:                  "Kód",
		title:                 "Název",
		status:                "Stav",
		completed:             "Splněno",
		inBlueprint:           "Blueprint",
		unassigned:            "Nezařazen",
		notCompleted:          "Nesplněno",
		credits:               "Kredity",
		creditsShort:          "Kr.",
		needed:                "potřeba",
		winter:                "ZS",
		summer:                "LS",
		both:                  "Oba",
		guarantors:            "Garant(i)",
		blueprint:             "Blueprint",
		assign:                "Přiřadit",
		courseIsUnassigned:    "Kurz je v blueprintu, ale není zařazen.",
		year:                  "ročník",
		winterAssign:          "ZS",
		summerAssign:          "LS",
		// language
		language:           language.CS,
		errCannotGetUserDP: "Nebylo možné získat studijní plán uživatele",
		// errors
		errCannotGetDP:            "Nebylo možné získat vybraný studijní plán",
		errCannotSaveDP:           "Nebylo možné uložit studijní plán",
		errCannotDeleteSavedDP:    "Nebylo možné smazat uložený studijní plán",
		errMissingMaxYearParam:    "Chybí parametr maximálního ročníku",
		errInvalidMaxYearParam:    "Neplatný parametr maximálního ročníku",
		errCannotRewriteBlueprint: "Nebylo možné přepsat blueprint doporučeným průběhem",
		errCannotMergeToBlueprint: "Nebylo možné sloučit doporučený průběh s blueprintem",
		errDPNotFound:             "Studijní plán nenalezen",
		errPageNotFound:           "Stránka nenalezena",
		// tooltips
		ttUncheckAll:       "Zrušit zaškrtnutí všech kurzů",
		ttAssignedCredits:  "počet kreditů přiřazených do ročníků / limit skupiny",
		ttBlueprintCredits: "počet kreditů v blueprintu / limit skupiny",
		ttCompletedCredits: "počet splněných kreditů / limit skupiny",
	},
	language.EN: {
		pageTitle:             "Degree Plan",
		studyField:            "Study branch",
		validity:              "Validity",
		offcanvasMenu:         "Menu",
		saveDegreePlan:        "Save degree plan",
		removeSavedDegreePlan: "Remove saved degree plan",
		searchDegreePlans:     "Search degree plans",
		tabDetail:             "Plan detail",
		tabRecommendedPlan:    "Recommended study plan",
		tabRecommendedPlanSm:  "Rec. plan",
		tabReqMap:             "Map of course requisites",
		recommendedPlan:       "Recommended study plan",
		addRecToBPBtn:         "Insert recommended plan to blueprint",
		addRecToBPTitle:       "Insert recommended plan to blueprint",
		addRecToBPText:        "Do you want to insert the recommended study plan to your own study plan (blueprint)? You can choose whether to add the courses to the existing courses in the blueprint, or to first clear the existing courses in the blueprint and then add only the courses from the recommended plan.",
		mergeRecToBP:          "Merge with blueprint",
		clearBPaddRec:         "Overwrite blueprint",
		showPrerequisites:     "Show prerequisites",
		showCorequisites:      "Show corequisites",
		showIncompatibilities: "Show incompatibilities",
		legendTitle:           "Legend",
		legendEdges:           "Edges",
		legendNodes:           "Nodes",
		legendDirection:       "Edges point from a course to its requisite.",
		prerequisite:          "Prerequisite",
		corequisite:           "Corequisite",
		incompatibility:       "Incompatibility",
		legendNodeInPlan:      "In degree plan",
		legendNodeOutsidePlan: "Outside degree plan",
		code:                  "Code",
		title:                 "Title",
		status:                "Status",
		completed:             "Completed",
		inBlueprint:           "Blueprint",
		unassigned:            "Unassigned",
		notCompleted:          "Not completed",
		credits:               "Credits",
		creditsShort:          "Cr.",
		needed:                "needed",
		winter:                "Winter",
		both:                  "Both",
		summer:                "Summer",
		guarantors:            "Guarantor(s)",
		blueprint:             "Blueprint",
		assign:                "Assign",
		courseIsUnassigned:    "Course is in the blueprint but not assigned.",
		year:                  "Year",
		winterAssign:          "Winter",
		summerAssign:          "Summer",
		// language
		language: language.EN,
		// errors
		errCannotGetUserDP:        "Unable to retrieve user degree plan",
		errCannotGetDP:            "Unable to retrieve selected degree plan",
		errCannotSaveDP:           "Unable to save degree plan",
		errCannotDeleteSavedDP:    "Unable to delete saved degree plan",
		errMissingMaxYearParam:    "Missing maximum year parameter",
		errInvalidMaxYearParam:    "Invalid maximum year parameter",
		errCannotRewriteBlueprint: "Unable to overwrite blueprint with recommended plan",
		errCannotMergeToBlueprint: "Unable to merge recommended plan to blueprint",
		errDPNotFound:             "Degree plan not found",
		errPageNotFound:           "Page not found",
		// tooltips
		ttUncheckAll:       "Uncheck all courses",
		ttAssignedCredits:  "sum of credits assigned to years / group limit",
		ttBlueprintCredits: "sum of credits in the blueprint / group limit",
		ttCompletedCredits: "sum of completed credits / group limit",
	},
}
