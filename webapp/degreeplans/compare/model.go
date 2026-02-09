package compare

import "fmt"

//================================================================================
// Constants
//================================================================================

const (
	dpBaseCompare = "dpBaseCompare"
	dpCompareWith = "dpCompareWith"

	leftSide  = "left"
	rightSide = "right"

	mobileLayout  = "mobile"
	desktopLayout = "desktop"
)

//================================================================================
// Data Types and Methods
//================================================================================

type degreePlanComparePage struct {
	basePlan    degreePlanData
	comparePlan degreePlanData
}

type degreePlanData struct {
	code   string
	title  string
	blocks []degreePlanBlock
}

type degreePlanBlock struct {
	code         string
	title        string
	limit        int
	isCompulsory bool
	isOptional   bool
	courses      []course
}

type course struct {
	code        string
	title       string
	credits     int
	isSupported bool
	otherPlan   courseFlags
}

func (c *course) creditsString() string {
	if !c.isSupported {
		return ""
	}
	return fmt.Sprintf("%d", c.credits)
}

func (c *course) isInOtherPlan() bool {
	return c.otherPlan.isIn
}

func (c *course) notInOtherPlan() bool {
	return !c.otherPlan.isIn
}

func (c *course) inOtherPlanSameType() bool {
	return c.otherPlan.isIn && c.otherPlan.isSameType
}

func (c *course) inOtherPlanDifferentType() bool {
	return c.otherPlan.isIn && !c.otherPlan.isSameType
}

type courseFlags struct {
	isIn       bool
	isSameType bool
}
