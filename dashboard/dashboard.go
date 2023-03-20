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

// "end_date": "",
// "role_id": 6,
// "taluk_id": "",
// "district_id": "",
// "trainerId": "",
// "emp_id": 88,
// "start_date": "",
// "somId": "",
// "gflId": "",
// "funder_id": "",
// "partner_id": "",
// "project_id": "",
// "opsManager": ""

type SSdashboard struct {
	EmpRole             int
	ProjectId           string
	StartDate           string
	EndDate             string
	EmpID               int
	TalukID             string
	DistrictID          string
	FunderID            string
	PartnerID           int
	TrainerID           string
	GalathiID           string
	OpsManager          string
	SOMID               string
	GFLID               string
	IsDateFilterApplied bool
}

// func getFunderList(con *sql.DB, funderListQuery string, startDate string, endDate string, dateFilter string, filter string, GalathiId int) ([]map[string]interface{}, error) {
// 	data := make([]map[string]interface{}, 0)
// 	summaryTarget := 0
// 	summaryActuals := 0
// 	summaryDay1 := 0
// 	summaryEnrolled := 0
// 	summaryVillages := 0
// 	summaryGreen := 0
// 	summaryVyapar := 0

// 	if len(funderListQuery) > 0 {
// 		res, err := con.Query(funderListQuery)
// 		if err != nil {
// 			return nil, err
// 		}
// 		defer res.Close()

// 		for res.Next() {
// 			var projectArray []int
// 			obj := make(map[string]interface{})
// 			var funderId int
// 			var funderName string
// 			err = res.Scan(&funderId, &funderName)
// 			if err != nil {
// 				return nil, err
// 			}

// 			getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter

// 			if startDate != "" && endDate != "" {
// 				getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + startDate + "' BETWEEN startDate and endDate and '" + endDate + "' BETWEEN startDate and endDate"
// 			}

// 			projResult, err := con.Query(getProj)
// 			if err != nil {
// 				return nil, err
// 			}
// 			defer projResult.Close()

// 			for projResult.Next() {
// 				var projectId int
// 				err = projResult.Scan(&projectId)
// 				if err != nil {
// 					return nil, err
// 				}
// 				projectArray = append(projectArray, projectId)
// 			}

// 			if len(projectArray) == 0 {
// 				obj["id"] = funderId
// 				obj["name"] = funderName
// 				obj["target"] = 0
// 				obj["actual"] = 0
// 				obj["day2"] = 0
// 				obj["women"] = 0
// 				obj["enrolled"] = 0
// 				obj["villages"] = 0
// 				obj["startDate"] = ""
// 				obj["endDate"] = ""
// 				obj["select_type"] = "2"
// 				data = append(data, obj)
// 				continue
// 			}

// 			obj["id"] = funderId
// 			obj["name"] = funderName
// 			obj["target"] = getTarget(con, startDate, endDate, projectArray)
// 			summaryTarget += obj["target"].(int)
// 			obj["actual"] = getActual(con, startDate, endDate, projectArray, "")
// 			summaryActuals += obj["actual"].(int)

// 			day1Count := getDay1Count(con, startDate, endDate, projectArray, "")
// 			summaryDay1 += day1Count
// 			if day1Count > 0 {
// 				day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
// 				obj["day2"] = int(math.Round(day2Turnout * 100))
// 			} else {
// 				obj["day2"] = 0
// 			}

// 			obj["women"] = obj["actual"]
// 			obj["enrolled"] = getGelathi(con, startDate, endDate, projectArray, GalathiId, funderId,"")
// 			summaryEnrolled += obj["enrolled"].(int)
// 			obj["villages"] = getVillages(con, startDate, endDate, projectArray,"")
// 			summaryVillages += obj["villages"].(int)

// 			obj["green"] = getGreenCount(con, startDate, endDate, projectArray)
// 			summaryGreen += obj["green"].(int)

// 			obj["vyapar"] = getVyaparCount(con, startDate, endDate, projectArray)
// 			summaryVyapar += obj["vyapar"].(int)

// 			obj["select_type"] = "1"
// 			obj["startDate"] = startDate
// 			obj["endDate"] = endDate

// 			data = append(data, obj)
// 		}

// 		if len(data) > 0 {
// 			summary := make(map[string]interface{})
// 			summary["id"] = 0
// 			summary["name"] = "Summary"
// 			summary["target"] = summaryTarget
// 			summary["actual"] = summaryActuals
// 			if summaryDay1 > 0 {
// 				summary["day2"] = int(math.Round(float64(summaryActuals) / float64(summaryDay1) * 100))
// 			} else {
// 				summary["day2"] = 0
// 			}
// 			summary["women"] = summaryActuals
// 			summary["enrolled"] = summaryEnrolled
// 			summary["villages"] = summaryVillages
// 			summary["green"] = summaryGreen
// 			summary["vyapar"] = summaryVyapar
// 			summary["select_type"] = "2"
// 			summary["startDate"] = startDate
// 			summary["endDate"] = endDate
// 			data = append(data, summary)
// 		}
// 	}

// 	return data, nil
// }

func HandleFunc() {
	db, err := sql.Open("mysql", "bdms_staff_admin:sfhakjfhyiqundfgs3765827635@tcp(buzzwomendatabase-new.cixgcssswxvx.ap-south-1.rds.amazonaws.com:3306)/bdms_staff?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database")
	defer db.Close()
	mux := http.NewServeMux()

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method != http.MethodPost {
			w.WriteHeader(405) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
			return
		}
		var dash SSdashboard
		err := json.NewDecoder(r.Body).Decode(&dash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Invalid Input Syntax", "Status Code": "400 ", "Error": err})
			return
		}
		i, err := strconv.Atoi(dash.ProjectId)
		// projectArray, _ := getAssociatedProjectList(db, i)

		projectArray, _ := getAssociatedProjectList(db, i)
		// if dash.EmpRole == 1 || dash.EmpRole == 2 || dash.EmpRole == 3 || dash.EmpRole == 5 || dash.EmpRole == 6 || dash.EmpRole == 11 || dash.EmpRole == 12 || dash.EmpRole == 13 {
		if dash.EmpRole == 1 || dash.EmpRole == 2 || dash.EmpRole == 3 || dash.EmpRole == 11 || dash.EmpRole == 12 {
			day1Count := getDay1Count(db, "", "", projectArray, "")
			summaryDay1 := 0
			var day2 int
			summaryDay1 += day1Count
			if day1Count > 0 {
				actual := getActual(db, dash.StartDate, dash.EndDate, projectArray, "")
				day2Turnout := float64(actual) / float64(day1Count)
				day2 = int(day2Turnout * 100)
			} else {
				day2 = 0
			}
			// var projectStrings []string
			// for _, projectID := range projectArray {
			// 	projectStrings = append(projectStrings, strconv.Itoa(projectID))
			// }

			json.NewEncoder(w).Encode(map[string]interface{}{"target": getTarget(db, dash.StartDate, dash.EndDate, projectArray), "actual": getActual(db, dash.StartDate, dash.EndDate, projectArray, ""), "no of villages": getVillages(db, dash.StartDate, dash.EndDate, projectArray, ""), "day2 turnout": day2, "greenmotivators": greenMotivators(db, dash.StartDate, dash.EndDate, projectArray, "", ""), "Enrolled gelathis": getGelathi(db, dash.StartDate, dash.EndDate, projectArray, "", "", ""), "Enrolled Vypar": Vyapar(db, dash.StartDate, dash.EndDate, projectArray, "", ""), "Women": getActual(db, dash.StartDate, dash.EndDate, projectArray, ""), "no of vyapar module completed": GetNoofVyaparModuleCompleted(db), "no of vyapar survey": GetNoOfVyaparSurvey(db, dash.StartDate, dash.EndDate, "")})

			// json.NewEncoder(w).Encode(map[string]interface{}{"no of vyapar survey": GetNoOfVyaparSurvey(db, dash.StartDate, dash.EndDate, ""), "day2 turnout": day2})

		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Invalid role", "Status Code": "400 Bad Request"})
			return
		}

	})

	mux.HandleFunc("/funder", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,application/json ")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method != http.MethodPost {
			w.WriteHeader(405) // Return 405 Method Not Allowed.
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "method not found", "Status Code": "405 "})
			return
		}
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		decoder := json.NewDecoder(r.Body)

		var dash struct {
			PartnerID  int    `json:"partnerId"`
			StartDate  string `json:"startDate"`
			EndDate    string `json:"endDate"`
			EmpRole    int    `json:"empRole"`
			DistrictID string `json:"districtId"`
			TalukID    string `json:"talukId"`
			TrainerID  string `json:"trainerId"`
			OpsManager string `json:"opsManager"`
			SOMID      string `json:"somid"`
			GFLID      string `json:"gflid"`
		}
		err := decoder.Decode(&dash)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Invalid Input Syntax", "Status Code": "400 ", "Error": err})
			return
		}
		if dash.EmpRole == 1 || dash.EmpRole == 2 || dash.EmpRole == 3 || dash.EmpRole == 5 || dash.EmpRole == 6 || dash.EmpRole == 11 || dash.EmpRole == 12 || dash.EmpRole == 13 {

			var dateFilter, filter, funderListQuery string
			var params []interface{}
			var isDateFilterApplied bool

			if isDateFilterApplied {
				dateFilter = " startDate >= ? AND endDate <= ?"
				params = append(params, dash.StartDate, dash.EndDate)
			} else {
				dateFilter = " endDate >= CURRENT_DATE()"
			}

			if dash.PartnerID > 0 {
				funderListQuery = "SELECT DISTINCT(p.funderId) as id, funderName as name FROM project p " +
					"INNER JOIN funder ON funder.funderID = p.funderID " +
					"WHERE p.partnerID = ? AND " + dateFilter + filter
				params = append(params, dash.PartnerID)
				filter += " AND p.partnerID = ?"
				params = append(params, dash.PartnerID)

				rows, err := db.Query(funderListQuery, params...)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Internal Server Error", "Status Code": "500 ", "Error": err})
					return
				}
				defer rows.Close()

				funders := make([]Funder, 0)
				for rows.Next() {
					var funder Funder
					err := rows.Scan(&funder.ID, &funder.Name)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Internal Server Error", "Status Code": "500 ", "Error": err})
						return
					}
					funders = append(funders, funder)
				}

				json.NewEncoder(w).Encode(funders)

			}
			// fmt.Println(funders)

		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"Message": "Invalid role", "Status Code": "400 Bad Request"})
			return
		}

	})
	mux.HandleFunc("/anand", func(w http.ResponseWriter, r *http.Request) {
		// isDateFilterApplied := true
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		// decoder := json.NewDecoder(r.Body)
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
			GalathiID  int    `json:"galathi_id"`
		}

		var reqBody ProjectRequest
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing request body: %v", err)
			return
		}
		isDateFilterApplied := false
		var dateFilter string
		if reqBody.StartDate != "" && reqBody.EndDate != "" {
			isDateFilterApplied = true
			dateFilter = fmt.Sprintf("startDate >= '%s' AND endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
		} else {
			dateFilter = "endDate >= CURRENT_DATE()"
		}

		// Define the SQL query
		// Build query based on request parameters
		var funderListQuery string
		var filter string

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

		// Execute the query
		rows, err := db.Query(funderListQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error executing query: %v", err)
			return
		}
		defer rows.Close()

		// Store the results in a slice of Funders
		funderList := make([]Funder, 0)
		for rows.Next() {
			var funder Funder
			err := rows.Scan(&funder.ID, &funder.Name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error scanning funder list: %v", err)
				return
			}
			funderList = append(funderList, funder)
		}

		// Encode the slice of Funders to JSON and send it as response
		w.Header().Set("Content-Type", "application/json")
		projectArray, _ := getAssociatedProjectList(db, reqBody.PartnerID)
		// json.NewEncoder(w).Encode(map[string]interface{}{"asd": GetNoofVyaparModuleCompleted(db)})
		// json.NewEncoder(w).Encode(map[string]interface{}{"no of vyapar enrolled": Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""), "no of villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "asd": funderList})
		json.NewEncoder(w).Encode(map[string]interface{}{"no of villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "no of vyapar enrolled": Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""), "no of vyapar survey": GetNoOfVyaparSurvey(db, reqBody.StartDate, reqBody.EndDate, ""), "no of vyapar module completed": GetNoofVyaparModuleCompleted(db), "funder lis": funderList})

	})

	mux.HandleFunc("/api/funders", func(w http.ResponseWriter, r *http.Request) {
		type Funder struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		// decoder := json.NewDecoder(r.Body)
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
			// Program Manager
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
				if reqBody.ProjectID > 0 {
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

			// Define the SQL query
			// Build query based on request parameters
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
					// projectList := ""
					// if projectId > 0 {
					// 	dateFilterNew := ""
					// 	if isDateFilterApplied {
					// 		dateFilterNew = " and startDate >= '" + startDate + "' and endDate <= '" + endDate + "'"
					// 	}
					// 	projectList = "SELECT id, projectName as name, startDate, endDate from project p where id = " + strconv.Itoa(projectId) + filter + dateFilterNew
					// 	summaryProjectsArray = append(summaryProjectsArray, projectId)
					// } else if trainerId > 0 {
					// 	projectList = "SELECT project_id as id, projectName as name, p.startDate, p.endDate from tbl_poa tp inner join project p on p.id = tp.project_id where user_id = " + strconv.Itoa(trainerId) + " and " + dateFilter + filter + " GROUP  by project_id"
					// 	summaryFilter = " and tp.user_id = " + strconv.Itoa(trainerId)
					// } else if opsManager > 0 {
					// 	if dateFilter == "" || (startDate == "" && endDate == "") {
					// 		projectList = "SELECT id, projectName as name, startDate, endDate from project p where operations_manager = " + strconv.Itoa(opsManager) + " and " + dateFilter + filter + " GROUP by id "
					// 	} else {
					// 		projectList = "SELECT p.id, p.projectName as name, p.startDate, p.endDate from project p join training_participants tp on p.id = tp.project_id where p.operations_manager = " + strconv.Itoa(opsManager) + " and tp.participant_day2 >= '" + startDate + "' and tp.participant_day2 <= '" + endDate + "' GROUP by p.id "
					// 	}
					// 	summaryFilter = " and p.operations_manager = " + strconv.Itoa(opsManager)
					// } else if somId > 0 {
					// 	projectList = "SELECT id, projectName as name, startDate, endDate from project p where operations_manager in(SELECT id from employee e where e.supervisorId =" + strconv.Itoa(somId) + ") and " + dateFilter + filter + " GROUP by id "
					// 	summaryFilter = " and p.operations_manager in (SELECT id from employee e where e.supervisorId =" + strconv.Itoa(somId) + ")"
					// } else if gflId > 0 {
					// 	projectList = "SELECT id, projectName as name, startDate, endDate from project p where operations_manager in(SELECT supervisorId from employee e where e.id =" + strconv.Itoa(gflId) + ") and " + dateFilter + filter + " GROUP by id "
					// 	summaryFilter = " and p.operations_manager in (SELECT supervisorId from employee e where e.id =" + strconv.Itoa(gflId) + ")"
					// } else if isDateFilterApplied && partnerId == 0 && dist == 0 && funderId == 0 || roleId == 4 && dist == 0 {
					// 	projectList = "SELECT id, projectName as name, startDate, endDate from project p where " + dateFilter + filter
					// }
					// if len(projectList) > 0 {
					// 	res, err := con.Query(projectList)
					// 	if err != nil {
					// 		// handle error
					// 	}
					// 	defer res.Close()
					// 	for res.Next() {
					// 		var obj map[string]interface{}
					// 		var projectArray []int
					// 		var tpFilter, tbFilter string
					// 		var day2Turnout float64
					// 		prList := make(map[string]interface{})
					// 		if err := res.Scan(&prList["id"], &prList["name"], &prList["startDate"], &prList["endDate"]); err != nil {
					// 			// handle error
					// 		}
					// 		obj["id"] = prList["id"]
					// 		obj["name"] = prList["name"]
					// 		projectArray = append(projectArray, prList["id"].(int))
					// 		if trainerId > 0 {
					// 			obj["target"] = getTrainerTarget(con, trainerId, projectArray)
					// 			summaryTarget += obj["target"].(int)
					// 			tpFilter = fmt.Sprintf(" and tp.trainer_id = %d", trainerId)
					// 			tbFilter = fmt.Sprintf(" and tp.user_id = %d", trainerId)
					// 		} else {
					// 			obj["target"] = getTarget(con, startDate, endDate, projectArray)
					// 			summaryTarget += obj["target"].(int)
					// 		}
					// 		obj["actual"] = getActual(con, startDate, endDate, projectArray, tpFilter)
					// 		summaryActuals += obj["actual"].(int)
					// 		day1Count := getDay1Count(con, startDate, endDate, projectArray, tpFilter)
					// 		summaryDay1 += day1Count
					// 		if day1Count > 0 {
					// 			day2Turnout = float64(obj["actual"].(int)) / float64(day1Count)
					// 			obj["day2"] = int(math.Round(day2Turnout * 100))
					// 		} else {
					// 			obj["day2"] = 0
					// 		}
					// 		obj["women"] = obj["actual"]
					// 		obj["enrolled"] = getGelathi(con, startDate, endDate, projectArray, tpFilter)
					// 		summaryEnrolled += obj["enrolled"].(int)
					// 		obj["villages"] = newVillageCount(con, startDate, endDate, projectArray, tbFilter)
					// 		summaryVillages += obj["villages"].(int)
					// 		obj["startDate"] = prList["startDate"]
					// 		obj["endDate"] = prList["endDate"]
					// 		obj["select_type"] = "1"
					// 		//green motivators
					// 		obj["greenMotivators"] = greenMotivators(con, startDate, endDate, projectArray, tpFilter)
					// 		obj["vyapar"] = Vyapar(con, startDate, endDate, projectArray, tpFilter)
					// 		summaryGreen += obj["greenMotivators"].(int)
					// 		summaryVyapar += obj["vyapar"].(int)

					// 		data = append(data, obj)
					// 	}
					// }
					// response := make(map[string]interface{})
					// response["summary_target"] = summaryTarget
					// response["summary_women"] = summaryActuals
					// response["summary_villages"] = getSummaryOfVillagesNew(con, startDate, endDate, summaryProjectsArray, filter) //new anas
					// response["summary_villages_count"] = summaryVillages //new anas
					// if summaryDay1 > 0 {
					// 	day2TurnOut := float64(summaryActuals) / float64(summaryDay1)
					// 	response["summary_day2"] = int(math.Round(day2TurnOut * 100))
					// } else {
					// 	response["summary_day2"] = 0
					// }
					// response["summary_enrolled"] = summaryEnrolled
					// response["summary_green"] = summary

					projectList := ""

					summaryFilter := ""

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
							// id := 0
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
							// var summaryTarget int = 0
							// var summaryActuals int = 0
							// var summaryEnrolled int = 0
							// var summaryVillages int = 0
							// var summaryDay1 int = 0
							// var summaryGreen int = 0
							// var summaryVyapar int = 0

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

						data = append(data, obj)
					}
					// json.NewEncoder(w).Encode(map[string]interface{}{"No of Vypar Cohorts": NoofVyaparCohorts(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Villages": getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, ""), "No Of Vypar Enrolled Vypar": Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", ""), "No Of Vyapar survey": GetNoOfVyaparSurvey(db, reqBody.StartDate, reqBody.EndDate, ""), "No Of Vyapar module completed": GetNoofVyaparModuleCompleted(db), "funder": data})
					json.NewEncoder(w).Encode(map[string]interface{}{"funder": data})
				}

				// fmt.Println(data...)

			}
			fmt.Println(summaryFilter)
			fmt.Println(funderList)
			// } else if reqBody.RoleID == 5 {
			// 	var dateFilter string
			// 	var isDateFilterApplied bool
			// 	if isDateFilterApplied {
			// 		dateFilter = " and p.startDate >= '" + reqBody.StartDate + "' and p.endDate <= '" + reqBody.EndDate + "'"
			// 	} else {
			// 		dateFilter = " and p.endDate >= CURRENT_DATE()"
			// 	}

			// 	var query string
			// 	data := []map[string]interface{}{}
			// 	// var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			// 	if reqBody.ProjectID > 0 {
			// 		query = "SELECT COALESCE(ISNULL(project_id, 0), 0) AS id, COALESCE(ISNULL(projectName, ''), '') AS name, COALESCE(ISNULL(p.startDate, '1900-01-01'), '1900-01-01') AS startDate, COALESCE(ISNULL(p.endDate, '9999-12-31'), '9999-12-31') AS endDate FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE user_id = " + strconv.Itoa(reqBody.EmpID) + " AND tp.project_id = " + strconv.Itoa(reqBody.ProjectID) + " GROUP BY tp.project_id"
			// 		summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			// 	} else {
			// 		query = "SELECT COALESCE(project_id, 0) AS id, COALESCE(projectName, '') AS name, COALESCE(p.startDate, '1900-01-01') AS startDate, COALESCE(p.endDate, '9999-12-31') AS endDate FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE user_id = " + strconv.Itoa(reqBody.EmpID) + dateFilter + " GROUP BY project_id"
			// 	}

			// 	res, err := db.Query(query)

			// 	if err != nil {
			// 		log.Fatal(err)
			// 	}

			// 	for res.Next() {
			// 		var id int
			// 		var name string
			// 		var startDate, endDate string

			// 		err := res.Scan(&id, &name, &startDate, &endDate)

			// 		if err != nil {
			// 			log.Fatal(err)
			// 		}

			// 		projectArray := []int{id}
			// 		obj := make(map[string]interface{})

			// 		obj["id"] = id
			// 		obj["name"] = name
			// 		obj["startDate"] = startDate
			// 		obj["endDate"] = endDate
			// 		obj["select_type"] = "1"

			// 		target := getTrainerTarget(db, reqBody.EmpID, projectArray)
			// 		obj["target"] = target
			// 		summaryTarget += target

			// 		filter := " and tp.trainer_id = " + strconv.Itoa(reqBody.EmpID)
			// 		actual := getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
			// 		obj["actual"] = actual
			// 		summaryActuals += actual

			// 		day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
			// 		summaryDay1 += day1Count

			// 		if day1Count > 0 {
			// 			day2TurnOut := float64(actual) / float64(day1Count)
			// 			obj["day2"] = int(math.Round(day2TurnOut * 100))
			// 		} else {
			// 			obj["day2"] = 0
			// 		}

			// 		obj["women"] = obj["actual"]
			// 		obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", "")
			// 		var tbFilter string

			// 		summaryEnrolled += obj["enrolled"].(int)
			// 		tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			// 		strSlice := make([]string, len(projectArray))

			// 		// loop through each element in intSlice and convert to string
			// 		for i, v := range projectArray {
			// 			strSlice[i] = strconv.Itoa(v)
			// 		}
			// 		obj["villages"] = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, tbFilter)
			// 		summaryVillages += obj["villages"].(int)
			// 		obj["startDate"] = obj["startDate"]
			// 		obj["endDate"] = obj["endDate"]
			// 		obj["select_type"] = "1"
			// 		obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
			// 		obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
			// 		summaryGreen += obj["greenMotivators"].(int)
			// 		summaryVyapar += obj["vyapar"].(int)
			// 		data = append(data, obj)
			// 		response := make(map[string]interface{})

			// 		response["summary_target"] = summaryTarget
			// 		response["summary_women"] = summaryActuals
			// 		tbFilter = fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			// 		intSlice := []int{}

			// 		// loop through each element in the []interface{} slice
			// 		for _, v := range summaryProjectsArray {
			// 			// check if the element is of type int
			// 			if i, ok := v.(int); ok {
			// 				// append the int value to the []int slice
			// 				intSlice = append(intSlice, i)
			// 			}
			// 		}
			// 		response["summary_villages"] = getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter)
			// 		response["summary_actual"] = summaryActuals
			// 		var day2Turnout float64

			// 		if summaryDay1 > 0 {
			// 			day2Turnout = float64(summaryActuals) / float64(summaryDay1)
			// 			response["summary_day2"] = int(math.Round(day2Turnout * 100))
			// 		} else {
			// 			day2Turnout = 0
			// 			response["summary_day2"] = 0
			// 		}

			// 		response["summary_enrolled"] = summaryEnrolled
			// 		response["summary_green"] = summaryGreen
			// 		response["summary_vyapar"] = summaryVyapar
			// 		response["data"] = data
			// 		response["code"] = 200
			// 		response["success"] = true
			// 		response["message"] = "Successfully"

			// 		jsonResponse, err := json.Marshal(response)
			// 		if err != nil {
			// 			log.Fatal(err)
			// 		}
			// 		fmt.Println(string(jsonResponse))
			// 	}

			// } else if reqBody.RoleID == 5 {
			// 	var dateFilter string
			// 	var isDateFilterApplied bool
			// 	if isDateFilterApplied {
			// 		dateFilter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
			// 	} else {
			// 		dateFilter = " and p.endDate >= CURRENT_DATE()"
			// 	}

			// 	query := fmt.Sprintf("SELECT COALESCE(project_id, 0) AS id, COALESCE(projectName, '') AS name, COALESCE(p.startDate, '1970-01-01') AS startDate, COALESCE(p.endDate, '2100-12-31') AS endDate FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE user_id = %d%s GROUP BY project_id", reqBody.EmpID, dateFilter)
			// 	if reqBody.ProjectID > 0 {
			// 		query = fmt.Sprintf("SELECT COALESCE(project_id, 0) as id, COALESCE(projectName, '') as name, COALESCE(p.startDate, '') as startDate, COALESCE(p.endDate, '') as endDate FROM tbl_poa tp INNER JOIN project p ON p.id = tp.project_id WHERE user_id = %d AND tp.project_id = %d GROUP BY tp.project_id", reqBody.EmpID, reqBody.ProjectID)
			// 		summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			// 	}

			// 	rows, err := db.Query(query)
			// 	fmt.Println(rows)
			// 	if err != nil {
			// 		// handle error
			// 	}
			// 	defer rows.Close()

			// 	var data []map[string]interface{}
			// 	var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			// 	for rows.Next() {
			// 		// var id sql.NullInt64
			// 		var name string
			// 		var startDate, endDate string
			// 		var id sql.NullInt64
			// 		var i int
			// 		if id.Valid {
			// 			i = int(id.Int64)
			// 		} else {
			// 			i = 0
			// 		}

			// 		if err := rows.Scan(&i, &name, &startDate, &endDate); err != nil {
			// 			// handle error
			// 		}

			// 		projectArray := []int{i}
			// 		obj := make(map[string]interface{})
			// 		obj["id"] = id
			// 		obj["name"] = name
			// 		obj["target"] = getTrainerTarget(db, reqBody.EmpID, projectArray)
			// 		summaryTarget += obj["target"].(int)
			// 		filter := fmt.Sprintf(" and tp.trainer_id = %d", reqBody.EmpID)
			// 		obj["actual"] = getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
			// 		summaryActuals += obj["actual"].(int)
			// 		day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, filter)
			// 		summaryDay1 += day1Count
			// 		if day1Count > 0 {
			// 			day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
			// 			obj["day2"] = math.Round(day2Turnout * 100)
			// 		} else {
			// 			obj["day2"] = 0
			// 		}
			// 		obj["women"] = obj["actual"]
			// 		obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "", filter)
			// 		summaryEnrolled += obj["enrolled"].(int)
			// 		tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			// 		strSlice := make([]string, len(projectArray))

			// 		// loop through each element in intSlice and convert to string
			// 		for i, v := range projectArray {
			// 			strSlice[i] = strconv.Itoa(v)
			// 		}
			// 		obj["villages"] = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, strSlice, tbFilter)
			// 		summaryVillages += obj["villages"].(int)
			// 		obj["startDate"] = reqBody.StartDate
			// 		obj["endDate"] = reqBody.EndDate
			// 		obj["select_type"] = "1"
			// 		// green motivators
			// 		obj["greenMotivators"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
			// 		obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", filter)
			// 		summaryGreen += obj["greenMotivators"].(int)
			// 		summaryVyapar += obj["vyapar"].(int)
			// 		// obj = map[string]interface{}{
			// 		// 	"greenMotivators": greenMotivators,
			// 		// 	"vyapar":          vyapar,
			// 		// }
			// 		data = append(data, obj)
			// 		fmt.Println(data)
			// 		intSlice := []int{}

			// 		// loop through each element in the []interface{} slice
			// 		for _, v := range summaryProjectsArray {
			// 			// check if the element is of type int
			// 			if i, ok := v.(int); ok {
			// 				// append the int value to the []int slice
			// 				intSlice = append(intSlice, i)
			// 			}
			// 		}
			// 		response := map[string]interface{}{
			// 			"summary_target":   summaryTarget,
			// 			"summary_women":    summaryActuals,
			// 			"summary_villages": getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, intSlice, tbFilter),
			// 			"summary_actual":   summaryActuals,
			// 		}
			// 		if summaryDay1 > 0 {
			// 			day2TurnOut := summaryActuals / summaryDay1
			// 			response["summary_day2"] = int(day2TurnOut * 100)
			// 		} else {
			// 			day2TurnOut := 0
			// 			response["summary_day2"] = 0
			// 			fmt.Println(day2TurnOut)
			// 		}
			// 		response["summary_enrolled"] = summaryEnrolled
			// 		response["summary_green"] = summaryGreen
			// 		response["summary_vyapar"] = summaryVyapar
			// 		response["data"] = data
			// 		response["code"] = 200
			// 		response["success"] = true
			// 		response["message"] = "Successfully"
			// 		json.NewEncoder(w).Encode(response)
			// 		fmt.Println(response)
			// 		json.NewEncoder(w).Encode(map[string]interface{}{"funder": data, "summary": response})

			// 	}
			// } else if reqBody.RoleID == 5 {
			// 	var dateFilter string
			// 	var isDateFilterApplied bool
			// 	if isDateFilterApplied {
			// 		dateFilter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate, reqBody.EndDate)
			// 	} else {
			// 		dateFilter = " and p.endDate >= CURRENT_DATE()"
			// 	}

			// 	query := fmt.Sprintf("SELECT project_id as id, projectName as name, p.startDate, p.endDate "+
			// 		"FROM tbl_poa tp "+
			// 		"INNER JOIN project p ON p.id = tp.project_id "+
			// 		"WHERE user_id = %d%s "+
			// 		"GROUP BY project_id", reqBody.EmpID, dateFilter)

			// 	var summaryProjectsArray []int

			// 	if reqBody.ProjectID > 0 {
			// 		query = fmt.Sprintf("SELECT project_id as id, projectName as name, p.startDate, p.endDate "+
			// 			"FROM tbl_poa tp "+
			// 			"INNER JOIN project p ON p.id = tp.project_id "+
			// 			"WHERE user_id = %d AND tp.project_id = %d "+
			// 			"GROUP BY tp.project_id", reqBody.EmpID, reqBody.ProjectID)
			// 		summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			// 	}

			// 	res, err := db.Query(query)

			// 	if err != nil {
			// 		// handle error
			// 	}
			// 	defer res.Close()

			// 	var data []map[string]interface{}
			// 	var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int

			// 	for res.Next() {
			// 		type project struct {
			// 			id              int
			// 			name            string
			// 			startDate       time.Time
			// 			endDate         time.Time
			// 			day2            int
			// 			women           int
			// 			enrolled        int
			// 			villages        int
			// 			greenMotivators int
			// 			vyapar          int
			// 			select_type     int
			// 		}

			// 		var projects []project
			// 		var projectIDs []int
			// 		var target, actual, day2Turnout, women, enrolled, villages, green, vyapar int

			// 		for res.Next() {
			// 			var prj project

			// 			err = res.Scan(&prj.id, &prj.name, &prj.startDate, &prj.endDate)
			// 			if err != nil {
			// 				// handle error
			// 			}

			// 			projects = append(projects, prj)
			// 			projectIDs = append(projectIDs, prj.id)
			// 		}

			// 		target = getTrainerTarget(db, reqBody.EmpID, projectIDs)
			// 		summaryTarget += target

			// 		filter := fmt.Sprintf(" and tp.trainer_id = %d", reqBody.EmpID)
			// 		actual = getActual(db, reqBody.StartDate, reqBody.EndDate, projectIDs, filter)
			// 		summaryActuals += actual

			// 		day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectIDs, filter)
			// 		summaryDay1 += day1Count

			// 		if day1Count > 0 {
			// 			day2Turnout = actual / day1Count
			// 		}

			// 		for i, prj := range projects {
			// 			projects[i].day2 = day2Turnout * 100

			// 			women = actual
			// 			projects[i].women = women

			// 			enrolled = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectIDs, "", "", filter)
			// 			summaryEnrolled += enrolled
			// 			projects[i].enrolled = enrolled

			// 			tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			// 			stringSlice := make([]string, len(projectIDs))

			// 			for i, v := range projectIDs {
			// 				stringSlice[i] = strconv.Itoa(v)
			// 			}

			// 			villages = newVillageCount(db, reqBody.StartDate, reqBody.EndDate, stringSlice, tbFilter)
			// 			summaryVillages += villages
			// 			projects[i].villages = villages

			// 			projects[i].select_type = 1

			// 			green = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectIDs, "", filter)
			// 			summaryGreen += green
			// 			projects[i].greenMotivators = green

			// 			vyapar = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectIDs, "", filter)
			// 			summaryVyapar += vyapar
			// 			projects[i].vyapar = vyapar
			// 			fmt.Println(prj)
			// 		}
			// 		data := make([]map[string]interface{}, 0, len(projects))

			// 		for _, prj := range projects {
			// 			m := map[string]interface{}{
			// 				"id":              prj.id,
			// 				"name":            prj.name,
			// 				"startDate":       prj.startDate,
			// 				"endDate":         prj.endDate,
			// 				"day2":            prj.day2,
			// 				"women":           prj.women,
			// 				"enrolled":        prj.enrolled,
			// 				"villages":        prj.villages,
			// 				"select_type":     prj.select_type,
			// 				"greenMotivators": prj.greenMotivators,
			// 				"vyapar":          prj.vyapar,
			// 			}
			// 			data = append(data, m)
			// 		}
			// 		// data =append(data,project...)

			// 		// data = append(data, projects...)
			// 	}

			// 	response := make(map[string]interface{})
			// 	response["summary_target"] = summaryTarget
			// 	response["summary_women"] = summaryActuals
			// 	tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
			// 	response["summary_villages"] = getSummaryOfVillagesNew(db, reqBody.StartDate, reqBody.EndDate, summaryProjectsArray, tbFilter)
			// 	response["summary_actual"] = summaryActuals
			// 	if summaryDay1 > 0 {
			// 		day2TurnOut := summaryActuals / summaryDay1
			// 		response["summary_day2"] = math.Round(float64(day2TurnOut * 100))
			// 	} else {
			// 		day2TurnOut := 0
			// 		response["summary_day2"] = 0
			// 		fmt.Println(day2TurnOut)
			// 	}
			// 	response["summary_enrolled"] = summaryEnrolled
			// 	response["summary_green"] = summaryGreen
			// 	response["summary_vyapar"] = summaryVyapar
			// 	response["data"] = data
			// 	response["code"] = 200
			// 	response["success"] = true
			// 	response["message"] = "Successfully"
			// 	jsonResponse, err := json.Marshal(response)
			// 	if err != nil {
			// 		// handle error
			// 	}
			// 	fmt.Println(string(jsonResponse))
			// 	os.Exit(0)

		} else if reqBody.RoleID == 5 {
			var dateFilter string
			var isDateFilterApplied bool
			if isDateFilterApplied {
				dateFilter = fmt.Sprintf(" and p.startDate >= '%s' and p.endDate <= '%s'", reqBody.StartDate)
			} else {
				dateFilter = " and p.endDate >= CURRENT_DATE()"
			}

			// trainer
			query := fmt.Sprintf("SELECT project_id as id,projectName as name,p.startDate,p.endDate "+
				"from tbl_poa tp "+
				"inner join project p on p.id = tp.project_id "+
				"where user_id = %d %s GROUP by project_id", reqBody.EmpID, dateFilter)

			if reqBody.ProjectID > 0 {
				query = fmt.Sprintf("SELECT project_id as id,projectName as name,p.startDate,p.endDate "+
					"from tbl_poa tp "+
					"inner join project p on p.id = tp.project_id "+
					"where user_id = %d and tp.project_id = %d GROUP by tp.project_id", reqBody.EmpID, reqBody.ProjectID)
				summaryProjectsArray = append(summaryProjectsArray, reqBody.ProjectID)
			}

			res, err := db.Query(query)
			if err != nil {
				// handle error
			}
			defer res.Close()

			var data []interface{}
			summary := make(map[string]interface{})
			var summaryTarget, summaryActuals, summaryDay1, summaryEnrolled, summaryVillages, summaryGreen, summaryVyapar int
			for res.Next() {
				var id, name string
				var startDate, endDate string
				err = res.Scan(&id, &name, &startDate, &endDate)
				if err != nil {
					// handle error
				}

				projectArray := []string{id}
				obj := make(map[string]interface{})
				obj["id"] = id
				obj["name"] = name
				obj["startDate"] = startDate
				obj["endDate"] = endDate
				obj["select_type"] = "1"

				intSlice := make([]int, len(projectArray))

				for i, str := range projectArray {
					num, err := strconv.Atoi(str)
					if err != nil {
						panic(err)
					}
					intSlice[i] = num
				}

				obj["target"] = getTrainerTarget(db, reqBody.EmpID, intSlice)
				if err != nil {
					// handle error
				}
				summaryTarget += obj["target"].(int)

				filter := fmt.Sprintf(" and tp.trainer_id = %d", reqBody.EmpID)
				obj["actual"] = getActual(db, startDate, endDate, intSlice, filter)
				if err != nil {
					// handle error
				}
				summaryActuals += obj["actual"].(int)

				day1Count := getDay1Count(db, startDate, endDate, intSlice, filter)
				if err != nil {
					// handle error
				}
				summaryDay1 += day1Count

				if day1Count > 0 {
					day2Turnout := obj["actual"].(float64) / float64(day1Count)
					obj["day2"] = math.Round(day2Turnout * 100)
				} else {
					obj["day2"] = 0
				}

				obj["women"] = obj["actual"]
				obj["enrolled"] = getGelathi(db, startDate, endDate, intSlice, "", "", filter)
				if err != nil {
					// handle error
				}
				summaryEnrolled += obj["enrolled"].(int)

				tbFilter := fmt.Sprintf(" and tp.user_id = %d", reqBody.EmpID)
				obj["villages"] = newVillageCount(db, startDate, endDate, projectArray, tbFilter)
				if err != nil {

					// handle error
				}
				summaryVillages += obj["villages"].(int)
				obj["green"] = greenMotivators(db, startDate, endDate, intSlice, filter, "")
				if err != nil {
					// handle error
				}
				summaryGreen += obj["green"].(int)

				obj["vyapar"] = Vyapar(db, startDate, endDate, intSlice, filter, "")
				if err != nil {
					// handle error
				}
				summaryVyapar += obj["vyapar"].(int)

				data = append(data, obj)
				fmt.Println(data...)
			}

			summary["target"] = summaryTarget
			summary["actual"] = summaryActuals
			summary["day1"] = summaryDay1
			summary["enrolled"] = summaryEnrolled
			summary["villages"] = summaryVillages
			summary["green"] = summaryGreen
			summary["vyapar"] = summaryVyapar
			json.NewEncoder(w).Encode(map[string]interface{}{"funder": data, "summary": summary})

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

			}
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

	// 	data := make([]map[string]interface{}, 0)
	// 	summaryTarget := 0
	// 	summaryActuals := 0
	// 	summaryDay1 := 0
	// 	summaryEnrolled := 0
	// 	summaryVillages := 0
	// 	summaryGreen := 0
	// 	summaryVyapar := 0

	// 	if len(funderListQuery) > 0 {
	// 		res, err := db.Query(funderListQuery)
	// 		if err != nil {
	// 			return
	// 		}
	// 		defer res.Close()

	// 		for res.Next() {
	// 			var projectArray []int
	// 			obj := make(map[string]interface{})
	// 			var funderId int
	// 			var funderName string
	// 			err = res.Scan(&funderId, &funderName)
	// 			if err != nil {
	// 				return
	// 			}

	// 			getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter

	// 			if reqBody.StartDate != "" && reqBody.EndDate != "" {
	// 				getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
	// 			}

	// 			projResult, err := db.Query(getProj)
	// 			if err != nil {
	// 				return
	// 			}
	// 			defer projResult.Close()

	// 			for projResult.Next() {
	// 				var projectId int
	// 				err = projResult.Scan(&projectId)
	// 				if err != nil {
	// 					return
	// 				}
	// 				projectArray = append(projectArray, projectId)
	// 			}

	// 			if len(projectArray) == 0 {
	// 				obj["id"] = funderId
	// 				obj["name"] = funderName
	// 				obj["target"] = 0
	// 				obj["actual"] = 0
	// 				obj["day2"] = 0
	// 				obj["women"] = 0
	// 				obj["enrolled"] = 0
	// 				obj["villages"] = 0
	// 				obj["startDate"] = ""
	// 				obj["endDate"] = ""
	// 				obj["select_type"] = "2"
	// 				data = append(data, obj)
	// 				continue
	// 			}

	// 			obj["id"] = funderId
	// 			obj["name"] = funderName
	// 			obj["target"] = getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray)
	// 			summaryTarget += obj["target"].(int)
	// 			obj["actual"] = getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, "")
	// 			summaryActuals += obj["actual"].(int)

	// 			day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, "")
	// 			summaryDay1 += day1Count
	// 			if day1Count > 0 {
	// 				day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
	// 				obj["day2"] = int(math.Round(day2Turnout * 100))
	// 			} else {
	// 				obj["day2"] = 0
	// 			}

	// 			obj["women"] = obj["actual"]
	// 			s := strconv.Itoa(funderId)
	// 			obj["enrolled"] = getGelathi(db, reqBody.StartDate, reqBody.EndDate, projectArray, reqBody.GalathiID, s, "")
	// 			summaryEnrolled += obj["enrolled"].(int)
	// 			obj["villages"] = getVillages(db, reqBody.StartDate, reqBody.EndDate, projectArray, "")
	// 			summaryVillages += obj["villages"].(int)

	// 			obj["green"] = greenMotivators(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "")

	// 			obj["vyapar"] = Vyapar(db, reqBody.StartDate, reqBody.EndDate, projectArray, "", "")
	// 			summaryVyapar += obj["vyapar"].(int)

	// 			obj["select_type"] = "1"
	// 			obj["startDate"] = reqBody.StartDate
	// 			obj["endDate"] = reqBody.EndDate

	// 			data = append(data, obj)
	// 		}

	// 		if len(data) > 0 {
	// 			summary := make(map[string]interface{})
	// 			summary["id"] = 0
	// 			summary["name"] = "Summary"
	// 			summary["target"] = summaryTarget
	// 			summary["actual"] = summaryActuals
	// 			if summaryDay1 > 0 {
	// 				summary["day2"] = int(math.Round(float64(summaryActuals) / float64(summaryDay1) * 100))
	// 			} else {
	// 				summary["day2"] = 0
	// 			}
	// 			summary["women"] = summaryActuals
	// 			summary["enrolled"] = summaryEnrolled
	// 			summary["villages"] = summaryVillages
	// 			summary["green"] = summaryGreen
	// 			summary["vyapar"] = summaryVyapar
	// 			summary["select_type"] = "2"
	// 			summary["startDate"] = reqBody.StartDate
	// 			summary["endDate"] = reqBody.EndDate
	// 			data = append(data, summary)
	// 		}
	// 	}

	// 	return data, nil
	// })

	// 	var funderList []map[string]interface{}

	// 	// gives funder list
	// 	if len(funderListQuery) > 0 {
	// 		res, err := db.Query(funderListQuery)
	// 		if err != nil {
	// 			// handle error
	// 		}
	// 		defer res.Close()

	// 		for res.Next() {
	// 			projectArray := []int{}
	// 			var funderId int
	// 			var funderName string
	// 			err := res.Scan(&funderId, &funderName)
	// 			if err != nil {
	// 				// handle error
	// 			}

	// 			getProj := "SELECT id from project p where funderID = " + strconv.Itoa(funderId) + " and " + dateFilter + filter
	// 			if reqBody.StartDate != "" && reqBody.EndDate != "" {
	// 				getProj = "SELECT id, startDate, endDate from project p where funderID = " + strconv.Itoa(funderId) + " and '" + reqBody.StartDate + "' BETWEEN startDate and endDate and '" + reqBody.EndDate + "' BETWEEN startDate and endDate"
	// 			}

	// 			projResult, err := db.Query(getProj)
	// 			if err != nil {
	// 				// handle error
	// 			}
	// 			defer projResult.Close()

	// 			for projResult.Next() {
	// 				var projectId int
	// 				err := projResult.Scan(&projectId)
	// 				if err != nil {
	// 					// handle error
	// 				}
	// 				projectArray = append(projectArray, projectId)
	// 			}

	// 			if len(projectArray) == 0 {
	// 				obj := make(map[string]interface{})
	// 				obj["id"] = funderId
	// 				obj["name"] = funderName
	// 				obj["target"] = 0
	// 				obj["actual"] = 0
	// 				obj["day2"] = 0
	// 				obj["women"] = 0
	// 				obj["enrolled"] = 0
	// 				obj["villages"] = 0
	// 				obj["startDate"] = ""
	// 				obj["endDate"] = ""
	// 				obj["select_type"] = "2"
	// 				data = append(data, obj)
	// 				continue
	// 			}

	// 			obj := make(map[string]interface{})
	// 			obj["id"] = funderId
	// 			obj["name"] = funderName
	// 			obj["target"] = getTarget(db, reqBody.StartDate, reqBody.EndDate, projectArray)
	// 			summaryTarget += obj["target"].(int)
	// 			obj["actual"] = getActual(db, reqBody.StartDate, reqBody.EndDate, projectArray, "")
	// 			summaryActuals += obj["actual"].(int)

	// 			day1Count := getDay1Count(db, reqBody.StartDate, reqBody.EndDate, projectArray, "")
	// 			summaryDay1 += day1Count
	// 			if day1Count > 0 {
	// 				day2Turnout := float64(obj["actual"].(int)) / float64(day1Count)
	// 				obj["day2"] = int(math.Round(day2Turnout * 100))
	// 			} else {
	// 				obj["day2"] = 0
	// 			}
	// 			obj["women"] = obj["actual"]
	// 			obj["enrolled"] = getGelathi(con, startDate, endDate, projectArray, GalathiId, funderId)
	// 			summaryEnrolled += obj["enrolled"].(int)
	// 			obj["villages"] = newVillageCount(con, startDate, endDate, projectArray)
	// 			summaryVillages += obj["villages"].(int)
	// 			obj["startDate"] = ""
	// 			obj["endDate"] = ""
	// 			obj["select_type"] = "2"
	// 			//green motivators
	// 			obj["greenMotivators"] = greenMotivators(con, startDate, endDate, projectArray, funderId)
	// 			obj["vyapar"] = Vyapar(con, startDate, endDate, projectArray, funderId)
	// 			summaryGreen += obj["greenMotivators"].(int)
	// 			summaryVyapar += obj["vyapar"].(int)
	// 			data = append(data, obj)
	// 		}
	// 	}
	// 	// Convert the slice of Funders to JSON and write the response
	// 	jsonResponse, err := json.Marshal(funders)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	fmt.Fprint(w, string(jsonResponse))
	// })

	// Start the HTTP server

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8081", handler)
}
