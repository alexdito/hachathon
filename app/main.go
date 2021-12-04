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
	UserId    int    `json:"user_id"`
	Timestamp string `json:"timestamp"`
	Category  string `json:"category"`
	Card      string `json:"card"`
	Amount    int    `json:"amount"`
}

type Reports struct {
	Report map[int]Report
}

type Report struct {
	UserId           int
	Sum              int
	CategoriesSum    map[string]int
	CategoriesAmount map[string]int
}

func (reports *Reports) CalculateSumReports()  {
	for _, report := range reports.Report {
		if thisProduct, ok := reports.Report[report.UserId]; ok {
			thisProduct.Sum = GetSumCategories(report.CategoriesSum)
			reports.Report[report.UserId] = thisProduct
		}
	}
}

func (report *Report) getSum() int{
	for _, v := range report.CategoriesSum {
		report.Sum += v
	}

	return report.Sum
}

func main() {
	start := time.Now()

	jsonFile, err := os.Open("transactions.json")
	defer jsonFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var rows []Row

	json.Unmarshal(byteValue, &rows)

	fmt.Println(fmt.Sprintf("Parse Json: %s", time.Since(start)))

	reports := Reports{}
	reports.Report = make(map[int]Report)

	for i := 0; i < len(rows); i++ {

		if _, isExist := reports.Report[rows[i].UserId]; !isExist {
			reports.Report[rows[i].UserId] = Report{UserId: rows[i].UserId, CategoriesSum: make(map[string]int), CategoriesAmount: make(map[string]int)}
		}

		if _, isExist := reports.Report[rows[i].UserId].CategoriesAmount[rows[i].Category]; !isExist {
			reports.Report[rows[i].UserId].CategoriesAmount[rows[i].Category] = 1
		} else {
			reports.Report[rows[i].UserId].CategoriesAmount[rows[i].Category] += 1
		}

		reports.Report[rows[i].UserId].CategoriesSum[rows[i].Category] += rows[i].Amount
	}

	reports.CalculateSumReports()

	go generateJson(reports)
	go generateCsv(reports)

	fmt.Println(fmt.Sprintf("General time: %s", time.Since(start)))

	time.Sleep(time.Millisecond)
}

func GetSumCategories(categories map[string]int) int {
	var sum int

	for _, v := range categories {
		sum += v
	}

	return sum
}

func generateJson(reports Reports) {
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

func generateCsv(reports Reports) {
	start := time.Now()

	fileCsv, err := os.Create("report-transactions.csv")
	defer fileCsv.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	writeCsv := csv.NewWriter(fileCsv)
	defer writeCsv.Flush()

	err = writeCsv.Write([]string{"User id", "General sum", "Category", "Number of category appearances", "Category sum"})

	for _, report := range reports.Report {
		for category, _ := range report.CategoriesSum {
			err = writeCsv.Write([]string{
				strconv.Itoa(report.UserId),
				strconv.Itoa(report.Sum),
				category,
				strconv.Itoa(report.CategoriesAmount[category]),
				strconv.Itoa(report.CategoriesSum[category]),
			})
		}
	}

	fmt.Println("Generated report-transactions.csv")

	fmt.Println(fmt.Sprintf("Generate Csv: %s", time.Since(start)))
}
