package keygen

import (
	"log"
	"net/http"
	//	"github.com/braintree/manners"
	"net"
	"fmt"
	"strconv"
	"path"
	"html/template"
	"time"
	"encoding/json"
	"context"
)

type WebCtl struct {
	host     net.IP
	port     uint16
	islisten bool
}

type Config struct {
	managerSrv managerSrv
}

//Представляет адрес сервера управления программой и порт
type managerSrv struct {
	Addr string
	Port uint16
}

type Page struct {
	Title   string
	Body    template.HTML
	LnkHome string
	DateNow template.HTML
	SeqCnt  int
}



var (
	GlobalConfig Config = Config{} //Глобальная переменная для хранения настроек
	home_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "index.gtpl")))
	WaitExit bool
	Quit = make(chan int,1 )  //канал для завершения сервера HTTP

)

/*Сервер*/
//Запускает goroutine ListenAndServe
//Может изменять accbook - справочник подписантов
func (w *WebCtl) StartServe() (err error) {
	//signal.Notify(Quit, os.Interrupt)
	srv := &http.Server{Addr : w.connString(), Handler: http.DefaultServeMux}

	// для отдачи сервером статичных файлов из папки public/static
	fs := http.FileServer(http.Dir("./static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	cssFileServer := http.StripPrefix("/static/", fs)
	http.Handle("/static/", cssFileServer)
	http.HandleFunc("/", urlhome) //Страница управления

	go func() {
		log.Println("Starting HTTP-server...")
		log.Fatalln("WebCtl:", srv.ListenAndServe())
	}()

	go func() {
		<-Quit
		fmt.Println("Shutting down HTTP-server...")
		ctx, _ := context.WithTimeout(context.Background(), 1 * time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalln("HTTP Shutdown error:", err)
		}
	}()
	w.islisten = true
	return err
}

//Обработчик запросов для home
func urlhome(w http.ResponseWriter, r *http.Request) {
	title := "TELEXGEN GO"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	//page := Page{title, template.HTML(body), lnkhome, "" }
	now := time.Now()
	datenow := now.Format("02/01/2006")
	page := Page{title, template.HTML(body), lnkhome, template.HTML(datenow), SeqCnt}

	if r.Method == "GET" {
		if err := home_template.ExecuteTemplate(w, "main", page); err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	} else {
		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()

		// Массив данных JSON для получения данных из формы (ajax)
		// Первый элемент должен содержать код действия
		type jsonPOSTData map[string]string

		var jh jsonPOSTData

		err := dec.Decode(&jh)
		if err != nil {
			log.Println("Handshake error: ", err)
		}

		fmt.Println(jh)

		enc := json.NewEncoder(w)
		switch jh["Post"] {
		case "SaveButton":
			sum, err := strconv.Atoi(jh["suminput"])
			if (err != nil) || (sum <= 0) {
				enc.Encode("SaveNotOkSUM")
				break
			}
			cnt, err := strconv.Atoi(jh["seqcounter"])
			if (err != nil) || (cnt <= 0 && cnt >= 128) {
				enc.Encode("SaveNotOk")
				break
			}
			key := CalcKey(int64(sum), jh["selectcur"], SeqCnt, false, 0)
			prefix := fmt.Sprintf("%03d", SeqCnt)
			telexkey := prefix + strconv.Itoa(key)

			err = Msg.SetParams(int64(sum), jh["selectcur"], jh["textarea"], jh["dateinput"], cnt, telexkey)
			if err != nil {
				enc.Encode("Ошибка: " + err.Error())
				break
			}
			err = WriteTelex()
			if err != nil {
				enc.Encode("Ошибка создания файла сообщения: " + err.Error())
				break
			}
			enc.Encode("SaveOk")

		case "CalcButton":
			sum, err := strconv.Atoi(jh["suminput"])
			if (err != nil) || (sum <= 0) {
				enc.Encode("CalcNotOk")
				break
			}

			cnt, err := strconv.Atoi(jh["seqcounter"])
			if (err != nil) || (cnt <= 0 && cnt >= 128) {
				enc.Encode("CalcNotOk")
				break
			}

			key := CalcKey(int64(sum), jh["selectcur"], SeqCnt, false, 0)
			prefix := fmt.Sprintf("%03d", SeqCnt)
			telexkey := prefix + strconv.Itoa(key)

			err = Msg.SetParams(int64(sum), jh["selectcur"], jh["textarea"], jh["dateinput"], cnt, telexkey)
			if err != nil {
				enc.Encode("Ошибка: " + err.Error())
				break
			}

			enc.Encode([]string{"CalcOk", telexkey, CalcLog})
		case "ExitButton":
			enc.Encode("ExitOk")
			WaitExit = true
		default:
			//Отправляем ответ на POST-запрос
			//для предотвращения ошибки JSON parse error в ajax методе
			enc.Encode("No action requested.")
		}
	}
}


//Функции установки значений
func (w *WebCtl) SetHost(host net.IP) {
	w.host = host
}

func (w *WebCtl) SetPort(port uint16) {
	w.port = port
}

/**/
func (w WebCtl) connString() string {
	return fmt.Sprintf("%s:%d", w.host.String(), w.port)
}

func (c *Config) SetManagerSrv(addr string, port uint16) {
	c.managerSrv = managerSrv{
		Addr: addr,
		Port: port,
	}
}

func (c *Config) ManagerSrvAddr() string {
	return c.managerSrv.Addr
}

func (c *Config) ManagerSrvPort() uint16 {
	return c.managerSrv.Port
}
