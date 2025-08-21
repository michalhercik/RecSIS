package blueprint

import (
	"fmt"

	"github.com/michalhercik/RecSIS/language"
)

type Adapter interface {
	blueprint(userID string, lang language.Language) (*blueprintPage, error)
	moveCourses(userID string, lang language.Language, year int, semester semesterAssignment, position int, courses ...int) error
	appendCourses(userID string, lang language.Language, year int, semester semesterAssignment, courses ...int) error
	unassignSemester(userID string, lang language.Language, year int, semester semesterAssignment) error
	removeCourses(userID string, lang language.Language, courses ...int) error
	removeCoursesBySemester(userID string, lang language.Language, year int, semester semesterAssignment) error
	addYear(userID string, lang language.Language) error
	removeYear(userID string, lang language.Language, shouldUnassign bool) error
	foldSemester(userID string, lang language.Language, year int, semester semesterAssignment, folded bool) error
}

type Cache struct {
	Source     Adapter
	blueprints map[string]*blueprintPage
}

func (c Cache) blueprint(userID string, lang language.Language) (*blueprintPage, error) {
	if c.blueprints == nil {
		c.blueprints = make(map[string]*blueprintPage)
	}
	key := generateKey(userID, lang)
	bp, ok := c.blueprints[key]
	if ok {
		return bp, nil
	}
	bp, err := c.Source.blueprint(userID, lang)
	if err != nil {
		return nil, err
	}
	c.blueprints[key] = bp
	return bp, nil
}

func generateKey(userID string, lang language.Language) string {
	return string(fmt.Sprintf("%s:%s", userID, lang))
}

func (c Cache) invalidate(key string) {
	delete(c.blueprints, key)
}

func (c Cache) moveCourses(userID string, lang language.Language, year int, semester semesterAssignment, position int, courses ...int) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.moveCourses(userID, lang, year, semester, position, courses...)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) appendCourses(userID string, lang language.Language, year int, semester semesterAssignment, courses ...int) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.appendCourses(userID, lang, year, semester, courses...)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) unassignSemester(userID string, lang language.Language, year int, semester semesterAssignment) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.unassignSemester(userID, lang, year, semester)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) removeCourses(userID string, lang language.Language, courses ...int) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.removeCourses(userID, lang, courses...)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) removeCoursesBySemester(userID string, lang language.Language, year int, semester semesterAssignment) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.removeCoursesBySemester(userID, lang, year, semester)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) addYear(userID string, lang language.Language) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.addYear(userID, lang)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) removeYear(userID string, lang language.Language, shouldUnassign bool) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.removeYear(userID, lang, shouldUnassign)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) foldSemester(userID string, lang language.Language, year int, semester semesterAssignment, folded bool) error {
	key := generateKey(userID, lang)
	c.invalidate(key)
	err := c.Source.foldSemester(userID, lang, year, semester, folded)
	if err != nil {
		return err
	}
	return nil
}
