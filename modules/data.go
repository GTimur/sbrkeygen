package keygen

import (
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"time"
	"os"
	"path/filepath"
	"log"
	"errors"
	"os/user"
)

type Telex struct {
	Sum    int64
	Cur    string
	Msg    string
	Date   string
	SeqCnt int
	Key    string
}

var (
	Amount map[int64]int
	Currency map[string]int
	Fixed map[string]int
	Dispdate map[int64]int
	Seqfrom map[int64]int
	Seqto map[int64]int
	Msg Telex

	CalcLog string          // лог вычисления ключа
	SeqCnt int                 // номер сообщения в СБЕР в течение года
	copytelexpath string    // файли логов и
	copylogpath string      // телекса копируются в папку пользователя
)

const (
	datapath = "..\\sbrkeygen-data\\system"
	seqfile = "seqcount.dat"
	logpath = "..\\sbrkeygen-data\\log"
	telexpath = "..\\sbrkeygen-data\\telex"
)

func (t *Telex) SetParams(sum int64, cur string, msg string, date string, seqcnt int, key string) error {
	if sum <= 0 {
		return errors.New("Не указано значение суммы")
	}
	t.Sum = sum
	if len(cur) == 0 {
		return errors.New("Не указана валюта сделки")
	}
	t.Cur = cur
	if len(msg) == 0 {
		return errors.New("Не заполнено сопровождающее сообщение к сделке")
	}
	t.Msg = msg
	if len(date) == 0 {
		return errors.New("Не указана дата сделки")
	}
	t.Date = date
	if seqcnt <= 0 {
		return errors.New("Не указан номер последовательности для сделки")
	}
	t.SeqCnt = seqcnt
	if len(key) == 0 {
		return errors.New("Ключ TELEX не может быть пустым")
	}
	t.Key = key

	return nil
}

// Выполняет разложение числа на разряды в массив
func SplitByAmount(num int64) (splitted map[int64]int64) {
	// Amount: 1 10 100 1000 10000 100000 1000000 10000000 100000000 1000000000 10000000000
	var amtlist = make(map[int64]int64)

	amtlist [0] = 0
	amtlist [10000000000] = 0

	for i := int64(1); i <= 10000000000; i *= 10 {
		amtlist[i] = 0
	}

	//fmt.Println(amtlist)

	for i := int64(10000000000); i != 0 && num != 0; i /= 10 {
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
func ReadAmount(filename string) (amount map[int64]int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return amount, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	amount = make(map[int64]int)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, "\t")
		k, err := strconv.ParseInt(str[0], 10, 64)//strconv.Atoi(str[0])
		if err != nil {
			return amount, err
		}
		v, err := strconv.ParseInt(str[1], 10, 64)//strconv.Atoi(str[1])
		if err != nil {
			return amount, err
		}
		amount[int64(k)] = int(v)
	}
	return amount, err
}

// Расчет суммы AMOUNT с учетом таблицы разбивки по раздрядам (amount.txt)
func CalcAmount(sum int64) int {
	// Amount: 1 10 100 1000 10000 100000 1000000 10000000 100000000 1000000000 10000000000
	splitted := SplitByAmount(sum)
	//fmt.Println(splitted)
	log := ""
	if sum >= 100000000000 {
		res := Amount[100000000000]
		CalcLog += "\nAMOUNT: " + strconv.Itoa(res)
		CalcLog += "\n" + strconv.FormatInt(sum, 10) + ">=100000000000 ===> " + strconv.Itoa(Amount[100000000000])
		return res
	}
	if sum == 0 {
		res := Amount[0]
		CalcLog += "\nAMOUNT: " + strconv.Itoa(res)
		CalcLog += "\n" + strconv.FormatInt(sum, 10) + "=0 ===> " + strconv.Itoa(Amount[0])
		return res
	}

	res := 0
	for k, v := range splitted {
		if v == 0 {
			continue
		}
		log += "\n" + strconv.FormatInt(sum, 10) + "=" + strconv.FormatInt(v, 10) + "*" + strconv.FormatInt(k, 10) + " ===> " + strconv.Itoa(Amount[k + v])
		if k == 1 && v == 1 {
			v -= 1
		}
		res += Amount[k + v]
		//fmt.Println("res:", Amount[k + v], k + v)
	}

	CalcLog += "\nAMOUNT: " + strconv.Itoa(res) + log
	return res
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
	return Currency[cur]
}

// Возврат кода для Seq
func CalcSeq(seq int, isfrom bool) int {
	if seq <= 0 || seq > 128 {
		seq = 1
		SeqCnt = 1
		UpdateSeqCnt(1)
	}
	if !isfrom {
		CalcLog += "\nSEQUENCE (TO SBER): " + strconv.Itoa(seq) + " ===> " + strconv.Itoa(Seqto[int64(seq)])
		return Seqto[int64(seq)]
	}
	CalcLog += "\nSEQUENCE (FROM SBER): " + strconv.Itoa(seq) + " ===> " + strconv.Itoa(Seqfrom[int64(seq)])
	return Seqfrom[int64(seq)]

}


// Считываем calendar.txt значения из файла
func ReadDispDate(filename string) (dispdate map[int64]int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return dispdate, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	dispdate = make(map[int64]int)
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

		dispdate[int64(k1 * 100 + k2)] = v
	}
	return dispdate, err
}

// Считываем SEQUENCE number либо из from.txt либо из to.txt
// заивсит от флага isfrom
func ReadFromTo(from string, to string, isfrom bool) (res map[int64]int, err error) {
	filename := from
	if !isfrom {
		filename = to
	}
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return res, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	res = make(map[int64]int)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, "\t")
		k, err := strconv.Atoi(str[0])
		if err != nil {
			return res, err
		}
		v, err := strconv.Atoi(str[1])
		if err != nil {
			return res, err
		}
		res[int64(k)] = v
	}
	return res, err
}


// Считываем SEQ COUNT значения из файла
// счетчик отправленных сообщений в телекс
func ReadSeqCount(filename string) (counter int, err error) {
	// Open the file.
	f, err := os.Open(filename)
	if err != nil {
		return counter, err
	}
	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		v, err := strconv.Atoi(line)
		if err != nil {
			return counter, err
		}
		counter = v
	}
	return counter, err
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

// Задает номер счетчика сообщений в файле (перезапишет файл)
func UpdateSeqCnt(cnt int) (err error) {
	/* Создадим/перезапишем файл */
	file, err := os.Create(filepath.Join(datapath, seqfile))
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(strconv.Itoa(cnt))

	return err
}


// Сохраняет сообщение с логом расчета
func WriteCalcLog() (err error) {
	/* Создадим/перезапишем файл */
	prefix := fmt.Sprintf("%03d", SeqCnt)
	now := time.Now()
	datefix := strings.Replace(Msg.Date, "/", "", -1) + fmt.Sprintf("%02d%02d", now.Hour(), now.Second())
	file, err := os.Create(filepath.Join(logpath, prefix + "-" + datefix + "-calc.txt"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(CalcLog)
	if err != nil {
		return err
	}

	/**MAKE COPY**/
	filecopy, err := os.Create(filepath.Join(copylogpath, prefix + "-" + datefix + "-calc.txt"))
	if err != nil {
		return err
	}
	defer filecopy.Close()

	_, err = filecopy.WriteString(CalcLog)
	if err != nil {
		return err
	}
	/******/


	/* Если данныые записаны на диск - увеличиваем счетчик */

	return err
}

// Сохраняет сообщение TELEX
func WriteTelex() (err error) {
	/* Создадим/перезапишем файл */
	prefix := fmt.Sprintf("%03d", SeqCnt)
	now := time.Now()
	datefix := strings.Replace(Msg.Date, "/", "", -1) + fmt.Sprintf("%02d%02d", now.Hour(), now.Second())

	file, err := os.Create(filepath.Join(telexpath, prefix + "-" + datefix + "-telex.txt"))
	if err != nil {
		return err
	}
	defer file.Close()

	TelexMessage := "UMK BANK\n\n" +
		"DATE\n" + Msg.Date + "\n\n" +
		"FREE FORMAT MESSAGE\n\n" +
		Msg.Msg + "\n\n" +
		"\tTELEX KEY IS " + Msg.Key + "\n\n" +
		"\tBEST REGARDS,\n" +
		"\tAO UMK BANK\n" +
		"\t(861)2100553\n\n" +
		"END OF MESSAGE\n\n"

	_, err = file.WriteString(TelexMessage)
	if err != nil {
		return err
	}

	/**MAKE FILE COPY**/
	filecopy, err := os.Create(filepath.Join(copytelexpath, prefix + "-" + datefix + "-telex.txt"))
	if err != nil {
		return err
	}
	defer filecopy.Close()

	_, err = filecopy.WriteString(TelexMessage)
	if err != nil {
		return err
	}
	/*****/

	err = WriteCalcLog()
	if err != nil {
		return err
	}

	/* Если данныые записаны на диск - увеличиваем счетчик */
	UpdateSeqCnt(SeqCnt)
	SeqCnt++

	return err
}

func InitData() (err error) {
	Amount, err = ReadAmount(filepath.Join(datapath, "amount.txt"))
	if err != nil {
		log.Println("Ошибка чтения файла кодов сумм (amount.txt):", err)
		return err
	}
	Currency, err = ReadCurrency(filepath.Join(datapath, "currency.txt"))
	if err != nil {
		log.Println("Ошибка чтения файла кодов валют (currency.txt):", err)
		return err
	}
	Fixed, err = ReadFixed(filepath.Join(datapath, "fixed.txt"))
	if err != nil {
		log.Println("Ошибка чтения файла fixed.txt:", err)
		return err
	}
	Dispdate, err = ReadDispDate(filepath.Join(datapath, "calendar.txt"))
	if err != nil {
		log.Println("Ошибка чтения файла calendar.txt:", err)
		return err
	}
	Seqfrom, err = ReadFromTo(filepath.Join(datapath, "seqfrom.txt"), filepath.Join(datapath, "seqto.txt"), true)
	if err != nil {
		log.Println("Ошибка чтения файла seqfrom.txt:", err)
		return err
	}
	Seqto, err = ReadFromTo(filepath.Join(datapath, "seqfrom.txt"), filepath.Join(datapath, "seqto.txt"), false)
	if err != nil {
		log.Println("Ошибка чтения файла seqto.txt:", err)
		return err
	}
	SeqCnt, err = ReadSeqCount(filepath.Join(datapath, "seqcount.dat"))
	if err != nil {
		log.Println("Ошибка чтения файла seqcount.dat", err)
		return err
	}

	userhome, err := GetUserHome()
	if err != nil {
		log.Println("Ошибка чтения домашнего каталога пользователя:", err)
		return err
	}


	copylogpath = filepath.Join(userhome, "DOCUMENTS", "TELEXGEN", "LOG")
	copytelexpath = filepath.Join(userhome, "DOCUMENTS", "TELEXGEN", "TELEX")
	if _, err := os.Stat(copylogpath); os.IsNotExist(err) {
		os.MkdirAll(copylogpath,0777);
	}

	if _, err := os.Stat(copytelexpath); os.IsNotExist(err) {
		os.MkdirAll(copytelexpath,0777);
	}

	// Зададим номер следующего сообщения
	SeqCnt += 1;
	if SeqCnt <= 0 || SeqCnt > 128 {
		SeqCnt = 1
		UpdateSeqCnt(1)
	}
	/*err = UpdateSeqCnt(100)
	if err != nil {
		log.Println("Ошибка записи файла seqcount.dat", err)
		return err
	}*/

	return err
}

func GetUserHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, err
}

// Вычислет ключ на основе следующих параметров
// SUM = сумма
// CUR = строковый код валюты
// SEQ = номер по порядку сообщения
// calshift = сдвиг в днях относительно текущей даты (код по календарю)

func CalcKey(sum int64, cur string, seq int, isseqfrom bool, calshift int) (key int) {
	CalcLog = "";
	CalcLog += "CURRENCY: " + cur + " ===> " + strconv.Itoa(CalcCurrency(cur))
	key += CalcCurrency(cur)
	key += CalcAmount(sum) //содержит CalcLog
	CalcLog += "\nDATE OF DISPATCH: " + strconv.Itoa(GetDate(0)) + " ===> " + strconv.Itoa(Dispdate[int64(GetDate(0))])
	key += Dispdate[int64(GetDate(0))]
	key += CalcSeq(seq, isseqfrom)
	CalcLog += "\nFIXED NUMBER: " + strconv.Itoa(Fixed["FIXED"])
	key += Fixed["FIXED"]
	CalcLog += "\n";
	return key
}
