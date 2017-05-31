package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
	"time"
	"path/filepath"
	"sbrkeygen/modules"
	"os/signal"
	"net"
)

var (
	amount map[int]int
	currency map[string]int
	fixed map[string]int
	dispdate map[int]int
)

const datapath = "..\\sbrkeygen-data"

func main() {
	// Баннер
	fmt.Println("SBERBANK TELEX KEY GENERATOR (C) 2017 ver.0.1")
	// Инициализация web-сервера
	keygen.WaitExit = false; //флаг для завершения работы
	var web keygen.WebCtl
	keygen.GlobalConfig.SetManagerSrv("127.0.0.1", 4040)
	fmt.Println("Web control configured: " + "http://" + keygen.GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(keygen.GlobalConfig.ManagerSrvPort())))
	web.SetHost(net.ParseIP(keygen.GlobalConfig.ManagerSrvAddr()))
	web.SetPort(keygen.GlobalConfig.ManagerSrvPort())
	/* Запускаем сервер обслуживания WebCtl */
	err := web.StartServe()
	if err != nil {
		log.Println("HTTP сервер: Ошибка. ", err)
		os.Exit(1)
	}

	/* Перехват CTRL+C для завершения приложения */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("\nReceived %v, shutdown procedure initiated.\n\n", sig)
			keygen.WaitExit = true
		}
	}()

	amount, err = ReadAmount(filepath.Join(datapath, "amount.txt"))
	if err != nil {
		log.Fatal("Ошибка чтения файла кодов сумм (amount.txt):", err)
	}
	currency, err = ReadCurrency(filepath.Join(datapath, "currency.txt"))
	if err != nil {
		log.Fatal("Ошибка чтения файла кодов валют (currency.txt):", err)
	}
	fixed, err = ReadFixed(filepath.Join(datapath, "fixed.txt"))
	if err != nil {
		log.Fatal("Ошибка чтения файла fixed.txt:", err)
	}

	dispdate, err = ReadDispDate(filepath.Join(datapath, "calendar.txt"))
	if err != nil {
		log.Fatal("Ошибка чтения файла calendar.txt:", err)
	}

	fmt.Println("res:", dispdate[GetDate(2)], GetDate(3))
	fmt.Println(CalcAmount(1000))

//	reader := bufio.NewReader(os.Stdin)
//	reader.ReadString('\n')

	ticker := time.NewTicker(time.Second * 2)

	for range ticker.C {
		// Запускаем обработчик каждую минуту
		if !keygen.WaitExit {
			continue
		}
		break
	}
	ticker.Stop()
}


// Выполняет разложение числа на разряды в массив
func SplitByAmount(num int) (splitted map[int]int) {
	// Amount: 1 10 100 1000 10000 100000 1000000 10000000 100000000 1000000000 10000000000
	var amtlist = make(map[int]int)

	amtlist [0] = 0
	amtlist [10000000000] = 0

	for i := 1; i <= 10000000000; i *= 10 {
		amtlist[i] = 0
	}

	//fmt.Println(amtlist)

	for i := 10000000000; i != 0 && num != 0; i /= 10 {
		//fmt.Println("NUM:", num, " ", i, " ", num / i)
		if num / i == 0 {
			continue
		}
		amtlist[i] = num / i
		num = num - i * (num / i)
	}

	/*	for i := 1; i <= 10000000000; i *= 10 {
			if amtlist[i] == 0 {
				continue
			}
			//fmt.Println(i, " ", amtlist[i])
		}*/

	return amtlist
}

// Считываем AMOUNT значения из файла
func ReadAmount(filename string) (amount map[int]int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return amount, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	amount = make(map[int]int)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, "\t")
		k, err := strconv.Atoi(str[0])
		if err != nil {
			return amount, err
		}
		v, err := strconv.Atoi(str[1])
		if err != nil {
			return amount, err
		}
		amount[k] = v
	}
	return amount, err
}

// Расчет суммы AMOUNT с учетом таблицы разбивки по раздрядам (amount.txt)
func CalcAmount(sum int) int {
	// Amount: 1 10 100 1000 10000 100000 1000000 10000000 100000000 1000000000 10000000000
	splitted := SplitByAmount(sum)
	fmt.Println(splitted)
	if sum >= 100000000000 {
		return amount[100000000000]
	}
	if sum == 0 {
		return amount[0]
	}

	sum = 0
	for k, v := range splitted {
		if v == 0 {
			continue
		}
		if k == 1 && v == 1 {
			v -= 1
		}
		sum += amount[k + v]
		fmt.Println("res:", amount[k + v], k + v)
	}

	return sum
}


// Считываем CURRENCY значения из файла
func ReadCurrency(filename string) (currency map[string]int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return currency, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	currency = make(map[string]int)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, "\t")
		v, err := strconv.Atoi(str[1])
		if err != nil {
			return currency, err
		}
		currency[str[0]] = v
	}
	return currency, err
}

// Считываем FIXED NUMBER значения из файла
func ReadFixed(filename string) (fixed map[string]int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return fixed, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	fixed = make(map[string]int)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, "\t")
		v, err := strconv.Atoi(str[1])
		if err != nil {
			return fixed, err
		}
		fixed[str[0]] = v
	}
	return fixed, err
}

// Возврат кода по валюте
func CalcCurrency(cur string) int {
	return currency[cur]
}

// Считываем calendar.txt значения из файла
func ReadDispDate(filename string) (dispdate map[int]int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return dispdate, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	dispdate = make(map[int]int)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, "\t")
		k1, err := strconv.Atoi(str[0])
		if err != nil {
			return dispdate, err
		}
		k2, err := strconv.Atoi(str[1])
		if err != nil {
			return dispdate, err
		}
		v, err := strconv.Atoi(str[2])
		if err != nil {
			return dispdate, err
		}

		dispdate[k1 * 100 + k2] = v
	}
	return dispdate, err
}


// Возвращает код для текущей даты календаря
// shift = двиг от текущей даты в днях
func GetDate(shift int) int {
	date := time.Now()
	date = date.Add(time.Duration(shift) * 24 * time.Hour)
	//res := strconv.Itoa(date.Year()) //ГГГГ
	res := fmt.Sprintf("%2d", int(date.Month())) //ММ
	res += fmt.Sprintf("%02d", date.Day()) //ДД
	i, err := strconv.Atoi(strings.TrimLeft(res, " "))
	if err != nil {

		return -1
	}
	return i
}
