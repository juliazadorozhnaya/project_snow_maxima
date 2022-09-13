package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TstRef() (x []int, y []int, z []int) {

	var xRef = []int{
		0, 1, 2, 3, 4, 5, 6,
		0, 2, 3, 4, 5, 6,
		0, 1, 2, 3, 4, 5, 6,
		0, 1, 3, 4, 5, 6,
		0, 1, 2, 3, 4, 5, 6,
	}
	var yRef = []int{
		0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4, 4,
	}
	var zRef = []int{
		2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2,
	}
	return xRef, yRef, zRef
}

func TstFull() (x []int, y []int, z []int) {

	var xFull = []int{
		0, 1, 2, 3, 4, 5, 6,
		0, 1, 2, 3, 4, 5, 6,
		0, 1, 2, 3, 6,
		0, 2, 3, 4, 5, 6,
		0, 1, 2, 3, 4, 5, 6,
	}
	var yFull = []int{
		0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2,
		3, 3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4, 4,
	}
	var zFull = []int{
		2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2,
		3, 3, 3, 3, 3,
		2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2,
	}
	return xFull, yFull, zFull
}

func getArg0[T any](a T, b error) T {
	if b != nil {
		fmt.Println(b)
	}
	return a
}

func GetScanData(iip int) (x []int, y []int, z []int) {
	fileName := strconv.Itoa(iip) + ".bin"
	if iip == -1 {
		x, y, z = TstFull()
	} else if iip == -2 {
		x, y, z = TstRef()
	} else {
		var body []byte
		fmt.Println("loading .bin file from PAC...")
		ip_adr := xxx
		usr_nm := xxx
		passwd := xxx
		database_nm := xxx
		db_port := xxx
		db_conn, err := pgxpool.Connect(context.Background(),
			"postgres://"+usr_nm+":"+passwd+"@"+ip_adr+":"+db_port+"/"+database_nm,
		)
		if err != nil {
			fmt.Println(err)
		} else {
			var binFilePath string
			err = db_conn.QueryRow(context.Background(), "SELECT cargo_scan_bin FROM iips WHERE id = $1", iip).Scan(&binFilePath)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Got bin file path:", binFilePath)
			}
			binFilePath = "https://pac.msk.vmet.ro/media/2" + binFilePath
			fmt.Println("Loading bin file from:", binFilePath)

			resp := getArg0(http.Get(binFilePath))
			defer resp.Body.Close()

			body = getArg0(io.ReadAll(resp.Body))
			fmt.Println("Bin file size =", len(body), "bytes")

			file := getArg0(os.Create(fileName))
			defer file.Close()
			file.Write(body)
			fmt.Println("File \"" + fileName + "\" saved.")
		}
		defer db_conn.Close()
		x, y, z = parseBinFile(body)
	}
	return x, y, z
}

func loadBinFileFromHDD(fileName string) []byte {
	file := getArg0(os.Open(fileName))
	defer file.Close()

	fsize := getArg0(file.Stat()).Size()
	fmt.Println("File\""+fileName+"\" size:", fsize)
	data := make([]byte, fsize)

	if fsize != int64(getArg0(file.Read(data))) { // если недочитали
		fmt.Println("Can't load full file.")
		os.Exit(1)
	}
	return data
}

func parseBinFile(data []byte) (x, y, z []int) {
	var p int

	MeasurementID := binary.LittleEndian.Uint32(data[p:])
	p += 4
	fmt.Println("MeasurementID", MeasurementID)

	IncomingMeasurementID := binary.LittleEndian.Uint32(data[p:])
	p += 4
	fmt.Println("IncomingMeasurementID", IncomingMeasurementID)

	CustomerMeasurementID := binary.LittleEndian.Uint32(data[p:])
	p += 4
	fmt.Println("CustomerMeasurementID", CustomerMeasurementID)

	IncomingUnitID := binary.LittleEndian.Uint32(data[p:])
	p += 4
	fmt.Println("IncomingUnitID", IncomingUnitID)

	TruckIDLen := data[p]
	p++
	fmt.Println("TruckIDLen", TruckIDLen)

	TruckID := string(data[p : p+int(TruckIDLen)])
	p += int(TruckIDLen)
	fmt.Println("TruckID", TruckID)

	// всегда ноль. Кажется признак "конца строки с TruckID"
	p += 4

	Delimiter := binary.LittleEndian.Uint32(data[p:])
	p += 4
	fmt.Println("Delimiter", Delimiter)

	//data:time? по описанию должны быть поля: год(с 2000 вроде), месяц, день, час, минута, секунда.
	// как раз 6 штук по 4 байта
	p += 24

	CountOfBlocks := binary.LittleEndian.Uint32(data[p:])
	p += 4
	fmt.Println("CountOfBlocks", CountOfBlocks)

	for i := uint32(0); i < CountOfBlocks; i++ {
		tmpX, tmpY, tmpZ, tmpP := parseBlock(data[p:], int(Delimiter))
		x = append(x, tmpX...)
		y = append(y, tmpY...)
		z = append(z, tmpZ...)
		p += tmpP
		fmt.Println("Block №", i, " parsed")
	}

	if len(x) < 3 {
		x = append(x, 0, 1, 3)
		y = append(y, 0, 1, 1)
		z = append(z, 0, 1, 1)
	}
	return x, y, z
}

func parseBlock(data []byte, delimiter int) (x, y, z []int, p int) {
	p += 28
	countOfLines := binary.LittleEndian.Uint32(data[p:])
	//fmt.Println("Block consists of", countOfLines, "lines")
	p += 4
	var tmpX, tmpY, tmpZ int
	for i := uint32(0); i < countOfLines; i++ {
		p += 32
		tmpX = int(int32(binary.LittleEndian.Uint32(data[p:])))
		p += 4
		tmpY = int(int32(binary.LittleEndian.Uint32(data[p:])))
		p += 4
		tmpZ = int(int32(binary.LittleEndian.Uint32(data[p:])))
		p += 4
		x = append(x, tmpX)
		y = append(y, tmpY)
		z = append(z, tmpZ)
		p += 8
	}
	p += 132
	return x, y, z, p
}
