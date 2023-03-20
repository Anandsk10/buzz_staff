package dashboard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getVillages(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var villageQuery, subVillageQuery string
	var villageCount, subVillageCount int

	if len(projectArray) > 0 {

		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0 AND project_id IN (%s)%s", intsToString(projectArray), filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0 AND project_id IN (%s)%s", intsToString(projectArray), filter)
	} else {
		if len(startDate) > 0 {
			dateFilter := fmt.Sprintf(" AND date >= '%s' AND date <= '%s'", startDate, endDate)
			villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
			subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
		} else {
			dateFilter := " AND date >= CURRENT_DATE()"
			villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
			subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
		}
	}

	err := db.QueryRow(villageQuery).Scan(&villageCount)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
	if err != nil {
		log.Fatal(err)
	}

	return villageCount + subVillageCount
}

func GetNoOfVyaparSurvey(db *sql.DB, startDate string, endDate string, gfId string) int {
	var getActualsQuery string
	var noofvyaparsurvey int

	// if len(projectArray) > 0 {
	getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline")
	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline where entry_date BETWEEN '%s' AND '%s'", startDate, endDate)
	} else {
		if gfId != "" {
			getActualsQuery = fmt.Sprintf("SELECT COUNT(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline tp WHERE entry_date BETWEEN '%s' AND '%s' and gfid '%s'", startDate, endDate, gfId)

		}
	}
	//  } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
	// }
	//  }
	// } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline")
	// }
	// }
	err := db.QueryRow(getActualsQuery).Scan(&noofvyaparsurvey)
	if err != nil {
		log.Fatal(err)
	}
	return noofvyaparsurvey

}

func NoofVyaparCohorts(db *sql.DB, startDate string, endDate string, project_id string) int {
	var getActualsQuery string
	var noofvyaparcohorts int

	if startDate != "" && endDate != "" {
		getActualsQuery = fmt.Sprintf("select count(session_type) as noofvyaparcohorts from tbl_poa where type=2 and session_type=1 and date BETWEEN '%s' AND '%s'", startDate, endDate)
	} else if project_id != "" {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(session_type) as noofvyaparcohorts from tbl_poa where type=2 and session_type=1 and project_id '%s'", project_id)

	} else {
		getActualsQuery = fmt.Sprintf("SELECT COUNT(session_type) as noofvyaparcohorts from tbl_poa where type=2 and session_type=1")

	}
	err := db.QueryRow(getActualsQuery).Scan(&noofvyaparcohorts)
	if err != nil {
		log.Fatal(err)
	}
	return noofvyaparcohorts

}

func GetNoofVyaparModuleCompleted(db *sql.DB) int {
	var getActualsQuery string
	var noofvyaparmodulecompleted int

	// if len(projectArray) > 0 {
	// if startDate != "" && endDate != "" {
	getActualsQuery = fmt.Sprintf("select count(module1=1 and module2=1 and module3=1 and module4=1 and module5=1) from BuzzVyaparProgramBaseline")
	//  } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
	// }
	//  }
	// } else {
	//  getActualsQuery = fmt.Sprintf("select count(id) as noofvyaparsurvey from BuzzVyaparProgramBaseline")
	// }
	// }
	err := db.QueryRow(getActualsQuery).Scan(&noofvyaparmodulecompleted)
	if err != nil {
		log.Fatal(err)
	}
	return noofvyaparmodulecompleted

}

func getTarget(db *sql.DB, startDate string, endDate string, projectArray []int) int {
	var getTargetQuery string
	var target sql.NullInt64

	if len(projectArray) > 0 {
		getTargetQuery = fmt.Sprintf("SELECT COALESCE(sum(training_target), 0) as target from project p where id in (%s)", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	} else {
		getTargetQuery = fmt.Sprintf("SELECT sum(training_target) as target from project p where (startDate >= '%s' and endDate <= '%s')", startDate, endDate)
	}

	// ...

	err := db.QueryRow(getTargetQuery).Scan(&target)
	if err != nil {
		log.Fatal(err)
	}

	if target.Valid {
		return int(target.Int64)
	} else {
		return 0
	}

}

func getActual(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and (participant_day2 BETWEEN '%s' and '%s' )and project_id in (%s) %s", startDate, endDate, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		} else {
			getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), filter)
		}
	} else {
		getActualsQuery = fmt.Sprintf("select count(tp.id) as actual from training_participants tp inner join project p on p.id = tp.project_id where day2 = 1 and startDate >= '%s' and endDate <= '%s'", startDate, endDate)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Fatal(err)
	}
	return actual
}

func getParticipantFilterActual(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var getActualsQuery string
	var actual int

	if len(projectArray) > 0 {
		getActualsQuery = fmt.Sprintf("select count(id) as actual from training_participants tp where day2 = 1 and project_id in (%s) and participant_day2 >= '%s' and participant_day2 <= '%s' %s", strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"), startDate, endDate, filter)
	} else {
		getActualsQuery = fmt.Sprintf("select count(tp.id) as actual from training_participants tp inner join project p on p.id = tp.project_id where day2 = 1 and participant_day2 >= '%s' and participant_day2 <= '%s' %s", startDate, endDate, filter)
	}

	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Fatal(err)
	}
	return actual
}

func getDay1Count(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var day1CountQuery string
	if len(projectArray) > 0 {
		projectIDs := make([]string, len(projectArray))
		for i, projectID := range projectArray {
			projectIDs[i] = strconv.Itoa(projectID)
		}
		day1CountQuery = "select count(id) as day1Count from training_participants tp where day1 = 1 and project_id in (" + strings.Join(projectIDs, ",") + ")" + filter
	} else {
		day1CountQuery = "select count(tp.id) as day1Count from training_participants tp inner join project p on p.id = tp.project_id where day1 = 1 and startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
	}
	row := db.QueryRow(day1CountQuery)
	var day1Count int
	err := row.Scan(&day1Count)
	if err != nil {
		fmt.Println(err)
	}
	return day1Count
}

func getParticipantFilterDay1Count(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var day1CountQuery string
	if len(projectArray) > 0 {
		projectIDs := make([]string, len(projectArray))
		for i, projectID := range projectArray {
			projectIDs[i] = strconv.Itoa(projectID)
		}
		day1CountQuery = "select count(id) as day1Count from training_participants tp where day1 = 1 and project_id in (" + strings.Join(projectIDs, ",") + ") and  participant_day1 >= '" + startDate + "' and participant_day1 <= '" + endDate + "'" + filter
	} else {
		day1CountQuery = "select count(tp.id) as day1Count from training_participants tp inner join project p on p.id = tp.project_id where day1 = 1 and participant_day1 >= '" + startDate + "' and participant_day1 <= '" + endDate + "'" + filter
	}
	row := db.QueryRow(day1CountQuery)
	var day1Count int
	err := row.Scan(&day1Count)
	if err != nil {
		fmt.Println(err)
	}
	return day1Count
}

func getGelathi(db *sql.DB, startDate string, endDate string, projectArray []int, GalathiId string, funderId string, filter string) int {
	var gelathiCount int
	// get associated projects for each project
	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectArray = uniqueIntSlice(projectArray)
	// build query
	var gelatiCountQuery string
	if len(projectArray) > 0 {

		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as gelathiCount FROM training_participants tp WHERE enroll = 1 AND enrolledProject IN (%s) AND enroll_date BETWEEN '%s' AND '%s'", intSliceToString(projectArray), startDate, endDate)
		} else {
			var dateRangeQuery []string
			if funderId != "" {
				dateRangeQuery = []string{fmt.Sprintf("(SELECT min(startDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE())", funderId), fmt.Sprintf("(SELECT max(endDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE())", funderId)}
			} else {
				dateRangeQuery = []string{"(SELECT min(startDate) FROM project WHERE endDate >= CURRENT_DATE())", "(SELECT max(endDate) FROM project WHERE endDate >= CURRENT_DATE())"}
			}
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as gelathiCount FROM training_participants tp WHERE enroll = 1 AND enrolledProject IN (%s) AND enroll_date BETWEEN %s AND %s", intSliceToString(projectArray), dateRangeQuery[0], dateRangeQuery[1])
		}
		if GalathiId != "" {
			gelatiCountQuery += strings.Replace(GalathiId, "trainer_id", "gelathi_id", -1)
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as gelathiCount FROM training_participants tp INNER JOIN project p ON p.id = tp.enrolledProject WHERE enroll = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}
	// execute query
	row := db.QueryRow(gelatiCountQuery)
	row.Scan(&gelathiCount)
	return gelathiCount
}

// helper function to remove duplicates from an int slice
func uniqueIntSlice(slice []int) []int {
	keys := make(map[int]bool)
	var unique []int
	for _, val := range slice {
		if _, value := keys[val]; !value {
			keys[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}

// helper function to convert an int slice to a string of comma-separated values
func intSliceToString(slice []int) string {
	var strSlice []string
	for _, val := range slice {
		strSlice = append(strSlice, strconv.Itoa(val))
	}
	return strings.Join(strSlice, ", ")
}

func getParticipantFilterGelathi(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var gelathiCount int
	var gelathiCountQuery string
	if len(projectArray) > 0 {
		gelathiCountQuery = fmt.Sprintf("SELECT COUNT(id) as gelathiCount FROM training_participants tp where enroll = 1 and enrolledProject in (%s) and (date(enroll_date) >= '%s' and date(enroll_date) <= '%s') %s", intsToString(projectArray), startDate, endDate, filter)
	} else {
		gelathiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as gelathiCount FROM training_participants tp inner join project p on p.id = tp.enrolledProject where enroll = 1 and (date(enroll_date) >= '%s' and date(enroll_date) <= '%s') %s", startDate, endDate, filter)
	}
	row := db.QueryRow(gelathiCountQuery)
	err := row.Scan(&gelathiCount)
	if err != nil {
		log.Fatal(err)
	}
	return gelathiCount
}

// func getGreenMotivators(db *sql.DB, startDate string, endDate string, projectArray []string, funderID string, filter string) int {
// 	if funderID != "" {
// 		getProj := "SELECT id, startDate, endDate from project p WHERE funderID = ?"
// 		projResult, err := db.Query(getProj, funderID)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer projResult.Close()

// 		for projResult.Next() {
// 			var id string
// 			var startDate string
// 			var endDate string
// 			err = projResult.Scan(&id, &startDate, &endDate)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			projectArray = append(projectArray, id)
// 		}
// 	}

// 	// Perform the required checks and filters on the projectArray

// 	// Form the SQL query
// 	var gelatiCountQuery string
// 	if len(projectArray) > 0 {
// 		if startDate != "" && endDate != "" {
// 			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) AS greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN ("+strings.Join(projectArray, ",")+") AND GreenMotivatorsDate BETWEEN '%s' AND '%s'", startDate, endDate)
// 			// Execute the query and get the result
// 			row := db.QueryRow(gelatiCountQuery, startDate, endDate)
// 			var count int
// 			err := row.Scan(&count)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			return count
// 		}

// 	}
// 	return 0
// }
func greenMotivators(con *sql.DB, startDate string, endDate string, projectArray []int, funderId string, filter string) int {
	if funderId != "" {
		getProj := fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %s", funderId)
		projResult, err := con.Query(getProj)
		if err != nil {
			panic(err)
		}
		defer projResult.Close()

		for projResult.Next() {
			var id int
			var startDate string
			var endDate string

			err = projResult.Scan(&id, &startDate, &endDate)
			if err != nil {
				panic(err)
			}

			projectArray = append(projectArray, id)
		}
	}

	for _, proj := range projectArray {
		projs, _ := getAssociatedProjectList(con, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectArray = uniqueIntSlice(projectArray)

	var gelatiCountQuery string
	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN (%s) AND GreenMotivatorsDate BETWEEN '%s' AND '%s'", sliceToString(projectArray), startDate, endDate)
		} else {
			if funderId != "" {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN (%s) AND GreenMotivatorsDate BETWEEN (SELECT MIN(startDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE()) AND (SELECT MAX(endDate) FROM project p WHERE funderID = %s AND endDate >= CURRENT_DATE())", sliceToString(projectArray), funderId, funderId)
			} else {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp WHERE GreenMotivators = 1 AND project_id IN (%s) AND GreenMotivatorsDate BETWEEN (SELECT MIN(startDate) FROM project p WHERE id IN (%s)) AND (SELECT MAX(endDate) FROM project p WHERE id IN (%s))", sliceToString(projectArray), sliceToString(projectArray), sliceToString(projectArray))
			}
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as greenMoti FROM training_participants tp INNER JOIN project p ON p.id = tp.GreenMotivatorsEnrolledProject WHERE GreenMotivators = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}

	rows, err := con.Query(gelatiCountQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var greenMoti int
		err = rows.Scan(&greenMoti)
		if err != nil {
			panic(err)
		}

		return greenMoti
	}

	return 0
}

func sliceToString(slice []int) string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(strSlice, ", ")
}

func Vyapar(db *sql.DB, startDate string, endDate string, projectArray []int, funderId string, filter string) int {
	if funderId != "" {
		getProj := fmt.Sprintf("SELECT id, startDate, endDate FROM project p WHERE funderID = %s", funderId)
		projResult, _ := db.Query(getProj)
		for projResult.Next() {
			var id int
			var startDate string
			var endDate string
			projResult.Scan(&id, &startDate, &endDate)
			projectArray = append(projectArray, id)
		}
	}

	for _, proj := range projectArray {
		// check if there are any associated project for each project
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}

	projectArray = uniqueIntSlice(projectArray)
	var gelatiCountQuery string
	if len(projectArray) > 0 {
		if startDate != "" && endDate != "" {
			gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN '%s' AND '%s'", intSliceToString(projectArray), startDate, endDate)
		} else {
			if funderId != "" {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN (SELECT min(startDate) from project p where funderID = %s and endDate >= CURRENT_DATE()) and (SELECT max(endDate) from project p where funderID = %s and endDate >= CURRENT_DATE())", intSliceToString(projectArray), funderId, funderId)
			} else {
				gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp WHERE VyaparEnrollment = 1 AND project_id IN (%s) AND VyaparEnrollmentDate BETWEEN (SELECT min(startDate) from project p where id IN (%s)) and (SELECT max(endDate) from project p where id IN (%s))", intSliceToString(projectArray), intSliceToString(projectArray), intSliceToString(projectArray))
			}
		}
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as Vyapar FROM training_participants tp INNER JOIN project p ON p.id = tp.VyaparEnrollmentEnrolledProject WHERE VyaparEnrollment = 1 AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
	}

	rows, _ := db.Query(gelatiCountQuery)
	defer rows.Close()
	var VyaparCount int
	for rows.Next() {
		rows.Scan(&VyaparCount)
	}
	return VyaparCount
}

func ParticipantFiltergreenMotivators(db *sql.DB, startDate, endDate string, projectArray []int, filter string) (int, error) {
	var gelatiCountQuery string
	var greenMoti int
	var err error

	if len(projectArray) > 0 {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(id) as greenMoti FROM training_participants tp where GreenMotivators = 1 and GreenMotivatorsEnrolledProject in (%s) and (date(GreenMotivatorsDate) >= '%s' and date(GreenMotivatorsDate) <= '%s') %s", intsToString(projectArray), startDate, endDate, filter)
	} else {
		gelatiCountQuery = fmt.Sprintf("SELECT COUNT(tp.id) as greenMoti FROM training_participants tp inner join project p on p.id = tp.GreenMotivatorsEnrolledProject where GreenMotivators = 1 and startDate >= '%s' and endDate <= '%s'", startDate, endDate)
	}

	err = db.QueryRow(gelatiCountQuery).Scan(&greenMoti)
	if err != nil {
		return 0, err
	}

	return greenMoti, nil
}

func getParticipantFilterVyaparapar(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var vyaparCount int

	if len(projectArray) > 0 {
		vyaparCountQuery := fmt.Sprintf("SELECT COUNT(id) as Vyapar FROM training_participants tp where VyaparEnrollment = 1 and VyaparEnrollmentEnrolledProject in (%v) and (date(VyaparEnrollmentDate) >= '%v' and date(VyaparEnrollmentDate) <= '%v') %v", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]"), startDate, endDate, filter)
		err := db.QueryRow(vyaparCountQuery).Scan(&vyaparCount)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		vyaparCountQuery := fmt.Sprintf("SELECT COUNT(tp.id) as Vyapar FROM training_participants tp inner join project p on p.id = tp.VyaparEnrollmentEnrolledProject where VyaparEnrollment = 1 and startDate >= '%v' and endDate <= '%v'", startDate, endDate)
		err := db.QueryRow(vyaparCountQuery).Scan(&vyaparCount)
		if err != nil {
			log.Fatal(err)
		}
	}
	return vyaparCount
}

func getTrainingBatches(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	//previously considering training batch count for villages.
	/*if (count($projectArray) > 0){
	      $tbQuery = "SELECT count(*) as tbCount from
	                      (SELECT max(id) from tbl_poa tp where check_out is not null and `type`  = 1 and added  = 0 and project_id in (".implode(',',$projectArray).") $filter
	                      GROUP BY tb_id ) as tb";
	  }else{
	      $tbQuery = "SELECT count(*) as tbCount from
	                      (SELECT max(tp.id) from tbl_poa tp
	                              inner join project p on p.id = tp.project_id
	                              where check_out is not null and `type` = 1 and added = 0 and startDate >= '$startDate' and endDate <= '$endDate'
	                      GROUP BY tb_id ) as tb ";
	  }
	  return mysqli_query($db,$tbQuery)->fetch_assoc()['tbCount'];*/

	for _, proj := range projectArray {
		//check if there are any associated project for each project
		projs, _ := getAssociatedProjectList(db, proj)
		if len(projs) > 1 {
			projectArray = append(projectArray, projs...)
		}
	}
	projectArray = removeDuplicates(projectArray)

	//get the count of villages
	var villageQuery string
	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as villages from tbl_poa tp where check_out is not null and primary_id != tb_id and `type`  = 1 and added  = 0 and project_id in (%s) %s", intSliceToString(projectArray), filter)
	} else {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as villages from tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null and primary_id != tb_id and `type`  = 1 and added  = 0 and startDate >= '%s' and endDate <= '%s' %s", startDate, endDate, filter)
	}
	// fmt.Println(villageQuery)
	var villages int
	err := db.QueryRow(villageQuery).Scan(&villages)
	if err != nil {
		log.Fatal(err)
	}
	return villages
}

func removeDuplicates(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func newVillageCount(db *sql.DB, startDate string, endDate string, projectArray []string, filter string) int {
	// Create a variable to store the project IDs
	var projIDs []string

	// Iterate over the projectArray and get all associated projects
	for _, proj := range projectArray {
		projInt, err := strconv.Atoi(proj)
		if err != nil {
			return 0
		}
		projs, _ := getAssociatedProjectList(db, projInt)
		if len(projs) > 1 {
			projIDs = append(projIDs, strconv.Itoa(projInt))
		}
	}

	// Add the original project IDs to the list and remove any duplicates
	projIDs = append(projIDs, projectArray...)
	projIDs = unique(projIDs)

	// Construct the village query
	villageQuery := "SELECT COUNT(DISTINCT location_id) as village FROM tbl_poa tp WHERE check_out IS NOT NULL AND sub_village='' AND tb_id != primary_id AND type = 1 AND added = 0 AND project_id IN (" + strings.Join(projIDs, ",") + ")"
	subVillageQuery := "SELECT COUNT(DISTINCT sub_village) as subVillage FROM tbl_poa tp WHERE check_out IS NOT NULL AND sub_village!='' AND tb_id != primary_id AND type = 1 AND added = 0 AND project_id IN (" + strings.Join(projIDs, ",") + ")"

	// Add the date filter to the queries if provided
	if startDate != "" && endDate != "" {
		villageQuery += " AND date BETWEEN '" + startDate + "' AND '" + endDate + "'"
		subVillageQuery += " AND date BETWEEN '" + startDate + "' AND '" + endDate + "'"
	}

	// Execute the queries and get the village and sub-village counts
	var villageCount int
	var subVillageCount int
	err := db.QueryRow(villageQuery).Scan(&villageCount)
	if err != nil {
		return 0
	}
	err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
	if err != nil {
		return 0
	}

	// Return the total village count
	return villageCount + subVillageCount
}

// Helper function to remove duplicates from a string slice
func unique(slice []string) []string {
	seen := make(map[string]bool)
	j := 0
	for _, val := range slice {
		if seen[val] {
			continue
		}
		seen[val] = true
		slice[j] = val
		j++
	}
	return slice[:j]
}

func getParticipantFilterTrainingBatches(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, trainerId int) int {
	if trainerId > 0 {
		filter = " and tp.user_id = " + strconv.Itoa(trainerId)
	}

	var villageQuery string
	if len(projectArray) > 0 {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as villages  from tbl_poa tp where check_out is not null and primary_id != tb_id and `type`  = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "' and project_id in (" + intSliceToString(projectArray) + ")" + filter
	} else {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as villages  from tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null and primary_id != tb_id and `type`  = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
	}

	row := db.QueryRow(villageQuery)

	var villages int
	err := row.Scan(&villages)
	if err != nil {
		log.Fatal(err)
	}

	return villages
}

func getParticipantFilterTrainingBatchesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, trainerId int) int {
	if trainerId > 0 {
		filter = " and tp.user_id = " + strconv.Itoa(trainerId)
	}
	var villageQuery, subVillageQuery string
	if len(projectArray) > 0 {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp where check_out is not null AND sub_village='' and tb_id != primary_id and `type` = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp where check_out is not null AND sub_village!='' and tb_id != primary_id and `type` = 1 and added  = 0 and tp.date >= '" + startDate + "' and tp.date <= '" + endDate + " 23:59:59' AND project_id in (" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(projectArray)), ","), "[]") + ")" + filter
	} else {
		villageQuery = "SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp  inner join project p on p.id = tp.project_id where check_out is not null AND sub_village='' and tb_id != primary_id and `type` = 1 and added  = 0 tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
		subVillageQuery = "SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp  inner join project p on p.id = tp.project_id where check_out is not null AND sub_village!='' and tb_id != primary_id and `type` = 1 and added  = 0 tp.date >= '" + startDate + "' and tp.date <= '" + endDate + "'" + filter
	}
	villageResult, err := db.Query(villageQuery)
	if err != nil {
		panic(err.Error())
	}
	subVillageResult, err := db.Query(subVillageQuery)
	if err != nil {
		panic(err.Error())
	}
	var village, subVillage int
	for villageResult.Next() {
		err := villageResult.Scan(&village)
		if err != nil {
			panic(err.Error())
		}
	}
	for subVillageResult.Next() {
		err := subVillageResult.Scan(&subVillage)
		if err != nil {
			panic(err.Error())
		}
	}
	return village + subVillage
}

func getSummaryOfVillages(db *sql.DB, startDate, endDate string, projectArray []int, filter string) int {
	var villageQuery string

	if len(projectArray) > 0 {
		projectIDs := make([]string, len(projectArray))
		for i, v := range projectArray {
			projectIDs[i] = fmt.Sprintf("%d", v)
		}
		projectIDsStr := strings.Join(projectIDs, ",")

		villageQuery = fmt.Sprintf(`
            SELECT COUNT(DISTINCT location_id) as villageCount
            FROM tbl_poa tp
            INNER JOIN project p ON p.id = tp.project_id 
            WHERE check_out IS NOT NULL
            AND primary_id != tb_id
            AND type = 1
            AND added = 0
            AND project_id IN (%s)
            %s`,
			projectIDsStr, filter)
	} else {
		var dateFilter string
		if len(startDate) > 0 {
			dateFilter = fmt.Sprintf("AND startDate >= '%s' AND endDate <= '%s'", startDate, endDate)
		} else {
			dateFilter = "AND endDate >= CURRENT_DATE()"
		}

		villageQuery = fmt.Sprintf(`
            SELECT COUNT(DISTINCT location_id) as villageCount
            FROM tbl_poa tp 
            INNER JOIN project p ON p.id = tp.project_id
            WHERE check_out IS NOT NULL
            AND primary_id != tb_id
            AND type = 1
            AND added = 0
            %s
            %s`,
			dateFilter, filter)
	}

	var villageCount int
	err := db.QueryRow(villageQuery).Scan(&villageCount)
	if err != nil {
		// handle error
	}

	return villageCount
}

func getSummaryOfVillagesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string) int {
	var villageQuery, subVillageQuery string
	var villageCount, subVillageCount int

	if len(projectArray) > 0 {
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0 AND project_id IN (%s)%s", intsToString(projectArray), filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0 AND project_id IN (%s)%s", intsToString(projectArray), filter)
	} else {
		if len(startDate) > 0 {
			dateFilter := fmt.Sprintf(" AND date >= '%s' AND date <= '%s'", startDate, endDate)
			villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
			subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
		} else {
			dateFilter := " AND date >= CURRENT_DATE()"
			villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
			subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE check_out IS NOT NULL AND sub_village!='' AND `type` = 1 AND added = 0%s%s", dateFilter, filter)
		}
	}

	err := db.QueryRow(villageQuery).Scan(&villageCount)
	if err != nil {
		log.Fatal(err)
	}
	err = db.QueryRow(subVillageQuery).Scan(&subVillageCount)
	if err != nil {
		log.Fatal(err)
	}

	return villageCount + subVillageCount
}

func getParticipantFilterSummaryOfVillagesNew(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, trainerId int) int {
	if trainerId > 0 {
		filter = fmt.Sprintf(" and tp.user_id = %d", trainerId)
	}
	var villageQuery, subVillageQuery string
	if len(projectArray) > 0 {
		dateFilter := ""
		if startDate != "" && endDate != "" {
			dateFilter = fmt.Sprintf(" and (tp.date >= '%s' and tp.date <= '%s 23:59:59')", startDate, endDate)
		}
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null AND sub_village='' and type = 1 and added = 0 and project_id in (%s) %s %s", intsToString(projectArray), dateFilter, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null AND sub_village!='' and type = 1 and added = 0 and project_id in (%s) %s %s", intsToString(projectArray), dateFilter, filter)
	} else {
		dateFilter := fmt.Sprintf(" and (tp.date >= '%s' and tp.date <= '%s 23:59:59')", startDate, endDate)
		villageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT location_id) as 'village' FROM tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null AND sub_village='' and type = 1 and added = 0 %s %s", dateFilter, filter)
		subVillageQuery = fmt.Sprintf("SELECT COUNT(DISTINCT sub_village) as 'subVillage' FROM tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null AND sub_village!='' and type = 1 and added = 0 %s %s", dateFilter, filter)
	}
	villageResult, _ := db.Query(villageQuery)
	subVillageResult, _ := db.Query(subVillageQuery)
	var villageCount, subVillageCount int
	villageResult.Next()
	villageResult.Scan(&villageCount)
	subVillageResult.Next()
	subVillageResult.Scan(&subVillageCount)
	return villageCount + subVillageCount
}

func intsToString(ints []int) string {
	var stringSlice []string
	for _, i := range ints {
		stringSlice = append(stringSlice, strconv.Itoa(i))
	}
	return strings.Join(stringSlice, ",")
}

func getParticipantFilterSummaryOfVillages(db *sql.DB, startDate string, endDate string, projectArray []int, filter string, trainerId int) int {
	if trainerId > 0 {
		filter = fmt.Sprintf(" and tp.user_id = %d", trainerId)
	}
	if len(projectArray) > 0 {
		villageQuery := fmt.Sprintf("SELECT count(DISTINCT location_id) as villageCount from tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null and primary_id != tb_id and type = 1 and added = 0 and project_id in (%s)%s", intsToString(projectArray), filter)
		var villageCount int
		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			log.Fatal(err)
		}
		return villageCount
	} else {
		dateFilter := ""
		if len(startDate) > 0 {
			dateFilter = fmt.Sprintf(" and (tp.date >= '%s' and tp.date <= '%s')", startDate, endDate)
		} else {
			dateFilter = " and tp.date >= CURRENT_DATE()"
		}
		villageQuery := fmt.Sprintf("SELECT count(DISTINCT location_id) as villageCount from tbl_poa tp inner join project p on p.id = tp.project_id where check_out is not null and primary_id != tb_id and type = 1 and added = 0%s%s", dateFilter, filter)
		var villageCount int
		err := db.QueryRow(villageQuery).Scan(&villageCount)
		if err != nil {
			log.Fatal(err)
		}
		return villageCount
	}
}

func getOpsManagers(db *sql.DB, empId int) []int {
	getOpsIds := fmt.Sprintf("SELECT id FROM employee WHERE supervisorId = %d AND empRole = 4", empId)
	ids := make([]int, 0)
	res, err := db.Query(getOpsIds)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		if err := res.Scan(&id); err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	return ids
}

func getSupervisor(db *sql.DB, empId int) []int {
	getOpsIds := fmt.Sprintf("SELECT supervisorId FROM employee WHERE id = %d", empId)
	ids := make([]int, 0)
	res, err := db.Query(getOpsIds)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var supervisorId int
		if err := res.Scan(&supervisorId); err != nil {
			log.Fatal(err)
		}
		ids = append(ids, supervisorId)
	}
	return ids
}

func getReportingOpsManagers(db *sql.DB, empId int) []int {
	ids := []int{}

	getOpsIds := "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
	res, err := db.Query(getOpsIds, empId)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}

	getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 12"
	res, err = db.Query(getOpsIds, empId)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()

	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
		som := id
		getOpsIds = "SELECT id FROM employee WHERE supervisorId = ? AND empRole = 4"
		res, err = db.Query(getOpsIds, som)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Close()

		for res.Next() {
			var id int
			err := res.Scan(&id)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, id)
		}
	}

	return ids
}

func getOpProjects(db *sql.DB, empID int) []int {
	getProjIds := fmt.Sprintf("SELECT id FROM project WHERE operations_manager = %d GROUP BY id", empID)
	rows, err := db.Query(getProjIds)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ids
}

func getOpParticipantFilterProjects(db *sql.DB, empID int) []int {
	getProjIds := fmt.Sprintf("SELECT id FROM project WHERE operations_manager = %d GROUP BY id", empID)
	rows, err := db.Query(getProjIds)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return ids
}

func getTrainerTarget(db *sql.DB, empId int, projectArray []int) int {
	targetQuery := fmt.Sprintf("SELECT sum(target) as total from project_emps pe where emp_id = %d and project_id in (%s)",
		empId, strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	var total int
	err := db.QueryRow(targetQuery).Scan(&total)
	if err != nil {
		log.Fatal(err)
	}
	return total
}

func getTrainerActual(db *sql.DB, empId int, projectArray []int) int {
	getActualsQuery := fmt.Sprintf("select count(tp.id) as actual from training_participants tp "+
		"inner join project p on p.id = tp.project_id "+
		"where day2 = 1 and tp.trainer_id = %d and tp.project_id in (%s)", empId,
		strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	var actual int
	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Fatal(err)
	}
	return actual
}

func getTrainerDay1(db *sql.DB, empId int, projectArray []int) int {
	getActualsQuery := fmt.Sprintf("select count(tp.id) as actual from training_participants tp "+
		"inner join project p on p.id = tp.project_id "+
		"where day1 = 1 and tp.trainer_id = %d and tp.project_id in (%s)", empId,
		strings.Trim(strings.Replace(fmt.Sprint(projectArray), " ", ",", -1), "[]"))
	var actual int
	err := db.QueryRow(getActualsQuery).Scan(&actual)
	if err != nil {
		log.Fatal(err)
	}
	return actual
}

func getGFData(db *sql.DB, filter string, sessionType int, empId int) int {
	filter += fmt.Sprintf(" and tp.user_id = %d", empId)
	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) as visit from tbl_poa tp inner join project p on p.id = tp.project_id where `type` = 2 and session_type = %d AND tp.check_out is NOT NULL %s", sessionType, filter)
	var visit int
	row := db.QueryRow(getVisit)
	err := row.Scan(&visit)
	if err != nil {
		log.Fatal(err)
	}
	return visit
}

func getGFDataN(db *sql.DB, filter string, sessionType int, empId []int) int {
	filter += fmt.Sprintf(" and tp.user_id in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(empId)), ","), "[]"))
	getVisit := fmt.Sprintf("SELECT COUNT(tp.tb_id) as visit from tbl_poa tp inner join project p on p.id = tp.project_id where `type` = 2 and session_type = %d AND tp.check_out is NOT NULL %s", sessionType, filter)
	var visit int
	row := db.QueryRow(getVisit)
	err := row.Scan(&visit)
	if err != nil {
		log.Fatal(err)
	}
	return visit
}

func getGFCircle(db *sql.DB, filter string, empId int) int {
	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id = %d%s", empId, filter)
	row := db.QueryRow(getCircle)
	var visit int
	row.Scan(&visit)
	return visit
}

func getGFCircleN(db *sql.DB, filter string, empId []int) int {
	getCircle := fmt.Sprintf("SELECT COUNT(*) as visit from circle tp inner join project p on p.id = tp.project_id where tp.gelathi_created_id in (%s)%s", strings.Trim(strings.Replace(fmt.Sprint(empId), " ", ",", -1), "[]"), filter)
	row := db.QueryRow(getCircle)
	var visit int
	row.Scan(&visit)
	return visit
}

func getGfEnrolled(db *sql.DB, filter string, empID int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(tp.id) as enrolled FROM training_participants tp "+
		"INNER JOIN project p ON tp.project_id = p.id "+
		"WHERE enroll = 1 AND gelathi_id = %d %s", empID, filter)

	row := db.QueryRow(query)
	var enrolled int
	err := row.Scan(&enrolled)
	if err != nil {
		return 0, err
	}
	return enrolled, nil
}

func getGfEnrolledN(db *sql.DB, filter string, empId []int) int {
	empIdStr := make([]string, len(empId))
	for i, id := range empId {
		empIdStr[i] = strconv.Itoa(id)
	}
	query := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
		"inner join project p on tp.project_id = p.id " +
		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " + filter
	row := db.QueryRow(query)
	var enrolled int
	row.Scan(&enrolled)
	return enrolled
}

func getParticipantFilterGfEnrolled(db *sql.DB, filter string, empId int, startDate string, endDate string) int {
	getEnrolled := fmt.Sprintf("SELECT COUNT(tp.id) as enrolled from training_participants tp "+
		"inner join project p on tp.project_id = p.id "+
		"where enroll = 1 and gelathi_id = %d and "+
		"((participant_day1 >= '%s' and participant_day1 <= '%s') "+
		"or (participant_day2 >= '%s' and participant_day2 <= '%s')) %s", empId, startDate, endDate, startDate, endDate, filter)
	row := db.QueryRow(getEnrolled)
	var enrolled int
	err := row.Scan(&enrolled)
	if err != nil {
		log.Fatal(err)
	}
	return enrolled
}

func getParticipantFilterGfEnrolledN(db *sql.DB, filter string, empId []int, startDate, endDate time.Time) int {
	enrolled := 0
	empIdStr := make([]string, len(empId))
	for i, v := range empId {
		empIdStr[i] = strconv.Itoa(v)
	}
	getEnrolled := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
		"inner join project p on tp.project_id = p.id " +
		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " +
		"and ((participant_day1 >= ? and participant_day1 <= ?) or (participant_day2 >= ? and participant_day2 <= ?)) " + filter
	row := db.QueryRow(getEnrolled, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	row.Scan(&enrolled)
	return enrolled
}

func showNoProj() {
	data := make([]interface{}, 0)
	response := make(map[string]interface{})
	response["summary_target"] = 0
	response["summary_women"] = 0
	response["summary_villages"] = 0
	response["summary_actual"] = 0
	response["summary_day2"] = 0
	response["summary_enrolled"] = 0
	response["summary_green"] = 0
	response["summary_vyapar"] = 0
	response["data"] = data
	response["code"] = 200
	response["success"] = true
	response["message"] = "Successfully"

	json.NewEncoder(os.Stdout).Encode(response)
	os.Exit(0)
}

func Kann(db *sql.DB, filter string, empId []int, startDate, endDate time.Time) int {
	enrolled := 0
	empIdStr := make([]string, len(empId))
	for i, v := range empId {
		empIdStr[i] = strconv.Itoa(v)
	}
	getEnrolled := "SELECT COUNT(tp.id) as enrolled from training_participants tp " +
		"inner join project p on tp.project_id = p.id " +
		"where enroll = 1 and gelathi_id in (" + strings.Join(empIdStr, ",") + ") " +
		"and ((participant_day1 >= ? and participant_day1 <= ?) or (participant_day2 >= ? and participant_day2 <= ?)) " + filter
	row := db.QueryRow(getEnrolled, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	row.Scan(&enrolled)
	return enrolled
}
