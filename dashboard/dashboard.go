package dashboard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/cors"
)

func HandleFunc() {
	db, err := sql.Open("mysql", "bdms_staff_admin:sfhakjfhyiqundfgs3765827635@tcp(buzzwomendatabase-new.cixgcssswxvx.ap-south-1.rds.amazonaws.com:3306)/bdms_staff?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database")
	defer db.Close()
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboard/vyapar", func(w http.ResponseWriter, r *http.Request) {
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		type ProjectRequest struct {
			PartnerID  int    `json:"partner_id"`
			Dist       int    `json:"dist"`
			Taluk      int    `json:"taluk"`
			Filter     string `json:"filter"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
			FunderId   int    `json:"funder_id"`
			ProjectID  int    `json:"project_id"`
			TrainerID  int    `json:"trainer_id"`
			OpsManager int    `json:"opsmanager"`
			SOMID      int    `json:"somid"`
			GFLID      int    `json:"gflid"`
			RoleID     int    `json:"roleid"`
			GalathiID  string `json:"galathi_id"`
			EmpID      int    `json:"emp_id"`
		}

		var reqBody ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing request body: %v", err)
			return
		}
		var summaryProjectsArray []interface{}
		data := []interface{}{}
		villagesArray := []interface{}{}
		summaryProjectsArray = []interface{}{}
		summaryTarget := 0
		summaryWomen := 0
		summaryVillages := 0
		summaryActuals := 0
		summaryDay1 := 0
		summaryDay2 := 0
		// summaryEnrolled := 0
		// summaryGreen := 0
		// summaryVyapar := 0

		if reqBody.RoleID == 1 || reqBody.RoleID == 9 || reqBody.RoleID == 3 || reqBody.RoleID == 4 || reqBody.RoleID == 12 {
			filter := ""
			summaryFilter := ""
			if reqBody.RoleID == 3 {
				var opsIds []int
				if reqBody.SOMID != 0 {
					opsIds = getReportingOpsManagers(db, reqBody.SOMID)
				} else if reqBody.GFLID != 0 {
					opsIds = getSupervisor(db, reqBody.GFLID)
				} else {
					opsIds = getReportingOpsManagers(db, reqBody.EmpID)
				}
				filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
			} else if reqBody.RoleID == 12 {
				opsIds := getOpsManagers(db, reqBody.EmpID)
				if len(opsIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					filter = " and p.operations_manager in (0)"
				}
			} else if reqBody.RoleID == 4 {
				// Ops Manager
				projectIds := getOpProjects(db, reqBody.EmpID)
				fmt.Println(projectIds)
				if len(projectIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager = %d", reqBody.EmpID)
				} else {
					showNoProj()
				}
			}

			isDateFilterApplied := false
			var dateFilter string
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				isDateFilterApplied = true
				dateFilter = fmt.Sprintf("startDate >= '%s' AND endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
			} else {
				dateFilter = "endDate >= CURRENT_DATE()"
			}

			var funderListQuery string

			if reqBody.PartnerID > 0 {
				funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) AS id, funderName AS name FROM project p "+
					"INNER JOIN funder ON funder.funderID = p.funderID "+
					"WHERE p.partnerID = %d AND %s %s", reqBody.PartnerID, dateFilter, filter)
				filter += fmt.Sprintf(" AND p.partnerID = %d", reqBody.PartnerID)
			} else if reqBody.Dist > 0 {
				if reqBody.Taluk > 0 {
					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID = %d AND %s %s GROUP BY p.funderID", reqBody.Taluk, dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID = %d", reqBody.Taluk)
				} else {
					// Get taluk of specified dist
					getTaluk := fmt.Sprintf("SELECT id FROM location WHERE `type` = 4 AND parentId = %d", reqBody.Dist)
					talukArray := []int{}
					talukRes, err := db.Query(getTaluk)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error getting taluk list: %v", err)
						return
					}
					defer talukRes.Close()

					for talukRes.Next() {
						var talukID int
						err := talukRes.Scan(&talukID)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							fmt.Fprintf(w, "Error scanning taluk list: %v", err)
							return
						}
						talukArray = append(talukArray, talukID)
					}

					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID IN (%s) AND %s %s GROUP BY p.funderID",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"),
						dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID IN (%s)",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"))
				}
			} else if reqBody.FunderId > 0 {
				funderListQuery = fmt.Sprintf("SELECT funderID as id, funderName as name FROM funder f WHERE funderID = %d", reqBody.FunderId)
				// summaryFilter := fmt.Sprintf(" AND p.funderID = %d", reqBody.FunderId)

			} else if reqBody.PartnerID == 0 && reqBody.TrainerID == 0 && reqBody.OpsManager == 0 && reqBody.SOMID == 0 && reqBody.GFLID == 0 && !isDateFilterApplied && reqBody.RoleID != 4 {
				// Role 4 OpsManager Default should be project list
				funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p INNER JOIN funder ON p.funderId = funder.funderID WHERE " + dateFilter + filter
			}

			funderList := []map[string]interface{}{}

			//gives funder list
			if len(funderListQuery) > 0 {
				res, err := db.Query(funderListQuery)
				if err != nil {
					// handle error
				}
				defer res.Close()
				for res.Next() {
					data := []interface{}{}

					projectArray := []int{}
					funderRow := map[string]interface{}{}
					var funderId int
					var funderName string
					err = res.Scan(&funderId, &funderName)
					if err != nil {
						// handle error
					}
					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter //commented by anas for ceo funder fiter not working
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id,startDate,endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}
					projResult, err := db.Query(getProj)
					if err != nil {
						// handle error
					}
					defer projResult.Close()
					for projResult.Next() {
						var projectId int
						err = projResult.Scan(&projectId)
						if err != nil {
							// handle error
						}
						projectArray = append(projectArray, projectId)
					}
					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":          funderId,
							"name":        funderName,
							"target":      0,
							"actual":      0,
							"day2":        0,
							"women":       0,
							"enrolled":    0,
							"villages":    0,
							"startDate":   "",
							"endDate":     "",
							"select_type": "2",
						}
						data = append(data, obj)
						continue
					}
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}
					obj := map[string]interface{}{
						"id":       funderId,
						"name":     funderName,
						"target":   getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":   getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"day2":     0,
						"women":    getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"enrolled": getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", ""),
						"villages": newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, ""), // New village count function anas

						"startDate":       "",
						"endDate":         "",
						"select_type":     "2",
						"greenMotivators": greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
						"vyapar":          Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
					}
					fmt.Println(funderRow)
					if day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""); day1Count > 0 {
						day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
						obj["day2"] = int(day2Turnout * 100)
					}
					data = append(data, obj)
				}
			}

			projectList := ""
			var summaryEnrolled, summaryGreen, summaryVyapar int

			// summaryFilter := ""

			if reqBody.ProjectID > 0 {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + strconv.Itoa(reqBody.ProjectID) + filter + dateFilterNew
				// summaryFilter := " and p.id = " + strdbv.Itoa(projectId)
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else if reqBody.TrainerID > 0 {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + strconv.Itoa(reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + strconv.Itoa(reqBody.TrainerID)
			} else if reqBody.OpsManager > 0 {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + strconv.Itoa(reqBody.OpsManager)
			} else if reqBody.SOMID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ")"
			} else if reqBody.GFLID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ")"
			} else if (isDateFilterApplied == true && reqBody.PartnerID == 0 && reqBody.Dist == 0 && reqBody.FunderId == 0) || (reqBody.RoleID == 4 && reqBody.Dist == 0) {
				//role 4 - OpsManager Default should be project list without location filter
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
			}
			fmt.Println(summaryFilter)

			if len(projectList) > 0 {
				res, err := db.Query(projectList)
				if err != nil {
					// handle error
				}
				defer res.Close()

				for res.Next() {
					var obj = make(map[string]interface{})
					var projectArray []int
					var id int
					var name string
					var startDate string
					var endDate string

					err := res.Scan(&id, &name, &startDate, &endDate)
					fmt.Println(err)
					if err != nil {
						fmt.Println(err)
					}

					obj["id"] = id
					obj["name"] = name

					projectArray = append(projectArray, id)

					var tpFilter string
					var tbFilter string

					if reqBody.TrainerID > 0 {
						target := getTrainerTarget(db, reqBody.TrainerID, projectArray)
						obj["target"] = target
						summaryTarget += target
						tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", reqBody.TrainerID)
						tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.TrainerID)
					} else {
						target := getTarget(db, startDate, endDate, projectArray)
						obj["target"] = target
						summaryTarget += target
					}

					actual := getActual(db, startDate, endDate, projectArray, tpFilter)
					obj["actual"] = actual
					summaryActuals += actual

					day1Count := getDay1Count(db, startDate, endDate, projectArray, tpFilter)
					summaryDay1 += day1Count

					if day1Count > 0 {
						day2Turnout := float64(actual) / float64(day1Count)
						obj["day2"] = int(math.Round(day2Turnout * 100))
					} else {
						obj["day2"] = 0
					}

					obj["women"] = actual
					obj["enrolled"] = getGelathi(db, startDate, endDate, projectArray, tpFilter, "", "")
					summaryEnrolled += obj["enrolled"].(int)
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}

					obj["villages"] = newVillageCount(db, startDate, endDate, strSlice, tbFilter)
					summaryVillages += obj["villages"].(int)

					obj["startDate"] = startDate
					obj["endDate"] = endDate
					obj["select_type"] = "1"

					obj["greenMotivators"] = greenMotivators(db, startDate, endDate, projectArray, tpFilter, "")
					obj["vyapar"] = Vyapar(db, startDate, endDate, projectArray, tpFilter, "")
					summaryGreen += obj["greenMotivators"].(int)
					summaryVyapar += obj["vyapar"].(int)

					data = append(data, obj)
					fmt.Println(data...)
					fmt.Println(summaryVillages)

					// data = append(data, obj)
					json.NewEncoder(w).Encode(map[string]interface{}{"No of Vypar Cohorts": NoofVyaparCohorts(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "No Of Vypar Enrolled Vypar": Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""), "No Of Vyapar survey": GetNoOfVyaparSurvey(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Vyapar module completed": GetNoofVyaparModuleCompleted(db), "funder": data})

				}
			}

			fmt.Println(summaryFilter)
			fmt.Println(funderList)
		} else if reqBody.RoleID == 5 {
			var dateFilter string
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if reqBody.ProjectID > 0 {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) + " and tp.project_id = " + strconv.Itoa(reqBody.ProjectID) +
					dateFilter +
					" GROUP by tp.project_id"
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := db.Query(query)

			if err != nil {
				log.Fatal(err)
			}
			var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Fatal(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				target := getTrainerTarget(db, reqBody.EmpID, projectArray)
				obj["target"] = target
				summaryTarget += target

				filter := " and tp.trainer_id = " + strconv.Itoa(reqBody.EmpID)
				actual := getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = actual
				summaryActuals += actual

				day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2TurnOut := float64(actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2TurnOut * 100))
				} else {
					obj["day2"] = 0
				}

				obj["women"] = obj["actual"]
				obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", "")
				var tbFilter string

				summaryEnrolled += obj["enrolled"].(int)
				tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = obj["startDate"]
				obj["endDate"] = obj["endDate"]
				obj["select_type"] = "1"
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			response := make(map[string]interface{})

			response["summary_target"] = summaryTarget
			response["summary_women"] = summaryActuals
			tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {
				// check if the element is of type int
				if i, ok := v.(int); ok {
					// append the int value to the []int slice
					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			var day2Turnout float64

			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				day2Turnout = 0
				response["summary_day2"] = 0
			}

			response["summary_enrolled"] = summaryEnrolled
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			// jsonResponse, err := json.Marshal(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"summary": response})
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(string(jsonResponse))

			// } else if reqBody.RoleID == 5 {
			// 	var dateFilter string
			// 	var isDateFilterApplied bool
			// 	if isDateFilterApplied {
			// 		dateFilter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate)
			// 	} else {
			// 		dateFilter = " and p.endDate >= CURRENT_DATE()"
			// 	}

			// 	// trainer
			// 	query := fmt.Sprintf("SELECT project_id as id,projectName as name,p.startDate,p.endDate "+
			// 		"from tbl_poa tp "+
			// 		"inner join project p on p.id = tp.project_id "+
			// 		"where user_id = %d %s GROUP by project_id", reqBody.EmpID, dateFilter)

			// 	if reqBody.ProjectID > 0 {
			// 		query = fmt.Sprintf("SELECT project_id as id,projectName as name,p.startDate,p.endDate "+
			// 			"from tbl_poa tp "+
			// 			"inner join project p on p.id = tp.project_id "+
			// 			"where user_id = %d and tp.project_id = %d GROUP by tp.project_id", reqBody.EmpID, reqBody.ProjectID)
			// 		summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			// 	}

			// 	res, err := db.Query(query)
			// 	if err != nil {
			// 		// handle error
			// 	}
			// 	defer res.Close()

			// 	var data []interface{}
			// 	summary := make(map[string]interface{})
			// 	var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int
			// 	for res.Next() {
			// 		var id, name string
			// 		var startDate, endDate string
			// 		err = res.Scan(&id, &name, &startDate, &endDate)
			// 		if err != nil {
			// 			// handle error
			// 		}

			// 		projectArray := []string{id}
			// 		obj := make(map[string]interface{})
			// 		obj["id"] = id
			// 		obj["name"] = name
			// 		obj["startDate"] = startDate
			// 		obj["endDate"] = endDate
			// 		obj["select_type"] = "1"

			// 		intSlice := make([]int, len(projectArray))

			// 		for i, str := range projectArray {
			// 			num, err := strconv.Atoi(str)
			// 			if err != nil {
			// 				panic(err)
			// 			}
			// 			intSlice[i] = num
			// 		}

			// 		obj["target"] = getTrainerTarget(db, reqBody.EmpID, intSlice)
			// 		if err != nil {
			// 			// handle error
			// 		}
			// 		summaryTarget += obj["target"].(int)

			// 		filter := fmt.Sprintf(" and tp.trainer_id = %d", reqBody.EmpID)
			// 		obj["actual"] = getActual(db, startDate, endDate, intSlice, filter)
			// 		if err != nil {
			// 			// handle error
			// 		}
			// 		summaryActuals += obj["actual"].(int)

			// 		day1Count := getDay1Count(db, startDate, endDate, intSlice, filter)
			// 		if err != nil {
			// 			// handle error
			// 		}
			// 		summaryDay1 += day1Count

			// 		if day1Count > 0 {
			// 			day2Turnout := obj["actual"].(float64) / float64(day1Count)
			// 			obj["day2"] = math.Round(day2Turnout * 100)
			// 		} else {
			// 			obj["day2"] = 0
			// 		}

			// 		obj["women"] = obj["actual"]
			// 		obj["enrolled"] = getGelathi(db, startDate, endDate, intSlice, "", "", filter)
			// 		if err != nil {
			// 			// handle error
			// 		}
			// 		summaryEnrolled += obj["enrolled"].(int)

			// 		tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			// 		obj["villages"] = newVillageCount(db, startDate, endDate, projectArray, tbFilter)
			// 		if err != nil {

			// 			// handle error
			// 		}
			// 		summaryVillages += obj["villages"].(int)
			// 		obj["green"] = greenMotivators(db, startDate, endDate, intSlice, filter, "")
			// 		if err != nil {
			// 			// handle error
			// 		}
			// 		summaryGreen += obj["green"].(int)

			// 		obj["vyapar"] = Vyapar(db, startDate, endDate, intSlice, filter, "")
			// 		if err != nil {
			// 			// handle error
			// 		}
			// 		summaryVyapar += obj["vyapar"].(int)

			// 		data = append(data, obj)
			// 		fmt.Println(data...)
			// 	}

			// 	summary["target"] = summaryTarget
			// 	summary["actual"] = summaryActuals
			// 	summary["day1"] = summaryDay1
			// 	summary["enrolled"] = summaryEnrolled
			// 	summary["villages"] = summaryVillages
			// 	summary["green"] = summaryGreen
			// 	summary["vyapar"] = summaryVyapar
			// 	json.NewEncoder(w).Encode(map[string]interface{}{"funder": data, "summary": summary})

			// return map[string]interface{}{"data": data, "summary": summary}

		} else if reqBody.RoleID == 6 {
			participantFilter := ""
			var filter, filterG, filterV string
			var isDateFilterApplied bool

			if reqBody.ProjectID > 0 {
				filter = fmt.Sprintf(" and tp.project_id = %d", reqBody.ProjectID)
			} else {
				if isDateFilterApplied {
					filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate >= '%s' and p.GreenMotivatorsDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate >= '%s' and p.VyaparEnrollmentDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
				} else {
					filter = " and p.endDate >= CURRENT_DATE()"
					filterG = " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
					filterV = " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
				}
			}

			circleMeet := getGFData(db, filter, 1, reqBody.EmpID)
			villageVisit := getGFData(db, filter, 2, reqBody.EmpID)
			beehive := getGFData(db, filter, 3, reqBody.EmpID)
			enrolled, _ := getGfEnrolled(db, filter, reqBody.EmpID)
			circleVisit := getGFCircle(db, filter, reqBody.EmpID)

			data := []map[string]interface{}{}
			getProjs := fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.user_id = %d %s GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d %s", reqBody.EmpID, filter, reqBody.EmpID, filter)

			if reqBody.ProjectID > 0 {
				getProjs = fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.project_id = %d and tp.user_id = %d GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d and tp.project_id = %d", reqBody.ProjectID, reqBody.EmpID, reqBody.EmpID, reqBody.ProjectID)
			}
			projectsList, err := db.Query(getProjs)
			if err != nil {
				fmt.Println(err)
			}
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					panic(err)
				}

				obj := make(map[string]interface{})
				obj["name"] = name
				// var villageProjvisit *int

				prjFilter := fmt.Sprintf(" and p.id = %d", id)
				circleProjMeet := getGFData(db, prjFilter, 1, reqBody.EmpID)
				obj["circle_meet"] = circleProjMeet
				obj["circles"] = getGFCircle(db, prjFilter, reqBody.EmpID)
				villageProjvisit := getGFData(db, prjFilter, 2, reqBody.EmpID)
				obj["villagevisit"] = villageProjvisit
				beehiveProj := getGFData(db, prjFilter, 3, reqBody.EmpID)
				obj["beehive"] = beehiveProj
				projEnrolled, _ := getGfEnrolled(db, prjFilter, reqBody.EmpID)
				bool1, _ := strconv.ParseBool(participantFilter)

				if bool1 {
					projEnrolled = getParticipantFilterGfEnrolled(db, prjFilter, reqBody.EmpID, reqBody.StartDate, reqBody.EndDate)
				}

				obj = make(map[string]interface{})
				obj["enroll"] = projEnrolled

				var projectArray []string
				str := strconv.Itoa(reqBody.EmpID)
				empIDs := strings.Split(str, ",")
				empIDsStr := "'" + strings.Join(empIDs, "','") + "'"
				projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", empIDsStr))

				// projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", reqBody.EmpID))
				if err != nil {
					// handle error
				}
				defer projectResult.Close()
				for projectResult.Next() {
					var prjIDs string
					err := projectResult.Scan(&prjIDs)
					if err != nil {
						// handle error
					}
					projectArray = append(projectArray, prjIDs)
				}
				intSlice := make([]int, len(projectArray))

				for i, str := range projectArray {
					num, err := strconv.Atoi(str)
					if err != nil {
						panic(err)
					}
					intSlice[i] = num
				}

				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterG)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterV)
				summaryGreen := 0
				summaryVyapar := 0

				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)

				data = append(data, obj)
				fmt.Println(data)

				response := make(map[string]interface{})
				response["summary_circle_meet"] = circleMeet
				response["summary_circles"] = circleVisit
				response["summary_villagevisit"] = villageVisit
				response["summary_beehive"] = beehive
				response["summary_enroll"] = enrolled
				response["summary_green"] = summaryGreen
				response["summary_vyapar"] = summaryVyapar
				response["data"] = data
				response["code"] = 200
				response["success"] = true
				response["message"] = "Successfully"

				js, err := json.Marshal(response)
				fmt.Println(response)
				if err != nil {
					// handle error
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}

			// } else if reqBody.RoleID == 13 {
			// 	filter := ""
			// 	filters := ""
			// 	filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			// 	filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"

			// 	if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 		filter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			// 		filters = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			// 		filterG = " and tp.GreenMotivatorsDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			// 		filterV = " and tp.VyaparEnrollmentDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			// 	}

			// 	if reqBody.ProjectID != 0 {
			// 		filters += " and p.id = " + reqBody.ProjectID
			// 	}

			// 	f := ""
			// 	if reqBody.GFLID != 0 {
			// 		f = " and id=" + reqBody.GFLID
			// 	}

			// 	summarycircleMeet := 0
			// 	summarycircles := 0
			// 	summaryvillagevisit := 0
			// 	summarybeehive := 0
			// 	summaryenroll := 0
			// 	summaryGreen := 0
			// 	summaryVyapar := 0
			// 	summarycircle_meet := 0

			// 	empIds := []string{}
			// 	rows, err := db.Query("SELECT id from employee e WHERE status =1 AND e.supervisorId = ?", reqBody.EmpID+f)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	defer rows.Close()

			// 	for rows.Next() {
			// 		var id string
			// 		if err := rows.Scan(&id); err != nil {
			// 			log.Fatal(err)
			// 		}
			// 		empIds = append(empIds, id)
			// 	}
			// 	if err := rows.Err(); err != nil {
			// 		log.Fatal(err)
			// 	}

			// 	getProjs := "Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where p.gfl_id = ? " + filters + " GROUP by tp.project_id"
			// 	rows, err = db.Query(getProjs, reqBody.EmpID)
			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	defer rows.Close()
			// 	var participantFilter string
			// 	projectsList, err := db.Query(getProjs)
			// 	if err != nil {
			// 		// handle error
			// 	}
			// 	var circleProjMeet string

			// 	for projectsList.Next() {
			// 		var obj map[string]interface{}
			// 		obj = make(map[string]interface{})
			// 		var id int
			// 		var name string
			// 		err = projectsList.Scan(&id, &name)
			// 		if err != nil {
			// 			// handle error
			// 		}

			// 		obj["name"] = name
			// 		var prjFilter string

			// 		if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 			prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", reqBody.StartDate, reqBody.EndDate, id)
			// 		} else {
			// 			prjFilter = fmt.Sprintf(" and p.id = %d", id)
			// 		}

			// 		// obj := make(map[string]interface{})
			// 		obj["circle_meet"] = "0"
			// 		if circleProjMeet != "" {
			// 			obj["circle_meet"] = circleProjMeet
			// 			summarycircle_meet += circleProjMeet
			// 		}

			// 		prjFilter = ""
			// 		if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 			prjFilter = fmt.Sprintf(" and p.endDate BETWEEN '%s' and '%s' and p.id = %d", reqBody.StartDate, reqBody.EndDate, id)
			// 		} else {
			// 			prjFilter = fmt.Sprintf(" and p.id = %d", id)
			// 		}
			// 		circles := getGFCircleN(db, prjFilter, ids)
			// 		obj["circles"] = circles
			// 		summarycircles += circles

			// 		prjFilter = ""
			// 		if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 			prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", reqBody.StartDate, reqBody.EndDate, id)
			// 		} else {
			// 			prjFilter = fmt.Sprintf(" and p.id = %d", id)
			// 		}
			// 		villageProjvisit := getGFDataN(db, prjFilter, 2, ids)
			// 		obj["villagevisit"] = "0"
			// 		if villageProjvisit != "" {
			// 			obj["villagevisit"] = villageProjvisit
			// 			summaryvillagevisit += villageProjvisit
			// 		}

			// 		prjFilter = ""
			// 		if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 			prjFilter = fmt.Sprintf(" and tp.date BETWEEN '%s' and '%s' and p.id = %d", reqBody.StartDate, reqBody.EndDate, id)
			// 		} else {
			// 			prjFilter = fmt.Sprintf(" and p.id = %d", id)
			// 		}
			// 		beehiveProj := getGFDataN(db, prjFilter, 3, ids)
			// 		obj["beehive"] = "0"
			// 		if beehiveProj != "" {
			// 			obj["beehive"] = beehiveProj
			// 			summarybeehive += beehiveProj
			// 		}

			// 		prjFilter = ""
			// 		if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 			prjFilter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s' and p.id = %d", reqBody.StartDate, reqBody.EndDate, id)
			// 		} else {
			// 			prjFilter = fmt.Sprintf(" and p.id = %d", id)
			// 		}
			// 		projEnrolled := getGfEnrolledN(db, prjFilter, ids)

			// 		if participantFilter {
			// 			if reqBody.StartDate != "" && reqBody.EndDate != "" {
			// 				prjFilter = fmt.Sprintf(" and tp.participant_day2 BETWEEN '%s' and '%s' and p.id = %d", reqBody.StartDate, reqBody.EndDate, id)
			// 			} else {
			// 				prjFilter = fmt.Sprintf(" and p.id = %d", id)
			// 			}
			// 			projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, reqBody.StartDate, reqBody.EndDate)
			// 		}
			// 		if projEnrolled != nil {
			// 			obj["enroll"] = *projEnrolled
			// 			summaryenroll += *projEnrolled
			// 		} else {
			// 			obj["enroll"] = "0"
			// 		}

			// 		obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, filterG)
			// 		obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, filterV)
			// 		summaryGreen += obj["greenMotivators"].(int)
			// 		summaryVyapar += obj["vyapar"].(int)
			// 		data = append(data, obj)
			// 		response := make(map[string]interface{})
			// 		response["summary_circle_meet"] = summarycircle_meet
			// 		response["summary_circles"] = summarycircles
			// 		response["summary_villagevisit"] = summaryvillagevisit
			// 		response["summary_beehive"] = summarybeehive
			// 		response["summary_enroll"] = summaryenroll
			// 		response["summary_green"] = summaryGreen
			// 		response["summary_vyapar"] = summaryVyapar
			// 		response["data"] = data
			// 		response["code"] = 200
			// 		response["success"] = true
			// 		response["message"] = "Successfully"

			// 		jsonResponse, err := json.Marshal(response)
			// 		if err != nil {
			// 			// handle error
			// 		}

			// 		fmt.Println(string(jsonResponse))
			// 	}

		} else if reqBody.RoleID == 13 {
			data := []map[string]interface{}{}
			filter := ""
			filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				filter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterG = " and tp.GreenMotivatorsDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterV = " and tp.VyaparEnrollmentDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			}
			// Additional code for roleId 13 can be added here
			filters := ""
			str := strconv.Itoa(reqBody.ProjectID)
			if str != "" {
				filters = " and p.id = " + str
			}
			// var request string

			// f := ""
			// if gfId, ok := request["gfId"]; ok && gfId != "" {
			// 	f = " and id=" + gfId
			// }

			summarycircleMeet := 0
			summarycircles := 0
			summaryvillagevisit := 0
			summarybeehive := 0
			summaryenroll := 0
			summaryGreen := 0
			summaryVyapar := 0
			summarycircle_meet := 0

			em, err := db.Query("SELECT id from employee e WHERE status =1 AND  e.supervisorId = ?", reqBody.EmpID)
			if err != nil {
				// handle error
			}
			ids := []int{}
			for em.Next() {
				var id int
				err := em.Scan(&id)
				if err != nil {
					// handle error
				}
				ids = append(ids, id)
			}

			getProjs := "Select project_id as id,p.projectName as name from tbl_poa tp " +
				"inner join project p on p.id = tp.project_id " +
				"where  p.gfl_id = ?" + filters + " GROUP by tp.project_id"

			projectsList, err := db.Query(getProjs, reqBody.EmpID)
			if err != nil {
				// handle error
			}
			participantFilter := ""
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					// handle error
				}

				var prjFilter string
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				circleProjMeet := getGFDataN(db, prjFilter, 1, ids)

				obj := make(map[string]interface{})
				obj["name"] = name
				obj["circle_meet"] = circleProjMeet
				summarycircleMeet += circleProjMeet

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and p.endDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				obj["circles"] = getGFCircleN(db, prjFilter, ids)
				summarycircles += obj["circles"].(int)

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var villageProjvisit interface{}
				villageProjvisit = getGFDataN(db, prjFilter, 2, ids)
				obj = make(map[string]interface{})
				obj["villagevisit"] = villageProjvisit
				if villageProjvisit != nil {
					summaryvillagevisit += villageProjvisit.(int)
				} else {
					obj["villagevisit"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var beehiveProj interface{}
				beehiveProj = getGFDataN(db, prjFilter, 3, ids)
				obj["beehive"] = beehiveProj
				if beehiveProj != nil {
					summarybeehive += beehiveProj.(int)
				} else {
					obj["beehive"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var projEnrolled interface{}
				b, err := strconv.ParseBool(participantFilter)
				reqBody.StartDate = "2023-03-15T15:30:45Z"

				// layout of the time string
				layout := ""

				// parse the time string into a time.Time object
				t, err := time.Parse(layout, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}
				reqBody.EndDate = ""

				// layout of the time string
				lay := "2006-01-02T15:04:05Z"

				// parse the time string into a time.Time object
				u, err := time.Parse(lay, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}

				// print the time.Time object
				fmt.Println("Parsed time:", u)
				if b {
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					} else {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					}
				} else {
					projEnrolled = getGfEnrolledN(db, prjFilter, ids)
				}
				obj["enroll"] = projEnrolled
				if projEnrolled != nil {
					summaryenroll += projEnrolled.(int)
				} else {
					obj["enroll"] = "0"
				}
				project_result, err := db.Query("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (" + strconv.Itoa(reqBody.EmpID) + ")")

				if err != nil {
					// handle error
				}
				var projectArray string
				if project_result.Next() {
					var ids interface{}
					err := project_result.Scan(&ids)
					if err != nil {
						// handle error
					}
					projectArray = ids.(string)
				}
				parts := strings.Split(projectArray, "")
				nums := make([]int, len(parts))
				for i, p := range parts {
					num, err := strconv.Atoi(p)
					if err != nil {
						panic(err)
					}
					nums[i] = num
				}
				// fmt.Println(nums)
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, nums, filterG, "")
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, nums, filterV, "")
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			fmt.Println(data)
			response := make(map[string]interface{})
			response["summary_circle_meet"] = summarycircle_meet
			response["summary_circles"] = summarycircles
			response["summary_villagevisit"] = summaryvillagevisit
			response["summary_beehive"] = summarybeehive
			response["summary_enroll"] = summaryenroll
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"
			fmt.Println(response)
			// json.NewEncoder(w).Encode(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"funder": response})

			fmt.Println(filter)

		} else {
			w.WriteHeader(http.StatusCreated)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid role id"
			json.NewEncoder(w).Encode(response)
		}
		fmt.Println(data...)
		fmt.Println(villagesArray...)
		fmt.Println(summaryWomen)
		fmt.Println(summaryDay2)

	})

	mux.HandleFunc("/dashboard/green", func(w http.ResponseWriter, r *http.Request) {
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		type ProjectRequest struct {
			PartnerID  int    `json:"partner_id"`
			Dist       int    `json:"dist"`
			Taluk      int    `json:"taluk"`
			Filter     string `json:"filter"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
			FunderId   int    `json:"funder_id"`
			ProjectID  int    `json:"project_id"`
			TrainerID  int    `json:"trainer_id"`
			OpsManager int    `json:"opsmanager"`
			SOMID      int    `json:"somid"`
			GFLID      int    `json:"gflid"`
			RoleID     int    `json:"roleid"`
			GalathiID  string `json:"galathi_id"`
			EmpID      int    `json:"emp_id"`
		}

		var reqBody ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing request body: %v", err)
			return
		}
		var summaryProjectsArray []interface{}
		data := []interface{}{}
		villagesArray := []interface{}{}
		summaryProjectsArray = []interface{}{}
		summaryTarget := 0
		summaryWomen := 0
		summaryVillages := 0
		summaryActuals := 0
		summaryDay1 := 0
		summaryDay2 := 0
		summaryEnrolled := 0
		summaryGreen := 0
		summaryVyapar := 0

		if reqBody.RoleID == 1 || reqBody.RoleID == 9 || reqBody.RoleID == 3 || reqBody.RoleID == 4 || reqBody.RoleID == 12 {
			filter := ""
			summaryFilter := ""
			if reqBody.RoleID == 3 {
				var opsIds []int
				if reqBody.SOMID != 0 {
					opsIds = getReportingOpsManagers(db, reqBody.SOMID)
				} else if reqBody.GFLID != 0 {
					opsIds = getSupervisor(db, reqBody.GFLID)
				} else {
					opsIds = getReportingOpsManagers(db, reqBody.EmpID)
				}
				filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
			} else if reqBody.RoleID == 12 {
				opsIds := getOpsManagers(db, reqBody.EmpID)
				if len(opsIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					filter = " and p.operations_manager in (0)"
				}
			} else if reqBody.RoleID == 4 {
				// Ops Manager
				projectIds := getOpProjects(db, reqBody.EmpID)
				fmt.Println(projectIds)
				if len(projectIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager = %d", reqBody.EmpID)
				} else {
					showNoProj()
				}
			}

			isDateFilterApplied := false
			var dateFilter string
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				isDateFilterApplied = true
				dateFilter = fmt.Sprintf("startDate >= '%s' AND endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
			} else {
				dateFilter = "endDate >= CURRENT_DATE()"
			}

			var funderListQuery string

			if reqBody.PartnerID > 0 {
				funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) AS id, funderName AS name FROM project p "+
					"INNER JOIN funder ON funder.funderID = p.funderID "+
					"WHERE p.partnerID = %d AND %s %s", reqBody.PartnerID, dateFilter, filter)
				filter += fmt.Sprintf(" AND p.partnerID = %d", reqBody.PartnerID)
			} else if reqBody.Dist > 0 {
				if reqBody.Taluk > 0 {
					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID = %d AND %s %s GROUP BY p.funderID", reqBody.Taluk, dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID = %d", reqBody.Taluk)
				} else {
					// Get taluk of specified dist
					getTaluk := fmt.Sprintf("SELECT id FROM location WHERE `type` = 4 AND parentId = %d", reqBody.Dist)
					talukArray := []int{}
					talukRes, err := db.Query(getTaluk)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error getting taluk list: %v", err)
						return
					}
					defer talukRes.Close()

					for talukRes.Next() {
						var talukID int
						err := talukRes.Scan(&talukID)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							fmt.Fprintf(w, "Error scanning taluk list: %v", err)
							return
						}
						talukArray = append(talukArray, talukID)
					}

					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID IN (%s) AND %s %s GROUP BY p.funderID",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"),
						dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID IN (%s)",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"))
				}
			} else if reqBody.FunderId > 0 {
				funderListQuery = fmt.Sprintf("SELECT funderID as id, funderName as name FROM funder f WHERE funderID = %d", reqBody.FunderId)
				// summaryFilter := fmt.Sprintf(" AND p.funderID = %d", reqBody.FunderId)

			} else if reqBody.PartnerID == 0 && reqBody.TrainerID == 0 && reqBody.OpsManager == 0 && reqBody.SOMID == 0 && reqBody.GFLID == 0 && !isDateFilterApplied && reqBody.RoleID != 4 {
				// Role 4 OpsManager Default should be project list
				funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p INNER JOIN funder ON p.funderId = funder.funderID WHERE " + dateFilter + filter
			}

			// funderList := []map[string]interface{}{}

			//gives funder list
			if len(funderListQuery) > 0 {
				res, err := db.Query(funderListQuery)
				if err != nil {
					// handle error
				}
				defer res.Close()
				for res.Next() {
					data := []interface{}{}

					projectArray := []int{}
					funderRow := map[string]interface{}{}
					var funderId int
					var funderName string
					err = res.Scan(&funderId, &funderName)
					if err != nil {
						// handle error
					}
					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter //commented by anas for ceo funder fiter not working
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id,startDate,endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}
					projResult, err := db.Query(getProj)
					if err != nil {
						// handle error
					}
					defer projResult.Close()
					for projResult.Next() {
						var projectId int
						err = projResult.Scan(&projectId)
						if err != nil {
							// handle error
						}
						projectArray = append(projectArray, projectId)
					}
					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":          funderId,
							"name":        funderName,
							"target":      0,
							"actual":      0,
							"day2":        0,
							"women":       0,
							"enrolled":    0,
							"villages":    0,
							"startDate":   "",
							"endDate":     "",
							"select_type": "2",
						}
						data = append(data, obj)
						continue
					}
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}
					obj := map[string]interface{}{
						"id":       funderId,
						"name":     funderName,
						"target":   getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":   getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"day2":     0,
						"women":    getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"enrolled": getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", ""),
						"villages": newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, ""), // New village count function anas

						"startDate":       "",
						"endDate":         "",
						"select_type":     "2",
						"greenMotivators": greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
						"vyapar":          Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
					}
					fmt.Println(funderRow)
					if day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""); day1Count > 0 {
						day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
						obj["day2"] = int(day2Turnout * 100)
					}
					data = append(data, obj)
				}
			}

			projectList := ""

			// summaryFilter := ""

			if reqBody.ProjectID > 0 {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + strconv.Itoa(reqBody.ProjectID) + filter + dateFilterNew
				// summaryFilter := " and p.id = " + strdbv.Itoa(projectId)
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else if reqBody.TrainerID > 0 {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + strconv.Itoa(reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + strconv.Itoa(reqBody.TrainerID)
			} else if reqBody.OpsManager > 0 {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + strconv.Itoa(reqBody.OpsManager)
			} else if reqBody.SOMID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ")"
			} else if reqBody.GFLID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ")"
			} else if (isDateFilterApplied == true && reqBody.PartnerID == 0 && reqBody.Dist == 0 && reqBody.FunderId == 0) || (reqBody.RoleID == 4 && reqBody.Dist == 0) {
				//role 4 - OpsManager Default should be project list without location filter
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
			}
			fmt.Println(summaryFilter)

			if len(projectList) > 0 {
				res, err := db.Query(projectList)
				if err != nil {
					// handle error
				}
				defer res.Close()

				for res.Next() {
					var obj = make(map[string]interface{})
					var projectArray []int
					var id int
					var name string
					var startDate string
					var endDate string

					err := res.Scan(&id, &name, &startDate, &endDate)
					fmt.Println(err)
					if err != nil {
						fmt.Println(err)
					}

					obj["id"] = id
					obj["name"] = name

					projectArray = append(projectArray, id)

					var tpFilter string
					var tbFilter string

					if reqBody.TrainerID > 0 {
						target := getTrainerTarget(db, reqBody.TrainerID, projectArray)
						obj["target"] = target
						summaryTarget += target
						tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", reqBody.TrainerID)
						tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.TrainerID)
					} else {
						target := getTarget(db, startDate, endDate, projectArray)
						obj["target"] = target
						summaryTarget += target
					}

					actual := getActual(db, startDate, endDate, projectArray, tpFilter)
					obj["actual"] = actual
					summaryActuals += actual

					day1Count := getDay1Count(db, startDate, endDate, projectArray, tpFilter)
					summaryDay1 += day1Count

					if day1Count > 0 {
						day2Turnout := float64(actual) / float64(day1Count)
						obj["day2"] = int(math.Round(day2Turnout * 100))
					} else {
						obj["day2"] = 0
					}

					obj["women"] = actual
					obj["enrolled"] = getGelathi(db, startDate, endDate, projectArray, tpFilter, "", "")
					summaryEnrolled += obj["enrolled"].(int)
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}

					obj["villages"] = newVillageCount(db, startDate, endDate, strSlice, tbFilter)
					summaryVillages += obj["villages"].(int)

					obj["startDate"] = startDate
					obj["endDate"] = endDate
					obj["select_type"] = "1"

					obj["greenMotivators"] = greenMotivators(db, startDate, endDate, projectArray, tpFilter, "")
					obj["vyapar"] = Vyapar(db, startDate, endDate, projectArray, tpFilter, "")
					summaryGreen += obj["greenMotivators"].(int)
					summaryVyapar += obj["vyapar"].(int)

					data = append(data, obj)
					fmt.Println(data...)
					// json.NewEncoder(w).Encode(map[string]interface{}{"funder": data})
					json.NewEncoder(w).Encode(map[string]interface{}{"No of Green Cohorts": NoofGreenCohorts(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "No Of Green Enrolled Vypar": greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""), "No Of Green survey": GetNoOfgreenSurvey(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Vyapar module completed": GetNoofGreenModuleCompleted(db), "funder": data})

				}

				// data = append(data, obj)
			}

		} else if reqBody.RoleID == 5 {
			var dateFilter string
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if reqBody.ProjectID > 0 {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) + " and tp.project_id = " + strconv.Itoa(reqBody.ProjectID) +
					dateFilter +
					" GROUP by tp.project_id"
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := db.Query(query)

			if err != nil {
				log.Fatal(err)
			}
			var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Fatal(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				target := getTrainerTarget(db, reqBody.EmpID, projectArray)
				obj["target"] = target
				summaryTarget += target

				filter := " and tp.trainer_id = " + strconv.Itoa(reqBody.EmpID)
				actual := getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = actual
				summaryActuals += actual

				day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2TurnOut := float64(actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2TurnOut * 100))
				} else {
					obj["day2"] = 0
				}

				obj["women"] = obj["actual"]
				obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", "")
				var tbFilter string

				summaryEnrolled += obj["enrolled"].(int)
				tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = obj["startDate"]
				obj["endDate"] = obj["endDate"]
				obj["select_type"] = "1"
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			response := make(map[string]interface{})

			response["summary_target"] = summaryTarget
			response["summary_women"] = summaryActuals
			tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {
				// check if the element is of type int
				if i, ok := v.(int); ok {
					// append the int value to the []int slice
					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			var day2Turnout float64

			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				day2Turnout = 0
				response["summary_day2"] = 0
			}

			response["summary_enrolled"] = summaryEnrolled
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			// jsonResponse, err := json.Marshal(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"summary": response})
			if err != nil {
				log.Fatal(err)
			}

			// return map[string]interface{}{"data": data, "summary": summary}

		} else if reqBody.RoleID == 6 {
			participantFilter := ""
			var filter, filterG, filterV string
			var isDateFilterApplied bool

			if reqBody.ProjectID > 0 {
				filter = fmt.Sprintf(" and tp.project_id = %d", reqBody.ProjectID)
			} else {
				if isDateFilterApplied {
					filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate >= '%s' and p.GreenMotivatorsDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate >= '%s' and p.VyaparEnrollmentDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
				} else {
					filter = " and p.endDate >= CURRENT_DATE()"
					filterG = " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
					filterV = " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
				}
			}

			circleMeet := getGFData(db, filter, 1, reqBody.EmpID)
			villageVisit := getGFData(db, filter, 2, reqBody.EmpID)
			beehive := getGFData(db, filter, 3, reqBody.EmpID)
			enrolled, _ := getGfEnrolled(db, filter, reqBody.EmpID)
			circleVisit := getGFCircle(db, filter, reqBody.EmpID)

			data := []map[string]interface{}{}
			getProjs := fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.user_id = %d %s GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d %s", reqBody.EmpID, filter, reqBody.EmpID, filter)

			if reqBody.ProjectID > 0 {
				getProjs = fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.project_id = %d and tp.user_id = %d GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d and tp.project_id = %d", reqBody.ProjectID, reqBody.EmpID, reqBody.EmpID, reqBody.ProjectID)
			}
			projectsList, err := db.Query(getProjs)
			if err != nil {
				fmt.Println(err)
			}
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					panic(err)
				}

				obj := make(map[string]interface{})
				obj["name"] = name
				// var villageProjvisit *int

				prjFilter := fmt.Sprintf(" and p.id = %d", id)
				circleProjMeet := getGFData(db, prjFilter, 1, reqBody.EmpID)
				obj["circle_meet"] = circleProjMeet
				obj["circles"] = getGFCircle(db, prjFilter, reqBody.EmpID)
				villageProjvisit := getGFData(db, prjFilter, 2, reqBody.EmpID)
				obj["villagevisit"] = villageProjvisit
				beehiveProj := getGFData(db, prjFilter, 3, reqBody.EmpID)
				obj["beehive"] = beehiveProj
				projEnrolled, _ := getGfEnrolled(db, prjFilter, reqBody.EmpID)
				bool1, _ := strconv.ParseBool(participantFilter)

				if bool1 {
					projEnrolled = getParticipantFilterGfEnrolled(db, prjFilter, reqBody.EmpID, reqBody.StartDate, reqBody.EndDate)
				}

				obj = make(map[string]interface{})
				obj["enroll"] = projEnrolled

				var projectArray []string
				str := strconv.Itoa(reqBody.EmpID)
				empIDs := strings.Split(str, ",")
				empIDsStr := "'" + strings.Join(empIDs, "','") + "'"
				projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", empIDsStr))

				// projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", reqBody.EmpID))
				if err != nil {
					// handle error
				}
				defer projectResult.Close()
				for projectResult.Next() {
					var prjIDs string
					err := projectResult.Scan(&prjIDs)
					if err != nil {
						// handle error
					}
					projectArray = append(projectArray, prjIDs)
				}
				intSlice := make([]int, len(projectArray))

				for i, str := range projectArray {
					num, err := strconv.Atoi(str)
					if err != nil {
						panic(err)
					}
					intSlice[i] = num
				}

				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterG)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterV)
				summaryGreen := 0
				summaryVyapar := 0

				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)

				data = append(data, obj)
				fmt.Println(data)

				response := make(map[string]interface{})
				response["summary_circle_meet"] = circleMeet
				response["summary_circles"] = circleVisit
				response["summary_villagevisit"] = villageVisit
				response["summary_beehive"] = beehive
				response["summary_enroll"] = enrolled
				response["summary_green"] = summaryGreen
				response["summary_vyapar"] = summaryVyapar
				response["data"] = data
				response["code"] = 200
				response["success"] = true
				response["message"] = "Successfully"

				js, err := json.Marshal(response)
				fmt.Println(response)
				if err != nil {
					// handle error
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}

		} else if reqBody.RoleID == 13 {
			data := []map[string]interface{}{}
			filter := ""
			filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				filter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterG = " and tp.GreenMotivatorsDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterV = " and tp.VyaparEnrollmentDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			}
			// Additional code for roleId 13 can be added here
			filters := ""
			str := strconv.Itoa(reqBody.ProjectID)
			if str != "" {
				filters = " and p.id = " + str
			}
			// var request string

			// f := ""
			// if gfId, ok := request["gfId"]; ok && gfId != "" {
			// 	f = " and id=" + gfId
			// }

			summarycircleMeet := 0
			summarycircles := 0
			summaryvillagevisit := 0
			summarybeehive := 0
			summaryenroll := 0
			summaryGreen := 0
			summaryVyapar := 0
			summarycircle_meet := 0

			em, err := db.Query("SELECT id from employee e WHERE status =1 AND  e.supervisorId = ?", reqBody.EmpID)
			if err != nil {
				// handle error
			}
			ids := []int{}
			for em.Next() {
				var id int
				err := em.Scan(&id)
				if err != nil {
					// handle error
				}
				ids = append(ids, id)
			}

			getProjs := "Select project_id as id,p.projectName as name from tbl_poa tp " +
				"inner join project p on p.id = tp.project_id " +
				"where  p.gfl_id = ?" + filters + " GROUP by tp.project_id"

			projectsList, err := db.Query(getProjs, reqBody.EmpID)
			if err != nil {
				// handle error
			}
			participantFilter := ""
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					// handle error
				}

				var prjFilter string
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				circleProjMeet := getGFDataN(db, prjFilter, 1, ids)

				obj := make(map[string]interface{})
				obj["name"] = name
				obj["circle_meet"] = circleProjMeet
				summarycircleMeet += circleProjMeet

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and p.endDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				obj["circles"] = getGFCircleN(db, prjFilter, ids)
				summarycircles += obj["circles"].(int)

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var villageProjvisit interface{}
				villageProjvisit = getGFDataN(db, prjFilter, 2, ids)
				obj = make(map[string]interface{})
				obj["villagevisit"] = villageProjvisit
				if villageProjvisit != nil {
					summaryvillagevisit += villageProjvisit.(int)
				} else {
					obj["villagevisit"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var beehiveProj interface{}
				beehiveProj = getGFDataN(db, prjFilter, 3, ids)
				obj["beehive"] = beehiveProj
				if beehiveProj != nil {
					summarybeehive += beehiveProj.(int)
				} else {
					obj["beehive"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var projEnrolled interface{}
				b, err := strconv.ParseBool(participantFilter)
				reqBody.StartDate = "2023-03-15T15:30:45Z"

				// layout of the time string
				layout := ""

				// parse the time string into a time.Time object
				t, err := time.Parse(layout, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}
				reqBody.EndDate = ""

				// layout of the time string
				lay := "2006-01-02T15:04:05Z"

				// parse the time string into a time.Time object
				u, err := time.Parse(lay, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}

				// print the time.Time object
				fmt.Println("Parsed time:", u)
				if b {
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					} else {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					}
				} else {
					projEnrolled = getGfEnrolledN(db, prjFilter, ids)
				}
				obj["enroll"] = projEnrolled
				if projEnrolled != nil {
					summaryenroll += projEnrolled.(int)
				} else {
					obj["enroll"] = "0"
				}
				project_result, err := db.Query("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (" + strconv.Itoa(reqBody.EmpID) + ")")

				if err != nil {
					// handle error
				}
				var projectArray string
				if project_result.Next() {
					var ids interface{}
					err := project_result.Scan(&ids)
					if err != nil {
						// handle error
					}
					projectArray = ids.(string)
				}
				parts := strings.Split(projectArray, "")
				nums := make([]int, len(parts))
				for i, p := range parts {
					num, err := strconv.Atoi(p)
					if err != nil {
						panic(err)
					}
					nums[i] = num
				}
				// fmt.Println(nums)
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, nums, filterG, "")
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, nums, filterV, "")
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			fmt.Println(data)
			response := make(map[string]interface{})
			response["summary_circle_meet"] = summarycircle_meet
			response["summary_circles"] = summarycircles
			response["summary_villagevisit"] = summaryvillagevisit
			response["summary_beehive"] = summarybeehive
			response["summary_enroll"] = summaryenroll
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"
			fmt.Println(response)
			// json.NewEncoder(w).Encode(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"funder": response})

			fmt.Println(filter)

		} else {
			w.WriteHeader(http.StatusCreated)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid role id"
			json.NewEncoder(w).Encode(response)
		}
		fmt.Println(data...)
		fmt.Println(villagesArray...)
		fmt.Println(summaryWomen)
		fmt.Println(summaryDay2)

	})

	mux.HandleFunc("/dashboard/selfsakthi", func(w http.ResponseWriter, r *http.Request) {
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		type ProjectRequest struct {
			PartnerID  int    `json:"partner_id"`
			Dist       int    `json:"dist"`
			Taluk      int    `json:"taluk"`
			Filter     string `json:"filter"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
			FunderId   int    `json:"funder_id"`
			ProjectID  int    `json:"project_id"`
			TrainerID  int    `json:"trainer_id"`
			OpsManager int    `json:"opsmanager"`
			SOMID      int    `json:"somid"`
			GFLID      int    `json:"gflid"`
			RoleID     int    `json:"roleid"`
			GalathiID  string `json:"galathi_id"`
			EmpID      int    `json:"emp_id"`
		}

		var reqBody ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing request body: %v", err)
			return
		}
		var summaryProjectsArray []interface{}
		data := []interface{}{}
		villagesArray := []interface{}{}
		summaryProjectsArray = []interface{}{}
		summaryTarget := 0
		summaryWomen := 0
		summaryVillages := 0
		summaryActuals := 0
		summaryDay1 := 0
		summaryDay2 := 0

		if reqBody.RoleID == 1 || reqBody.RoleID == 9 || reqBody.RoleID == 3 || reqBody.RoleID == 4 || reqBody.RoleID == 12 {
			filter := ""
			summaryFilter := ""
			var projectArray []int
			day1Count := getDay1Count(db, "", "", projectArray, "")
			// summaryDay1 := 0
			var day2 int
			summaryDay1 += day1Count
			if day1Count > 0 {
				actual := getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, "")
				day2Turnout := float64(actual) / float64(day1Count)
				day2 = int(day2Turnout * 100)
			} else {
				day2 = 0
			}
			if reqBody.RoleID == 3 {
				var opsIds []int
				if reqBody.SOMID != 0 {
					opsIds = getReportingOpsManagers(db, reqBody.SOMID)
				} else if reqBody.GFLID != 0 {
					opsIds = getSupervisor(db, reqBody.GFLID)
				} else {
					opsIds = getReportingOpsManagers(db, reqBody.EmpID)
				}
				filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
			} else if reqBody.RoleID == 12 {
				opsIds := getOpsManagers(db, reqBody.EmpID)
				if len(opsIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					filter = " and p.operations_manager in (0)"
				}
			} else if reqBody.RoleID == 4 {
				// Ops Manager
				projectIds := getOpProjects(db, reqBody.EmpID)
				fmt.Println(projectIds)
				if len(projectIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager = %d", reqBody.EmpID)
				} else {
					showNoProj()
				}
			}

			isDateFilterApplied := false
			var dateFilter string
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				isDateFilterApplied = true
				dateFilter = fmt.Sprintf("startDate >= '%s' AND endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
			} else {
				dateFilter = "endDate >= CURRENT_DATE()"
			}

			var funderListQuery string

			if reqBody.PartnerID > 0 {
				funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) AS id, funderName AS name FROM project p "+
					"INNER JOIN funder ON funder.funderID = p.funderID "+
					"WHERE p.partnerID = %d AND %s %s", reqBody.PartnerID, dateFilter, filter)
				filter += fmt.Sprintf(" AND p.partnerID = %d", reqBody.PartnerID)
			} else if reqBody.Dist > 0 {
				if reqBody.Taluk > 0 {
					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID = %d AND %s %s GROUP BY p.funderID", reqBody.Taluk, dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID = %d", reqBody.Taluk)
				} else {
					// Get taluk of specified dist
					getTaluk := fmt.Sprintf("SELECT id FROM location WHERE `type` = 4 AND parentId = %d", reqBody.Dist)
					talukArray := []int{}
					talukRes, err := db.Query(getTaluk)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error getting taluk list: %v", err)
						return
					}
					defer talukRes.Close()

					for talukRes.Next() {
						var talukID int
						err := talukRes.Scan(&talukID)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							fmt.Fprintf(w, "Error scanning taluk list: %v", err)
							return
						}
						talukArray = append(talukArray, talukID)
					}

					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID IN (%s) AND %s %s GROUP BY p.funderID",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"),
						dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID IN (%s)",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"))
				}
			} else if reqBody.FunderId > 0 {
				funderListQuery = fmt.Sprintf("SELECT funderID as id, funderName as name FROM funder f WHERE funderID = %d", reqBody.FunderId)
				// summaryFilter := fmt.Sprintf(" AND p.funderID = %d", reqBody.FunderId)

			} else if reqBody.PartnerID == 0 && reqBody.TrainerID == 0 && reqBody.OpsManager == 0 && reqBody.SOMID == 0 && reqBody.GFLID == 0 && !isDateFilterApplied && reqBody.RoleID != 4 {
				// Role 4 OpsManager Default should be project list
				funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p INNER JOIN funder ON p.funderId = funder.funderID WHERE " + dateFilter + filter
			}

			funderList := []map[string]interface{}{}

			//gives funder list
			if len(funderListQuery) > 0 {
				res, err := db.Query(funderListQuery)
				if err != nil {
					// handle error
				}
				defer res.Close()
				for res.Next() {
					data := []interface{}{}

					projectArray := []int{}
					funderRow := map[string]interface{}{}
					var funderId int
					var funderName string
					err = res.Scan(&funderId, &funderName)
					if err != nil {
						// handle error
					}
					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter //commented by anas for ceo funder fiter not working
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id,startDate,endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}
					projResult, err := db.Query(getProj)
					if err != nil {
						// handle error
					}
					defer projResult.Close()
					for projResult.Next() {
						var projectId int
						err = projResult.Scan(&projectId)
						if err != nil {
							// handle error
						}
						projectArray = append(projectArray, projectId)
					}
					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":          funderId,
							"name":        funderName,
							"target":      0,
							"actual":      0,
							"day2":        0,
							"women":       0,
							"enrolled":    0,
							"villages":    0,
							"startDate":   "",
							"endDate":     "",
							"select_type": "2",
						}
						data = append(data, obj)
						continue
					}
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}
					obj := map[string]interface{}{
						"id":       funderId,
						"name":     funderName,
						"target":   getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":   getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"day2":     0,
						"women":    getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"enrolled": getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", ""),
						"villages": newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, ""), // New village count function anas

						"startDate":       "",
						"endDate":         "",
						"select_type":     "2",
						"greenMotivators": greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
						"vyapar":          Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
					}
					fmt.Println(funderRow)
					if day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""); day1Count > 0 {
						day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
						obj["day2"] = int(day2Turnout * 100)
					}
					data = append(data, obj)
				}
			}

			projectList := ""
			var summaryEnrolled, summaryGreen, summaryVyapar int

			// summaryFilter := ""

			if reqBody.ProjectID > 0 {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + strconv.Itoa(reqBody.ProjectID) + filter + dateFilterNew
				// summaryFilter := " and p.id = " + strdbv.Itoa(projectId)
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else if reqBody.TrainerID > 0 {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + strconv.Itoa(reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + strconv.Itoa(reqBody.TrainerID)
			} else if reqBody.OpsManager > 0 {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + strconv.Itoa(reqBody.OpsManager)
			} else if reqBody.SOMID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ")"
			} else if reqBody.GFLID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ")"
			} else if (isDateFilterApplied == true && reqBody.PartnerID == 0 && reqBody.Dist == 0 && reqBody.FunderId == 0) || (reqBody.RoleID == 4 && reqBody.Dist == 0) {
				//role 4 - OpsManager Default should be project list without location filter
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
			}
			fmt.Println(summaryFilter)

			if len(projectList) > 0 {
				res, err := db.Query(projectList)
				if err != nil {
					// handle error
				}
				defer res.Close()

				for res.Next() {
					var obj = make(map[string]interface{})
					var projectArray []int
					var id int
					var name string
					var startDate string
					var endDate string

					err := res.Scan(&id, &name, &startDate, &endDate)
					fmt.Println(err)
					if err != nil {
						fmt.Println(err)
					}

					obj["id"] = id
					obj["name"] = name

					projectArray = append(projectArray, id)

					var tpFilter string
					var tbFilter string

					if reqBody.TrainerID > 0 {
						target := getTrainerTarget(db, reqBody.TrainerID, projectArray)
						obj["target"] = target
						summaryTarget += target
						tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", reqBody.TrainerID)
						tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.TrainerID)
					} else {
						target := getTarget(db, startDate, endDate, projectArray)
						obj["target"] = target
						summaryTarget += target
					}

					actual := getActual(db, startDate, endDate, projectArray, tpFilter)
					obj["actual"] = actual
					summaryActuals += actual

					day1Count := getDay1Count(db, startDate, endDate, projectArray, tpFilter)
					summaryDay1 += day1Count

					if day1Count > 0 {
						day2Turnout := float64(actual) / float64(day1Count)
						obj["day2"] = int(math.Round(day2Turnout * 100))
					} else {
						obj["day2"] = 0
					}

					obj["women"] = actual
					obj["enrolled"] = getGelathi(db, startDate, endDate, projectArray, tpFilter, "", "")
					summaryEnrolled += obj["enrolled"].(int)
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}

					obj["villages"] = newVillageCount(db, startDate, endDate, strSlice, tbFilter)
					summaryVillages += obj["villages"].(int)

					obj["startDate"] = startDate
					obj["endDate"] = endDate
					obj["select_type"] = "1"

					obj["greenMotivators"] = greenMotivators(db, startDate, endDate, projectArray, tpFilter, "")
					obj["vyapar"] = Vyapar(db, startDate, endDate, projectArray, tpFilter, "")
					summaryGreen += obj["greenMotivators"].(int)
					summaryVyapar += obj["vyapar"].(int)

					data = append(data, obj)
					fmt.Println(data...)

				}
				json.NewEncoder(w).Encode(map[string]interface{}{"Target": getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray), "No Of Villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "Actual": getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "No Of Batches": getTrainingBatches(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "No Of self shakthi survey": GetNoofGreenModuleCompleted(db), "2nd Day turnout %": day2, "funder": data})

			}

			fmt.Println(summaryFilter)
			fmt.Println(funderList)

		} else if reqBody.RoleID == 5 {
			var dateFilter string
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if reqBody.ProjectID > 0 {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) + " and tp.project_id = " + strconv.Itoa(reqBody.ProjectID) +
					dateFilter +
					" GROUP by tp.project_id"
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := db.Query(query)

			if err != nil {
				log.Fatal(err)
			}
			var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Fatal(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				target := getTrainerTarget(db, reqBody.EmpID, projectArray)
				obj["target"] = target
				summaryTarget += target

				filter := " and tp.trainer_id = " + strconv.Itoa(reqBody.EmpID)
				actual := getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = actual
				summaryActuals += actual

				day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2TurnOut := float64(actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2TurnOut * 100))
				} else {
					obj["day2"] = 0
				}

				obj["women"] = obj["actual"]
				obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", "")
				var tbFilter string

				summaryEnrolled += obj["enrolled"].(int)
				tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = obj["startDate"]
				obj["endDate"] = obj["endDate"]
				obj["select_type"] = "1"
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			response := make(map[string]interface{})

			response["summary_target"] = summaryTarget
			response["summary_women"] = summaryActuals
			tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {
				// check if the element is of type int
				if i, ok := v.(int); ok {
					// append the int value to the []int slice
					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			var day2Turnout float64

			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				day2Turnout = 0
				response["summary_day2"] = 0
			}

			response["summary_enrolled"] = summaryEnrolled
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			// jsonResponse, err := json.Marshal(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"summary": response})
			if err != nil {
				log.Fatal(err)
			}

			// return map[string]interface{}{"data": data, "summary": summary}

		} else if reqBody.RoleID == 6 {
			participantFilter := ""
			var filter, filterG, filterV string
			var isDateFilterApplied bool

			if reqBody.ProjectID > 0 {
				filter = fmt.Sprintf(" and tp.project_id = %d", reqBody.ProjectID)
			} else {
				if isDateFilterApplied {
					filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate >= '%s' and p.GreenMotivatorsDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate >= '%s' and p.VyaparEnrollmentDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
				} else {
					filter = " and p.endDate >= CURRENT_DATE()"
					filterG = " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
					filterV = " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
				}
			}

			circleMeet := getGFData(db, filter, 1, reqBody.EmpID)
			villageVisit := getGFData(db, filter, 2, reqBody.EmpID)
			beehive := getGFData(db, filter, 3, reqBody.EmpID)
			enrolled, _ := getGfEnrolled(db, filter, reqBody.EmpID)
			circleVisit := getGFCircle(db, filter, reqBody.EmpID)

			data := []map[string]interface{}{}
			getProjs := fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.user_id = %d %s GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d %s", reqBody.EmpID, filter, reqBody.EmpID, filter)

			if reqBody.ProjectID > 0 {
				getProjs = fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.project_id = %d and tp.user_id = %d GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d and tp.project_id = %d", reqBody.ProjectID, reqBody.EmpID, reqBody.EmpID, reqBody.ProjectID)
			}
			projectsList, err := db.Query(getProjs)
			if err != nil {
				fmt.Println(err)
			}
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					panic(err)
				}

				obj := make(map[string]interface{})
				obj["name"] = name
				// var villageProjvisit *int

				prjFilter := fmt.Sprintf(" and p.id = %d", id)
				circleProjMeet := getGFData(db, prjFilter, 1, reqBody.EmpID)
				obj["circle_meet"] = circleProjMeet
				obj["circles"] = getGFCircle(db, prjFilter, reqBody.EmpID)
				villageProjvisit := getGFData(db, prjFilter, 2, reqBody.EmpID)
				obj["villagevisit"] = villageProjvisit
				beehiveProj := getGFData(db, prjFilter, 3, reqBody.EmpID)
				obj["beehive"] = beehiveProj
				projEnrolled, _ := getGfEnrolled(db, prjFilter, reqBody.EmpID)
				bool1, _ := strconv.ParseBool(participantFilter)

				if bool1 {
					projEnrolled = getParticipantFilterGfEnrolled(db, prjFilter, reqBody.EmpID, reqBody.StartDate, reqBody.EndDate)
				}

				obj = make(map[string]interface{})
				obj["enroll"] = projEnrolled

				var projectArray []string
				str := strconv.Itoa(reqBody.EmpID)
				empIDs := strings.Split(str, ",")
				empIDsStr := "'" + strings.Join(empIDs, "','") + "'"
				projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", empIDsStr))

				// projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", reqBody.EmpID))
				if err != nil {
					// handle error
				}
				defer projectResult.Close()
				for projectResult.Next() {
					var prjIDs string
					err := projectResult.Scan(&prjIDs)
					if err != nil {
						// handle error
					}
					projectArray = append(projectArray, prjIDs)
				}
				intSlice := make([]int, len(projectArray))

				for i, str := range projectArray {
					num, err := strconv.Atoi(str)
					if err != nil {
						panic(err)
					}
					intSlice[i] = num
				}

				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterG)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterV)
				summaryGreen := 0
				summaryVyapar := 0

				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)

				data = append(data, obj)
				fmt.Println(data)

				response := make(map[string]interface{})
				response["summary_circle_meet"] = circleMeet
				response["summary_circles"] = circleVisit
				response["summary_villagevisit"] = villageVisit
				response["summary_beehive"] = beehive
				response["summary_enroll"] = enrolled
				response["summary_green"] = summaryGreen
				response["summary_vyapar"] = summaryVyapar
				response["data"] = data
				response["code"] = 200
				response["success"] = true
				response["message"] = "Successfully"

				js, err := json.Marshal(response)
				fmt.Println(response)
				if err != nil {
					// handle error
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}

		} else if reqBody.RoleID == 13 {
			data := []map[string]interface{}{}
			filter := ""
			filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				filter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterG = " and tp.GreenMotivatorsDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterV = " and tp.VyaparEnrollmentDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			}
			// Additional code for roleId 13 can be added here
			filters := ""
			str := strconv.Itoa(reqBody.ProjectID)
			if str != "" {
				filters = " and p.id = " + str
			}
			// var request string

			// f := ""
			// if gfId, ok := request["gfId"]; ok && gfId != "" {
			// 	f = " and id=" + gfId
			// }

			summarycircleMeet := 0
			summarycircles := 0
			summaryvillagevisit := 0
			summarybeehive := 0
			summaryenroll := 0
			summaryGreen := 0
			summaryVyapar := 0
			summarycircle_meet := 0

			em, err := db.Query("SELECT id from employee e WHERE status =1 AND  e.supervisorId = ?", reqBody.EmpID)
			if err != nil {
				// handle error
			}
			ids := []int{}
			for em.Next() {
				var id int
				err := em.Scan(&id)
				if err != nil {
					// handle error
				}
				ids = append(ids, id)
			}

			getProjs := "Select project_id as id,p.projectName as name from tbl_poa tp " +
				"inner join project p on p.id = tp.project_id " +
				"where  p.gfl_id = ?" + filters + " GROUP by tp.project_id"

			projectsList, err := db.Query(getProjs, reqBody.EmpID)
			if err != nil {
				// handle error
			}
			participantFilter := ""
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					// handle error
				}

				var prjFilter string
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				circleProjMeet := getGFDataN(db, prjFilter, 1, ids)

				obj := make(map[string]interface{})
				obj["name"] = name
				obj["circle_meet"] = circleProjMeet
				summarycircleMeet += circleProjMeet

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and p.endDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				obj["circles"] = getGFCircleN(db, prjFilter, ids)
				summarycircles += obj["circles"].(int)

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var villageProjvisit interface{}
				villageProjvisit = getGFDataN(db, prjFilter, 2, ids)
				obj = make(map[string]interface{})
				obj["villagevisit"] = villageProjvisit
				if villageProjvisit != nil {
					summaryvillagevisit += villageProjvisit.(int)
				} else {
					obj["villagevisit"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var beehiveProj interface{}
				beehiveProj = getGFDataN(db, prjFilter, 3, ids)
				obj["beehive"] = beehiveProj
				if beehiveProj != nil {
					summarybeehive += beehiveProj.(int)
				} else {
					obj["beehive"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var projEnrolled interface{}
				b, err := strconv.ParseBool(participantFilter)
				reqBody.StartDate = "2023-03-15T15:30:45Z"

				// layout of the time string
				layout := ""

				// parse the time string into a time.Time object
				t, err := time.Parse(layout, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}
				reqBody.EndDate = ""

				// layout of the time string
				lay := "2006-01-02T15:04:05Z"

				// parse the time string into a time.Time object
				u, err := time.Parse(lay, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}

				// print the time.Time object
				fmt.Println("Parsed time:", u)
				if b {
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					} else {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					}
				} else {
					projEnrolled = getGfEnrolledN(db, prjFilter, ids)
				}
				obj["enroll"] = projEnrolled
				if projEnrolled != nil {
					summaryenroll += projEnrolled.(int)
				} else {
					obj["enroll"] = "0"
				}
				project_result, err := db.Query("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (" + strconv.Itoa(reqBody.EmpID) + ")")

				if err != nil {
					// handle error
				}
				var projectArray string
				if project_result.Next() {
					var ids interface{}
					err := project_result.Scan(&ids)
					if err != nil {
						// handle error
					}
					projectArray = ids.(string)
				}
				parts := strings.Split(projectArray, "")
				nums := make([]int, len(parts))
				for i, p := range parts {
					num, err := strconv.Atoi(p)
					if err != nil {
						panic(err)
					}
					nums[i] = num
				}
				// fmt.Println(nums)
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, nums, filterG, "")
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, nums, filterV, "")
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			fmt.Println(data)
			response := make(map[string]interface{})
			response["summary_circle_meet"] = summarycircle_meet
			response["summary_circles"] = summarycircles
			response["summary_villagevisit"] = summaryvillagevisit
			response["summary_beehive"] = summarybeehive
			response["summary_enroll"] = summaryenroll
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"
			fmt.Println(response)
			// json.NewEncoder(w).Encode(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"funder": response})

			fmt.Println(filter)

		} else {
			w.WriteHeader(http.StatusCreated)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid role id"
			json.NewEncoder(w).Encode(response)
		}
		fmt.Println(data...)
		fmt.Println(villagesArray...)
		fmt.Println(summaryWomen)
		fmt.Println(summaryDay2)

	})
	mux.HandleFunc("/dashboard/gelathiprogram", func(w http.ResponseWriter, r *http.Request) {
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		type ProjectRequest struct {
			PartnerID  int    `json:"partner_id"`
			Dist       int    `json:"dist"`
			Taluk      int    `json:"taluk"`
			Filter     string `json:"filter"`
			StartDate  string `json:"start_date"`
			EndDate    string `json:"end_date"`
			FunderId   int    `json:"funder_id"`
			ProjectID  int    `json:"project_id"`
			TrainerID  int    `json:"trainer_id"`
			OpsManager int    `json:"opsmanager"`
			SOMID      int    `json:"somid"`
			GFLID      int    `json:"gflid"`
			RoleID     int    `json:"roleid"`
			GalathiID  string `json:"galathi_id"`
			EmpID      int    `json:"emp_id"`
		}

		var reqBody ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing request body: %v", err)
			return
		}
		var summaryProjectsArray []interface{}
		data := []interface{}{}
		villagesArray := []interface{}{}
		summaryProjectsArray = []interface{}{}
		summaryTarget := 0
		summaryWomen := 0
		summaryVillages := 0
		summaryActuals := 0
		summaryDay1 := 0
		summaryDay2 := 0
		// summaryEnrolled := 0
		// summaryGreen := 0
		// summaryVyapar := 0

		if reqBody.RoleID == 1 || reqBody.RoleID == 9 || reqBody.RoleID == 3 || reqBody.RoleID == 4 || reqBody.RoleID == 12 {
			filter := ""
			summaryFilter := ""
			if reqBody.RoleID == 3 {
				var opsIds []int
				if reqBody.SOMID != 0 {
					opsIds = getReportingOpsManagers(db, reqBody.SOMID)
				} else if reqBody.GFLID != 0 {
					opsIds = getSupervisor(db, reqBody.GFLID)
				} else {
					opsIds = getReportingOpsManagers(db, reqBody.EmpID)
				}
				filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
			} else if reqBody.RoleID == 12 {
				opsIds := getOpsManagers(db, reqBody.EmpID)
				if len(opsIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager in (%s)", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(opsIds)), ","), "[]"))
				} else {
					filter = " and p.operations_manager in (0)"
				}
			} else if reqBody.RoleID == 4 {
				// Ops Manager
				projectIds := getOpProjects(db, reqBody.EmpID)
				fmt.Println(projectIds)
				if len(projectIds) > 0 {
					filter = fmt.Sprintf(" and p.operations_manager = %d", reqBody.EmpID)
				} else {
					showNoProj()
				}
			}

			isDateFilterApplied := false
			var dateFilter string
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				isDateFilterApplied = true
				dateFilter = fmt.Sprintf("startDate >= '%s' AND endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
			} else {
				dateFilter = "endDate >= CURRENT_DATE()"
			}

			var funderListQuery string

			if reqBody.PartnerID > 0 {
				funderListQuery = fmt.Sprintf("SELECT DISTINCT(p.funderId) AS id, funderName AS name FROM project p "+
					"INNER JOIN funder ON funder.funderID = p.funderID "+
					"WHERE p.partnerID = %d AND %s %s", reqBody.PartnerID, dateFilter, filter)
				filter += fmt.Sprintf(" AND p.partnerID = %d", reqBody.PartnerID)
			} else if reqBody.Dist > 0 {
				if reqBody.Taluk > 0 {
					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID = %d AND %s %s GROUP BY p.funderID", reqBody.Taluk, dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID = %d", reqBody.Taluk)
				} else {
					// Get taluk of specified dist
					getTaluk := fmt.Sprintf("SELECT id FROM location WHERE `type` = 4 AND parentId = %d", reqBody.Dist)
					talukArray := []int{}
					talukRes, err := db.Query(getTaluk)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Error getting taluk list: %v", err)
						return
					}
					defer talukRes.Close()

					for talukRes.Next() {
						var talukID int
						err := talukRes.Scan(&talukID)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							fmt.Fprintf(w, "Error scanning taluk list: %v", err)
							return
						}
						talukArray = append(talukArray, talukID)
					}

					funderListQuery = fmt.Sprintf("SELECT p.funderID AS id, funderName AS name FROM project p "+
						"INNER JOIN funder ON funder.funderID = p.funderID "+
						"WHERE locationID IN (%s) AND %s %s GROUP BY p.funderID",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"),
						dateFilter, filter)
					filter += fmt.Sprintf(" AND locationID IN (%s)",
						strings.Trim(strings.Join(strings.Fields(fmt.Sprint(talukArray)), ","), "[]"))
				}
			} else if reqBody.FunderId > 0 {
				funderListQuery = fmt.Sprintf("SELECT funderID as id, funderName as name FROM funder f WHERE funderID = %d", reqBody.FunderId)
				// summaryFilter := fmt.Sprintf(" AND p.funderID = %d", reqBody.FunderId)

			} else if reqBody.PartnerID == 0 && reqBody.TrainerID == 0 && reqBody.OpsManager == 0 && reqBody.SOMID == 0 && reqBody.GFLID == 0 && !isDateFilterApplied && reqBody.RoleID != 4 {
				// Role 4 OpsManager Default should be project list
				funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p INNER JOIN funder ON p.funderId = funder.funderID WHERE " + dateFilter + filter
			}

			funderList := []map[string]interface{}{}

			//gives funder list
			if len(funderListQuery) > 0 {
				res, err := db.Query(funderListQuery)
				if err != nil {
					// handle error
				}
				defer res.Close()
				for res.Next() {
					data := []interface{}{}

					projectArray := []int{}
					funderRow := map[string]interface{}{}
					var funderId int
					var funderName string
					err = res.Scan(&funderId, &funderName)
					if err != nil {
						// handle error
					}
					getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter //commented by anas for ceo funder fiter not working
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						getProj = "SELECT id,startDate,endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
					}
					projResult, err := db.Query(getProj)
					if err != nil {
						// handle error
					}
					defer projResult.Close()
					for projResult.Next() {
						var projectId int
						err = projResult.Scan(&projectId)
						if err != nil {
							// handle error
						}
						projectArray = append(projectArray, projectId)
					}
					if len(projectArray) == 0 {
						obj := map[string]interface{}{
							"id":          funderId,
							"name":        funderName,
							"target":      0,
							"actual":      0,
							"day2":        0,
							"women":       0,
							"enrolled":    0,
							"villages":    0,
							"startDate":   "",
							"endDate":     "",
							"select_type": "2",
						}
						data = append(data, obj)
						continue
					}
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}
					obj := map[string]interface{}{
						"id":       funderId,
						"name":     funderName,
						"target":   getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray),
						"actual":   getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"day2":     0,
						"women":    getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""),
						"enrolled": getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", ""),
						"villages": newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, ""), // New village count function anas

						"startDate":       "",
						"endDate":         "",
						"select_type":     "2",
						"greenMotivators": greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
						"vyapar":          Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""),
					}
					fmt.Println(funderRow)
					if day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""); day1Count > 0 {
						day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
						obj["day2"] = int(day2Turnout * 100)
					}
					data = append(data, obj)
				}
			}

			projectList := ""
			var summaryEnrolled, summaryGreen, summaryVyapar int

			// summaryFilter := ""

			if reqBody.ProjectID > 0 {
				dateFilterNew := ""
				if isDateFilterApplied {
					dateFilterNew = " and startDate >= '" + reqBody.StartDate + "' and endDate <= '" + reqBody.EndDate + "'"
				}
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where id = " + strconv.Itoa(reqBody.ProjectID) + filter + dateFilterNew
				// summaryFilter := " and p.id = " + strdbv.Itoa(projectId)
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else if reqBody.TrainerID > 0 {
				projectList = "SELECT project_id as id,projectName as name,p.startDate,p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + strconv.Itoa(reqBody.TrainerID) + " and " + dateFilter + filter + " GROUP  by project_id"
				summaryFilter = " and tp.user_id = " + strconv.Itoa(reqBody.TrainerID)
			} else if reqBody.OpsManager > 0 {
				if dateFilter == "" || (reqBody.StartDate == "" && reqBody.EndDate == "") {
					projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and " + dateFilter + filter + " GROUP by id "
				} else {
					projectList = "SELECT p.id,p.projectName as name,p.startDate,p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + strconv.Itoa(reqBody.OpsManager) + " and tp.participant_day2 >= '" + reqBody.StartDate + "' and tp.participant_day2 <= '" + reqBody.EndDate + "' GROUP by p.id "
				}
				summaryFilter = " and p.operations_manager = " + strconv.Itoa(reqBody.OpsManager)
			} else if reqBody.SOMID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + strconv.Itoa(reqBody.SOMID) + ")"
			} else if reqBody.GFLID > 0 {
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ") and " + dateFilter + filter + " GROUP by id "
				summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + strconv.Itoa(reqBody.GFLID) + ")"
			} else if (isDateFilterApplied == true && reqBody.PartnerID == 0 && reqBody.Dist == 0 && reqBody.FunderId == 0) || (reqBody.RoleID == 4 && reqBody.Dist == 0) {
				//role 4 - OpsManager Default should be project list without location filter
				projectList = "SELECT id,projectName as name,p.startDate,p.endDate from project p where " + dateFilter + filter
			}
			fmt.Println(summaryFilter)

			if len(projectList) > 0 {
				res, err := db.Query(projectList)
				if err != nil {
					// handle error
				}
				defer res.Close()

				for res.Next() {
					var obj = make(map[string]interface{})
					var projectArray []int
					var id int
					var name string
					var startDate string
					var endDate string

					err := res.Scan(&id, &name, &startDate, &endDate)
					fmt.Println(err)
					if err != nil {
						fmt.Println(err)
					}

					obj["id"] = id
					obj["name"] = name

					projectArray = append(projectArray, id)

					var tpFilter string
					var tbFilter string

					if reqBody.TrainerID > 0 {
						target := getTrainerTarget(db, reqBody.TrainerID, projectArray)
						obj["target"] = target
						summaryTarget += target
						tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", reqBody.TrainerID)
						tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.TrainerID)
					} else {
						target := getTarget(db, startDate, endDate, projectArray)
						obj["target"] = target
						summaryTarget += target
					}

					actual := getActual(db, startDate, endDate, projectArray, tpFilter)
					obj["actual"] = actual
					summaryActuals += actual

					day1Count := getDay1Count(db, startDate, endDate, projectArray, tpFilter)
					summaryDay1 += day1Count

					if day1Count > 0 {
						day2Turnout := float64(actual) / float64(day1Count)
						obj["day2"] = int(math.Round(day2Turnout * 100))
					} else {
						obj["day2"] = 0
					}

					obj["women"] = actual
					obj["enrolled"] = getGelathi(db, startDate, endDate, projectArray, tpFilter, "", "")
					summaryEnrolled += obj["enrolled"].(int)
					var strSlice []string
					for _, num := range projectArray {
						strSlice = append(strSlice, strconv.Itoa(num))
					}

					obj["villages"] = newVillageCount(db, startDate, endDate, strSlice, tbFilter)
					summaryVillages += obj["villages"].(int)

					obj["startDate"] = startDate
					obj["endDate"] = endDate
					obj["select_type"] = "1"

					obj["greenMotivators"] = greenMotivators(db, startDate, endDate, projectArray, tpFilter, "")
					obj["vyapar"] = Vyapar(db, startDate, endDate, projectArray, tpFilter, "")
					summaryGreen += obj["greenMotivators"].(int)
					summaryVyapar += obj["vyapar"].(int)

					data = append(data, obj)
					fmt.Println(data...)
					json.NewEncoder(w).Encode(map[string]interface{}{"No Of Villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "No Of gelathi Enrolled Vypar": getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", ""), "No Of Sporthi survey": GetNoOfSporthiSurvey(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Sporthi module completed": GetNoofSporthiModuleCompleted(db), "No of beehives": getGFData(db, filter, 0, 0), "data": data})

				}

			}

			fmt.Println(summaryFilter)
			fmt.Println(funderList)

		} else if reqBody.RoleID == 5 {
			var dateFilter string
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			var query string
			if reqBody.ProjectID > 0 {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) + " and tp.project_id = " + strconv.Itoa(reqBody.ProjectID) +
					dateFilter +
					" GROUP by tp.project_id"
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			} else {
				query = "SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate " +
					"from tbl_poa tp " +
					"inner join project p on p.id = tp.project_id " +
					"where user_id = " + strconv.Itoa(reqBody.EmpID) +
					dateFilter +
					" GROUP by project_id"
			}

			res, err := db.Query(query)

			if err != nil {
				log.Fatal(err)
			}
			var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			for res.Next() {
				var obj = make(map[string]interface{})
				var projectArray []int
				var id int
				var name string
				var startDate, endDate string

				err := res.Scan(&id, &name, &startDate, &endDate)

				if err != nil {
					log.Fatal(err)
				}

				projectArray = append(projectArray, id)
				obj = make(map[string]interface{})

				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				target := getTrainerTarget(db, reqBody.EmpID, projectArray)
				obj["target"] = target
				summaryTarget += target

				filter := " and tp.trainer_id = " + strconv.Itoa(reqBody.EmpID)
				actual := getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				obj["actual"] = actual
				summaryActuals += actual

				day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2TurnOut := float64(actual) / float64(day1Count)
					obj["day2"] = int(math.Round(day2TurnOut * 100))
				} else {
					obj["day2"] = 0
				}

				obj["women"] = obj["actual"]
				obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", "")
				var tbFilter string

				summaryEnrolled += obj["enrolled"].(int)
				tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
				strSlice := make([]string, len(projectArray))

				// loop through each element in intSlice and convert to string
				for i, v := range projectArray {
					strSlice[i] = strconv.Itoa(v)
				}
				obj["villages"] = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, tbFilter)
				summaryVillages += obj["villages"].(int)
				obj["startDate"] = obj["startDate"]
				obj["endDate"] = obj["endDate"]
				obj["select_type"] = "1"
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			response := make(map[string]interface{})

			response["summary_target"] = summaryTarget
			response["summary_women"] = summaryActuals
			tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			intSlice := []int{}

			// loop through each element in the []interface{} slice
			for _, v := range summaryProjectsArray {
				// check if the element is of type int
				if i, ok := v.(int); ok {
					// append the int value to the []int slice
					intSlice = append(intSlice, i)
				}
			}
			response["summary_villages"] = getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			response["summary_actual"] = summaryActuals
			var day2Turnout float64

			if summaryDay1 > 0 {
				day2Turnout = float64(summaryActuals) / float64(summaryDay1)
				response["summary_day2"] = int(math.Round(day2Turnout * 100))
			} else {
				day2Turnout = 0
				response["summary_day2"] = 0
			}

			response["summary_enrolled"] = summaryEnrolled
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"

			// jsonResponse, err := json.Marshal(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"summary": response})
			if err != nil {
				log.Fatal(err)
			}

			// return map[string]interface{}{"data": data, "summary": summary}

		} else if reqBody.RoleID == 6 {
			participantFilter := ""
			var filter, filterG, filterV string
			var isDateFilterApplied bool

			if reqBody.ProjectID > 0 {
				filter = fmt.Sprintf(" and tp.project_id = %d", reqBody.ProjectID)
			} else {
				if isDateFilterApplied {
					filter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterG = fmt.Sprintf(" and tp.GreenMotivatorsDate >= '%s' and p.GreenMotivatorsDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
					filterV = fmt.Sprintf(" and tp.VyaparEnrollmentDate >= '%s' and p.VyaparEnrollmentDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
				} else {
					filter = " and p.endDate >= CURRENT_DATE()"
					filterG = " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
					filterV = " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
				}
			}

			circleMeet := getGFData(db, filter, 1, reqBody.EmpID)
			villageVisit := getGFData(db, filter, 2, reqBody.EmpID)
			beehive := getGFData(db, filter, 3, reqBody.EmpID)
			enrolled, _ := getGfEnrolled(db, filter, reqBody.EmpID)
			circleVisit := getGFCircle(db, filter, reqBody.EmpID)

			data := []map[string]interface{}{}
			getProjs := fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.user_id = %d %s GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d %s", reqBody.EmpID, filter, reqBody.EmpID, filter)

			if reqBody.ProjectID > 0 {
				getProjs = fmt.Sprintf("Select project_id as id,p.projectName as name from tbl_poa tp inner join project p on p.id = tp.project_id where tp.project_id = %d and tp.user_id = %d GROUP by tp.project_id UNION SELECT tp.project_id as id,p.projectName as name from training_participants tp inner join project p on tp.project_id = p.id where enroll = 1 and gelathi_id = %d and tp.project_id = %d", reqBody.ProjectID, reqBody.EmpID, reqBody.EmpID, reqBody.ProjectID)
			}
			projectsList, err := db.Query(getProjs)
			if err != nil {
				fmt.Println(err)
			}
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					panic(err)
				}

				obj := make(map[string]interface{})
				obj["name"] = name
				// var villageProjvisit *int

				prjFilter := fmt.Sprintf(" and p.id = %d", id)
				circleProjMeet := getGFData(db, prjFilter, 1, reqBody.EmpID)
				obj["circle_meet"] = circleProjMeet
				obj["circles"] = getGFCircle(db, prjFilter, reqBody.EmpID)
				villageProjvisit := getGFData(db, prjFilter, 2, reqBody.EmpID)
				obj["villagevisit"] = villageProjvisit
				beehiveProj := getGFData(db, prjFilter, 3, reqBody.EmpID)
				obj["beehive"] = beehiveProj
				projEnrolled, _ := getGfEnrolled(db, prjFilter, reqBody.EmpID)
				bool1, _ := strconv.ParseBool(participantFilter)

				if bool1 {
					projEnrolled = getParticipantFilterGfEnrolled(db, prjFilter, reqBody.EmpID, reqBody.StartDate, reqBody.EndDate)
				}

				obj = make(map[string]interface{})
				obj["enroll"] = projEnrolled

				var projectArray []string
				str := strconv.Itoa(reqBody.EmpID)
				empIDs := strings.Split(str, ",")
				empIDsStr := "'" + strings.Join(empIDs, "','") + "'"
				projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", empIDsStr))

				// projectResult, err := db.Query(fmt.Sprintf("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (%s)", reqBody.EmpID))
				if err != nil {
					// handle error
				}
				defer projectResult.Close()
				for projectResult.Next() {
					var prjIDs string
					err := projectResult.Scan(&prjIDs)
					if err != nil {
						// handle error
					}
					projectArray = append(projectArray, prjIDs)
				}
				intSlice := make([]int, len(projectArray))

				for i, str := range projectArray {
					num, err := strconv.Atoi(str)
					if err != nil {
						panic(err)
					}
					intSlice[i] = num
				}

				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterG)
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, intSlice, "", filterV)
				summaryGreen := 0
				summaryVyapar := 0

				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)

				data = append(data, obj)
				fmt.Println(data)

				response := make(map[string]interface{})
				response["summary_circle_meet"] = circleMeet
				response["summary_circles"] = circleVisit
				response["summary_villagevisit"] = villageVisit
				response["summary_beehive"] = beehive
				response["summary_enroll"] = enrolled
				response["summary_green"] = summaryGreen
				response["summary_vyapar"] = summaryVyapar
				response["data"] = data
				response["code"] = 200
				response["success"] = true
				response["message"] = "Successfully"

				js, err := json.Marshal(response)
				fmt.Println(response)
				if err != nil {
					// handle error
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(js)
				return
			}

		} else if reqBody.RoleID == 13 {
			data := []map[string]interface{}{}
			filter := ""
			filterG := " and tp.GreenMotivatorsDate >= CURRENT_DATE()"
			filterV := " and tp.VyaparEnrollmentDate >= CURRENT_DATE()"
			if reqBody.StartDate != "" && reqBody.EndDate != "" {
				filter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterG = " and tp.GreenMotivatorsDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
				filterV = " and tp.VyaparEnrollmentDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "'"
			}
			// Additional code for roleId 13 can be added here
			filters := ""
			str := strconv.Itoa(reqBody.ProjectID)
			if str != "" {
				filters = " and p.id = " + str
			}
			// var request string

			// f := ""
			// if gfId, ok := request["gfId"]; ok && gfId != "" {
			// 	f = " and id=" + gfId
			// }

			summarycircleMeet := 0
			summarycircles := 0
			summaryvillagevisit := 0
			summarybeehive := 0
			summaryenroll := 0
			summaryGreen := 0
			summaryVyapar := 0
			summarycircle_meet := 0

			em, err := db.Query("SELECT id from employee e WHERE status =1 AND  e.supervisorId = ?", reqBody.EmpID)
			if err != nil {
				// handle error
			}
			ids := []int{}
			for em.Next() {
				var id int
				err := em.Scan(&id)
				if err != nil {
					// handle error
				}
				ids = append(ids, id)
			}

			getProjs := "Select project_id as id,p.projectName as name from tbl_poa tp " +
				"inner join project p on p.id = tp.project_id " +
				"where  p.gfl_id = ?" + filters + " GROUP by tp.project_id"

			projectsList, err := db.Query(getProjs, reqBody.EmpID)
			if err != nil {
				// handle error
			}
			participantFilter := ""
			for projectsList.Next() {
				var id int
				var name string
				err := projectsList.Scan(&id, &name)
				if err != nil {
					// handle error
				}

				var prjFilter string
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				circleProjMeet := getGFDataN(db, prjFilter, 1, ids)

				obj := make(map[string]interface{})
				obj["name"] = name
				obj["circle_meet"] = circleProjMeet
				summarycircleMeet += circleProjMeet

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and p.endDate BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				obj["circles"] = getGFCircleN(db, prjFilter, ids)
				summarycircles += obj["circles"].(int)

				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var villageProjvisit interface{}
				villageProjvisit = getGFDataN(db, prjFilter, 2, ids)
				obj = make(map[string]interface{})
				obj["villagevisit"] = villageProjvisit
				if villageProjvisit != nil {
					summaryvillagevisit += villageProjvisit.(int)
				} else {
					obj["villagevisit"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.date BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var beehiveProj interface{}
				beehiveProj = getGFDataN(db, prjFilter, 3, ids)
				obj["beehive"] = beehiveProj
				if beehiveProj != nil {
					summarybeehive += beehiveProj.(int)
				} else {
					obj["beehive"] = "0"
				}
				if reqBody.StartDate != "" && reqBody.EndDate != "" {
					prjFilter = " and tp.participant_day2 BETWEEN '" + reqBody.StartDate + "' and '" + reqBody.EndDate + "' and p.id = " + strconv.Itoa(id)
				} else {
					prjFilter = " and p.id = " + strconv.Itoa(id)
				}
				var projEnrolled interface{}
				b, err := strconv.ParseBool(participantFilter)
				reqBody.StartDate = "2023-03-15T15:30:45Z"

				// layout of the time string
				layout := ""

				// parse the time string into a time.Time object
				t, err := time.Parse(layout, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}
				reqBody.EndDate = ""

				// layout of the time string
				lay := "2006-01-02T15:04:05Z"

				// parse the time string into a time.Time object
				u, err := time.Parse(lay, reqBody.StartDate)
				if err != nil {
					fmt.Println("Error parsing time string:", err)
					return
				}

				// print the time.Time object
				fmt.Println("Parsed time:", u)
				if b {
					if reqBody.StartDate != "" && reqBody.EndDate != "" {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					} else {
						projEnrolled = getParticipantFilterGfEnrolledN(db, prjFilter, ids, t, u)
					}
				} else {
					projEnrolled = getGfEnrolledN(db, prjFilter, ids)
				}
				obj["enroll"] = projEnrolled
				if projEnrolled != nil {
					summaryenroll += projEnrolled.(int)
				} else {
					obj["enroll"] = "0"
				}
				project_result, err := db.Query("SELECT GROUP_CONCAT(DISTINCT prj.id) as ids from project_emps em_pr left join project prj on em_pr.project_id = prj.id where emp_id IN (" + strconv.Itoa(reqBody.EmpID) + ")")

				if err != nil {
					// handle error
				}
				var projectArray string
				if project_result.Next() {
					var ids interface{}
					err := project_result.Scan(&ids)
					if err != nil {
						// handle error
					}
					projectArray = ids.(string)
				}
				parts := strings.Split(projectArray, "")
				nums := make([]int, len(parts))
				for i, p := range parts {
					num, err := strconv.Atoi(p)
					if err != nil {
						panic(err)
					}
					nums[i] = num
				}
				// fmt.Println(nums)
				obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, nums, filterG, "")
				obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, nums, filterV, "")
				summaryGreen += obj["greenMotivators"].(int)
				summaryVyapar += obj["vyapar"].(int)
				data = append(data, obj)
			}
			fmt.Println(data)
			response := make(map[string]interface{})
			response["summary_circle_meet"] = summarycircle_meet
			response["summary_circles"] = summarycircles
			response["summary_villagevisit"] = summaryvillagevisit
			response["summary_beehive"] = summarybeehive
			response["summary_enroll"] = summaryenroll
			response["summary_green"] = summaryGreen
			response["summary_vyapar"] = summaryVyapar
			response["data"] = data
			response["code"] = 200
			response["success"] = true
			response["message"] = "Successfully"
			fmt.Println(response)
			// json.NewEncoder(w).Encode(response)
			json.NewEncoder(w).Encode(map[string]interface{}{"funder": response})

			fmt.Println(filter)

		} else {
			w.WriteHeader(http.StatusCreated)
			response := make(map[string]interface{})
			response["success"] = false
			response["message"] = "Invalid role id"
			json.NewEncoder(w).Encode(response)
		}
		fmt.Println(data...)
		fmt.Println(villagesArray...)
		fmt.Println(summaryWomen)
		fmt.Println(summaryDay2)

	})

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8081", handler)
}
