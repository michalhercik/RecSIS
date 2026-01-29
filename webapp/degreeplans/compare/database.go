package compare

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/michalhercik/RecSIS/dbds"
	"github.com/michalhercik/RecSIS/degreeplans/compare/internal/sqlquery"
	"github.com/michalhercik/RecSIS/errorx"
	"github.com/michalhercik/RecSIS/language"
)

type DBManager struct {
	DB *sqlx.DB
}

type dbDegreePlanRecord struct {
	DegreePlanCode  string `db:"degree_plan_code"`
	DegreePlanTitle string `db:"degree_plan_title"`
	// DegreePlanValidFrom int            `db:"degree_plan_valid_from"`
	// DegreePlanValidTo   int            `db:"degree_plan_valid_to"`
	BlocCode  string `db:"bloc_subject_code"`
	BlocName  string `db:"bloc_name"`
	BlocLimit int    `db:"bloc_limit"`
	BlocType  string `db:"bloc_type"`
	dbds.Course
	CourseIsSupported bool `db:"course_is_supported"`
}

func (m DBManager) degreePlanCompareContent(basePlan, comparePlan string, lang language.Language) (*degreePlanComparePage, error) {
	var baseRecords []dbDegreePlanRecord
	err := m.DB.Select(&baseRecords, sqlquery.DegreePlan, basePlan, lang)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DegreePlan: %w", err), errorx.P("basePlan", basePlan), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetDP,
		)
	}
	if len(baseRecords) == 0 {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("degree plan not found: %s", basePlan), errorx.P("basePlan", basePlan)),
			http.StatusNotFound,
			texts[lang].errDPNotExisting,
		)
	}
	var compareRecords []dbDegreePlanRecord
	err = m.DB.Select(&compareRecords, sqlquery.DegreePlan, comparePlan, lang)
	if err != nil {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("sqlquery.DegreePlan: %w", err), errorx.P("comparePlan", comparePlan), errorx.P("lang", lang)),
			http.StatusInternalServerError,
			texts[lang].errCannotGetDP,
		)
	}
	if len(compareRecords) == 0 {
		return nil, errorx.NewHTTPErr(
			errorx.AddContext(fmt.Errorf("degree plan not found: %s", comparePlan), errorx.P("comparePlan", comparePlan)),
			http.StatusNotFound,
			texts[lang].errDPNotExisting,
		)
	}
	dp := buildDegreePlansComparePage(baseRecords, compareRecords)
	return &dp, nil
}

func buildDegreePlansComparePage(baseRecords, compareRecords []dbDegreePlanRecord) degreePlanComparePage {
	page := degreePlanComparePage{
		basePlan:    intoDegreePlan(baseRecords),
		comparePlan: intoDegreePlan(compareRecords),
	}
	addOtherPlanFlags(&page)
	return page
}

func intoDegreePlan(records []dbDegreePlanRecord) degreePlanData {
	var dp degreePlanData
	dp.code = records[0].DegreePlanCode
	dp.title = records[0].DegreePlanTitle
	// dp.validFrom = records[0].DegreePlanValidFrom
	// dp.validTo = records[0].DegreePlanValidTo
	for _, record := range records {
		add(&dp, record)
	}
	fixLimits(&dp)
	return dp
}

func add(dp *degreePlanData, record dbDegreePlanRecord) {
	blocIndex := -1
	for i, b := range dp.blocks {
		if b.code == record.BlocCode {
			blocIndex = i
			break
		}
	}
	if blocIndex == -1 {
		dp.blocks = append(dp.blocks, degreePlanBlock{
			title:        record.BlocName,
			code:         record.BlocCode,
			limit:        record.BlocLimit,
			isCompulsory: record.BlocType == "A",
			isOptional:   record.BlocType == "C",
		})
		blocIndex = len(dp.blocks) - 1
	}
	dp.blocks[blocIndex].courses = append(dp.blocks[blocIndex].courses, intoCourse(record))
}

func intoCourse(from dbDegreePlanRecord) course {
	return course{
		code:        from.Code,
		title:       from.Title,
		credits:     from.Credits,
		isSupported: from.CourseIsSupported,
	}
}

func fixLimits(dp *degreePlanData) {
	for i := range dp.blocks {
		if dp.blocks[i].isCompulsory {
			creditSum := 0
			for _, c := range dp.blocks[i].courses {
				creditSum += c.credits
			}
			dp.blocks[i].limit = creditSum
		} else if dp.blocks[i].isOptional {
			dp.blocks[i].limit = 0
		}
	}
}

func addOtherPlanFlags(page *degreePlanComparePage) {
	type course struct {
		code         string
		isCompulsory bool
		isOptional   bool
	}
	baseCourses := map[string]course{}
	for _, block := range page.basePlan.blocks {
		for _, c := range block.courses {
			baseCourses[c.code] = course{
				code:         c.code,
				isCompulsory: block.isCompulsory,
				isOptional:   block.isOptional,
			}
		}
	}
	compareCourses := map[string]course{}
	for _, block := range page.comparePlan.blocks {
		for _, c := range block.courses {
			compareCourses[c.code] = course{
				code:         c.code,
				isCompulsory: block.isCompulsory,
				isOptional:   block.isOptional,
			}
		}
	}
	for i, block := range page.basePlan.blocks {
		for j, c := range block.courses {
			if otherC, ok := compareCourses[c.code]; ok {
				page.basePlan.blocks[i].courses[j].otherPlan.isIn = true
				isSameType := (block.isCompulsory && otherC.isCompulsory) ||
					(block.isOptional && otherC.isOptional) ||
					(!block.isCompulsory && !block.isOptional && !otherC.isCompulsory && !otherC.isOptional)
				page.basePlan.blocks[i].courses[j].otherPlan.isSameType = isSameType
			}
		}
	}
	for i, block := range page.comparePlan.blocks {
		for j, c := range block.courses {
			if otherC, ok := baseCourses[c.code]; ok {
				page.comparePlan.blocks[i].courses[j].otherPlan.isIn = true
				isSameType := (block.isCompulsory && otherC.isCompulsory) ||
					(block.isOptional && otherC.isOptional) ||
					(!block.isCompulsory && !block.isOptional && !otherC.isCompulsory && !otherC.isOptional)
				page.comparePlan.blocks[i].courses[j].otherPlan.isSameType = isSameType
			}
		}
	}
}
