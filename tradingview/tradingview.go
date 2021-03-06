package tradingview

import (
	"CrawlerTable/lembarsaham"
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"regexp"
	"time"
)

//tradingvi

func Crawl(filename string, delay time.Duration) {
	var row []string

	var rows [][]string

	doc := lembarsaham.GetBody("https://id.tradingview.com/markets/stocks-indonesia/sectorandindustry-sector/")
	doc.Find("tbody a").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		tablebody := lembarsaham.GetBody("https://id.tradingview.com" + link)
		time.Sleep(delay * time.Second)
		fmt.Println("Delay " + delay.String() + " detik")
		tablebody.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
			tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
				rowhtml.Find("td").Each(func(indexth int, tablecell *goquery.Selection) {
					space := regexp.MustCompile(`\s+`)
					s := space.ReplaceAllString(tablecell.Text(), " ")
					row = append(row, s)
				})
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

	headings := []string{"KODE NAMA PT", "Terakhir", "%Perubahan", "Perrubahan", "Penialain", "Volume", "Cap Pasar", "P/E", "EPS(TTM)", "Pegawai", "Jenis Industri"}
	err = writer.Write(headings)
	lembarsaham.CheckError("Cannot write to file", err)

	for _, data := range rows {
		if len(data) > 0 {
			err := writer.Write(data)
			lembarsaham.CheckError("Cannot write to file", err)
		}
	}

}
