package database

var db DB

func SetDatabase(newDB DB) {
    db = newDB
}

func GetCoursesData() CoursesData {
    courses := db.GetAllCourses()
    return CoursesData{Courses: courses}
}

func GetCourseData(id int) CourseData {
    course, er := db.GetCourse(id)
    if er != nil {
        panic(er)
    }
    return CourseData{Course: course}
}

func GetBlueprintData(user int) BlueprintData {
    unassigned := db.BlueprintGetUnassigned(user)
    years := db.BlueprintGetAssigned(user)
    return BlueprintData{Unassigned: unassigned, Years: years}
}

func BlueprintRemoveUnassigned(user int, course int) {
    db.BlueprintRemoveUnassigned(user, course)
}

func BlueprintRemoveYear(user int, year int) {
    yearCourses := db.BlueprintRemoveYear(user, year)
    for _, course := range yearCourses {
        db.BlueprintAddUnassigned(user, course.Id)
    }
}

func BlueprintAddYear(user int) {
    db.BlueprintAddYear(user)
}

// TODO implement business logic