package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type fPos struct {
	id      int //144:受診者ID
	jname   int //147:受診者名
	jkana   int //145:ﾌﾘｶﾞﾅ
	sei     int //148:性別
	seinen  int //149:生年月日
	sno     int //153:社員No
	kname   int //21:企業名
	kcd     int //20:企業cd
	scd     int //26:所蔵cd１
	sname   int //27:所属名１
	gantei1 int //●眼底片眼
	gantei2 int //●眼底両眼
}

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	// ログファイル準備
	logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	failOnError(err)
	defer logfile.Close()

	log.SetOutput(logfile)
	log.Print("Start\r\n")

	//ドロップされたファイルの数で処理を分ける
	filesu := flag.NArg()
	if filesu == 1 {
		// ファイルを読み込んで二次元配列に入れる
		records := readfile(flag.Arg(0))

		// 眼底ノートPC用のファイル処理
		precs := processRecord(records)

		// ファイルを協会けんぽ資格確認用のCSVに出力
		saveCsv(precs)

	} else {
		log.Print("ドロップファイルエラー。処理を終了します。")
		os.Exit(1)
	}

	log.Print("Finesh !\r\n")

}

func readfile(filename string) [][]string {
	// 入力ファイル準備
	infile, err := os.Open(filename)
	failOnError(err)
	defer infile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'

	//CSVファイルを２次元配列に展開
	readrecords := make([][]string, 0)
	record, err := reader.Read() // 1行読み出す
	if err == io.EOF {
		return readrecords
	} else {
		failOnError(err)
	}

	colMax := len(record) - 1
	//readrecords = append(readrecords, record[:colMax])
	readrecords = append(readrecords, record[:colMax])

	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			// log.Print(record)
			// log.Print(len(record))
			failOnError(err)
		}

		readrecords = append(readrecords, record[:colMax])

	}

	return readrecords
}

func processRecord(precs [][]string) [][]string {

	// 眼底ノートPC用の受診者項目を抽出する
	filePos := fPos{id: -1, jname: -1, jkana: -1, sei: -1, seinen: -1, sno: -1, kname: -1, kcd: -1, scd: -1, sname: -1, gantei1: -1, gantei2: -1}

	hedRow := precs[0]
	for pos, colName := range hedRow {
		switch colName {
		case "受診者ID":
			filePos.id = pos
		case "受診者名":
			filePos.jname = pos
		case "ﾌﾘｶﾞﾅ":
			filePos.jkana = pos
		case "性別":
			filePos.sei = pos
		case "生年月日":
			filePos.seinen = pos
		case "社員No":
			filePos.sno = pos
		case "企業名":
			filePos.kname = pos
		case "企業cd":
			filePos.kcd = pos
		case "所属cd１":
			filePos.scd = pos
		case "所属名１":
			filePos.sname = pos
		case "●眼底片眼":
			filePos.gantei1 = pos
		case "●眼底両眼":
			filePos.gantei2 = pos
		}
	}

	// 眼底の対象者および項目のレコードを作成する
	grecs := make([][]string, 0)
	for i, v := range precs {
		if i == 0 || gCheck(v[filePos.gantei1], v[filePos.gantei2]) {
			k1 := v[filePos.id]
			k2 := v[filePos.jname]
			k3 := v[filePos.jkana]
			k4 := v[filePos.sei]
			k5 := v[filePos.seinen]
			k6 := v[filePos.sno]
			k7 := v[filePos.kname]
			k8 := v[filePos.kcd]
			k9 := v[filePos.scd]
			k10 := v[filePos.sname]
			//log.Print(grec)
			grecs = append(grecs, []string{k1, k2, k3, k4, k5, k6, k7, k8, k9, k10})
		}
	}

	return grecs

}

func saveCsv(recs [][]string) {

	// 出力ファイル準備
	outfile, err := os.Create("./眼底ノートPC用受診者名簿" + time.Now().Format("0102") + ".csv")
	failOnError(err)
	defer outfile.Close()

	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.Comma = ','
	writer.UseCRLF = true

	for _, recRow := range recs {
		writer.Write(recRow)
	}

	writer.Flush()

}

func gCheck(g1 string, g2 string) bool {
	// 眼底対象者かチェック
	var check bool
	if g1 == "●" || g2 == "●" {
		check = true
	} else {
		check = false
	}

	return check
}
