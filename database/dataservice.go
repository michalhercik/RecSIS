package database

var db DB

func SetDatabase(newDB DB) {
    db = newDB
}

func GetBlueprintData(user int) BlueprintData {
    unassigned := db.GetUnassignedBlueprint(user)
    years := db.GetAssignedBlueprint(user)
    return BlueprintData{Unassigned: unassigned, Years: years}
}

func GetCoursesData() CoursesData {
    courses := db.GetAllCourses()
    return CoursesData{Courses: courses}
}

func RemoveFromBlueprint(user int, course int) {
    db.RemoveBlueprintCourse(user, course)
}

// TODO implement bussiness logic