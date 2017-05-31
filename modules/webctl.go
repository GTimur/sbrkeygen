package keygen

import (
	"log"
	"net/http"
	"github.com/braintree/manners"
	"net"
	"fmt"
	"strconv"
	"path"
	"html/template"
	"time"
	"encoding/json"
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
}


//Глобальная переменная для хранения настроек
var (
	GlobalConfig Config = Config{}
	home_template = template.Must(template.ParseFiles(path.Join("static", "tpl", "main.gtpl"), path.Join("static", "tpl", "index.gtpl")))
	WaitExit bool
)

/*Сервер*/
//Запускает goroutine ListenAndServe
//Может изменять accbook - справочник подписантов
func (w *WebCtl) StartServe() (err error) {
	// для отдачи сервером статичных файлов из папки public/static

	fs := http.FileServer(http.Dir("./static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	cssFileServer := http.StripPrefix("/static/", fs)
	http.Handle("/static/", cssFileServer)
	http.HandleFunc("/", urlhome) //Страница управления
	go func() {
		log.Fatalln("WebCtl:", manners.ListenAndServe(w.connString(), http.DefaultServeMux))
	}()
	w.islisten = true
	return err
}

//Обработчик запросов для home
func urlhome(w http.ResponseWriter, r *http.Request) {
	title := "COURIER GO"
	body := ""
	lnkhome := "http://" + GlobalConfig.managerSrv.Addr + ":" + strconv.Itoa(int(GlobalConfig.managerSrv.Port))
	//page := Page{title, template.HTML(body), lnkhome, "" }
	now := time.Now()
	datenow := now.Format("02/01/2006")
	page := Page{title, template.HTML(body), lnkhome, template.HTML(datenow)}

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

		fmt.Println("POST:", jh["Post"])

		enc := json.NewEncoder(w)
		switch jh["Post"] {
		case "SaveButton":
			fmt.Println("SaveButton pressed")
			enc.Encode("SaveOk")
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
