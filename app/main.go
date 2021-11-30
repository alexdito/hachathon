package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type Row struct {
	UserId int `json:"user_id"`
	Timestamp string `json:"timestamp"`
	Category string `json:"category"`
	Card string `json:"card"`
	Amount int `json:"amount"`
}

type Category struct {
	Name string
	Count int
	Sum int
}

type Reports struct {
	Report []Report
}

type Report struct {
	UserId int
	Sum int
	CategoriesSum map[string]int
	CategoriesAmount map[string]int
}

func main() {
	start := time.Now()

	jsonFile, err := os.Open("transactions.json")
	defer jsonFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	rows := []Row{}

	json.Unmarshal(byteValue, &rows)

	fmt.Println(fmt.Sprintf("Parse Json: %s", time.Since(start)))

	userIds := make(map[int]int)

	categories := make(map[int][]Category)

	for i := 0; i < len(rows); i++ {
		userIds[rows[i].UserId] = rows[i].UserId

		categories[rows[i].UserId] = append(categories[rows[i].UserId], Category{Name: rows[i].Category, Sum: rows[i].Amount})
	}


	reports := generateReport(userIds, categories)

	generateJson(reports)
	generateCsv(reports)

	fmt.Println(fmt.Sprintf("General time: %s", time.Since(start)))
}

func GetSumCategories(categories map[string]int) int {
	var sum int

	for _, v := range categories {
		sum += v
	}

	return sum
}

func generateJson(reports Reports)  {
	start := time.Now()

	result, err := json.Marshal(reports)

	if err != nil {
		log.Println(err)
	}

	reportFile, err := os.Create("report-transactions.json")

	if err != nil {
		fmt.Println(err)
	}

	reportFile.Write(result)

	fmt.Println("Generated report-transactions.json")

	reportFile.Close()

	fmt.Println(fmt.Sprintf("Generate Json: %s", time.Since(start)))
}

func generateCsv(reports Reports)  {
	start := time.Now()

	fileCsv, err := os.Create("report-transactions.csv")
	defer fileCsv.Close()

	recordRows := [][]string{
		{"User id",	"General sum", "Category",	"Number of category appearances", "Category sum"},
	}


	for _, report := range reports.Report {
		setHead := true
		for category, _ := range report.CategoriesSum {
			 if setHead {
				 recordRows = append(
					 recordRows, []string{
						 strconv.Itoa(report.UserId),
						 strconv.Itoa(report.Sum),
						 category,
						 strconv.Itoa(report.CategoriesAmount[category]),
						 strconv.Itoa(report.CategoriesSum[category]),
					 },
				 )
				 setHead = false
			 } else {
				 recordRows = append(
					 recordRows, []string{
						 "",
						 "",
						 category,
						 strconv.Itoa(report.CategoriesAmount[category]),
						 strconv.Itoa(report.CategoriesSum[category]),
					 },
				 )
			 }
		}
	}


	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	writeCsv := csv.NewWriter(fileCsv)
	defer writeCsv.Flush()

	err = writeCsv.WriteAll(recordRows)
	fmt.Println("Generated report-transactions.csv")

	fmt.Println(fmt.Sprintf("Generate Csv: %s", time.Since(start)))
}

func generateReport(userIds map[int]int, categories map[int][]Category) Reports  {
	reports := Reports{}
	report := Report{}

	for i := 1; i <= len(userIds); i++ {
		report.UserId = userIds[i]
		report.CategoriesSum = make(map[string]int)
		report.CategoriesAmount = make(map[string]int)

		for j := 0; j < len(categories[userIds[i]]); j++ {

			if _, isExist := report.CategoriesSum[categories[i][j].Name]; !isExist {
				report.CategoriesSum[categories[i][j].Name] = 1
				report.CategoriesAmount[categories[i][j].Name] = categories[i][j].Sum
			} else {
				report.CategoriesSum[categories[i][j].Name] += 1
				report.CategoriesAmount[categories[i][j].Name] += categories[i][j].Sum
			}
		}

		report.Sum = GetSumCategories(report.CategoriesAmount)
		reports.Report = append(reports.Report, report)
	}

	return reports
}