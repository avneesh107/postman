package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func main() {
	var filePath string
	fmt.Printf("enter the name of thie file:\n")
	fmt.Scanln(&filePath)
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer f.Close()
	count := 0.0
	sheet := "CSF111_202425_01_GradeBook"
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}
	fields := map[int]string{1: "S.NO", 2: "Class No.", 3: "Employee ID", 4: "campus ID", 5: "Quiz", 6: "Mid-sem", 7: "Lab Tests", 8: "Weekly Labs", 9: "Precompre", 10: "Compre", 11: "Total"}
	var avg = map[int]float64{4: 0, 5: 0, 6: 0, 7: 0, 9: 0, 10: 0}
	var courseavg = make(map[string]float64)
	var coursecount = make(map[string]float64)
	maxindex := 10
	topsmark := [7][3]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	topsrow := [7][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	fmt.Println("Checking total marks")

	for i, row := range rows[1:] {
		count += 1
		branch := row[3][4:6]
		if val, err := strconv.ParseFloat(row[10], 64); err == nil {
			courseavg[branch] += val
			coursecount[branch]++

		} else {
			fmt.Printf("Invalid format")
		}

		total := 0.0
		indices := []int{4, 5, 6, 7, 9}

		for j := 4; j <= 10; j++ {
			if val, err := strconv.ParseFloat(row[j], 64); err == nil {
				avg[j] += val
				if val > topsmark[j-4][2] {
					if val > topsmark[j-4][1] {
						if val > topsmark[j-4][0] {
							topsmark[j-4][2] = topsmark[j-4][1]
							topsrow[j-4][2] = topsrow[j-4][1]
							topsmark[j-4][1] = topsmark[j-4][0]
							topsrow[j-4][1] = topsrow[j-4][0]
							topsmark[j-4][0] = val
							topsrow[j-4][0] = i

						} else {
							topsmark[j-4][2] = topsmark[j-4][1]
							topsrow[j-4][2] = topsrow[j-4][1]
							topsmark[j-4][1] = val
							topsrow[j-4][1] = i
						}
					} else {
						topsmark[j-4][2] = val
						topsrow[j-4][2] = i
					}
				}

			} else {
				fmt.Printf("Invalid format")
			}

		}

		for _, k := range indices {
			val, err := strconv.ParseFloat(row[k], 64)
			if err == nil {
				total += val
			} else {
				fmt.Printf("Invalid format")
			}
		}

		t2, err := strconv.ParseFloat(row[maxindex], 64)
		if err != nil {
			fmt.Printf("Invalid format")
		}

		if total-t2 >= 0.001 || total-t2 <= -0.001 {
			fmt.Printf("Row %d: Total is incorrect! calculated Total: %.2f, Present total: %.2f\n", i+2, total, t2)
		}
	}

	fmt.Println("\nColumn Averages:")
	for i := 4; i <= 10; i++ {
		fmt.Printf("Avg Of %s: %.2f\n", fields[i+1], avg[i]/count)

	}
	fmt.Println("\nBranch wise avgs:")
	for i, total := range courseavg {
		fmt.Printf("avg of %s : %.2f\n", i, total/coursecount[i])

	}
	fmt.Println("\nTop 3 scores")
	for i := 0; i <= 6; i++ {
		fmt.Printf("\nTop performer of: %s are\n", fields[i+5])
		for j := 0; j <= 2; j++ {
			fmt.Printf("%d.employee id = %s, marks = %.2f\n", j+1, rows[topsrow[i][j]][2], topsmark[i][j])
		}
	}
}
