package main

import (
	a "buzzstaff/dashboard"
)

func main() {
	a.HandleFunc()
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
