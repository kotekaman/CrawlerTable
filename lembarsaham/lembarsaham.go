package lembarsaham

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
func GetBody(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	for res.StatusCode != 200 {
		time.Sleep(20 * time.Second)
		res, err = http.Get(url)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}
func Crawl(filename string, delay time.Duration) {
	var headings, row []string
	var rows [][]string
	var file *os.File
	var data []string

	doc := GetBody("https://lembarsaham.com/daftar-emiten/9-sektor-bei")
	doc.Find("h4 a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		tablebody := GetBody("https://lembarsaham.com" + link)
		time.Sleep(delay * time.Second)
		fmt.Println("Delay " + delay.String() + " detik")
		tablebody.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
			tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
				rowhtml.Find("th").Each(func(indexth int, tableheading *goquery.Selection) {
					if len(headings) <= 5 {
						headings = append(headings, tableheading.Text())
					}
				})
				rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
					space := regexp.MustCompile(`\s+`)
					s := space.ReplaceAllString(tablecell.Text(), " ")
					row = append(row, s)
				})
				if len(row) > 5 {
					data = strings.Split(row[5], " ")
					row[5] = data[1]
					row = append(row, data[2]+" "+data[3]+" "+data[4])
				}
				rows = append(rows, row)
				row = nil
			})
		})
	})

	file, err := os.Open(filename)
	//file, err := os.OpenFile(filename,os.O_APPEND|os.O_RDWR,os.ModeAppend)
	if err != nil {
		file, err = os.Create(filename)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_RDWR, os.ModeAppend)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headings = append(headings, "waktu_update")
	err = writer.Write(headings)
	CheckError("Cannot write to file", err)

	for _, data := range rows {
		if len(data) > 0 {
			err := writer.Write(data)
			CheckError("Cannot write to file", err)
		}
	}
}
