//go:generate goversioninfo -icon=icon.ico
package main

import (
	"fmt"
	"os"
	"strconv"
	"log"
	"time"
	"sbrkeygen/modules"
	"os/signal"
	"net"
)

func main() {
	// Баннер
	fmt.Println("SBERBANK TELEX KEY GENERATOR (C) 2017 ver.0.61")
	err := keygen.InitData()
	if err != nil {
		log.Fatal("Ошибка инициализации. Возникла проблема с одним из файлов данных:", err)
	}
	// Инициализация web-сервера
	keygen.WaitExit = false; //флаг для завершения работы
	var web keygen.WebCtl
	keygen.GlobalConfig.SetManagerSrv("127.0.0.1", 4040)
	fmt.Println("Web control configured: " + "http://" + keygen.GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(keygen.GlobalConfig.ManagerSrvPort())))
	web.SetHost(net.ParseIP(keygen.GlobalConfig.ManagerSrvAddr()))
	web.SetPort(keygen.GlobalConfig.ManagerSrvPort())

	/* Запускаем сервер обслуживания WebCtl */
	err = web.StartServe()
	if err != nil {
		log.Println("HTTP сервер: Ошибка. ", err)
		os.Exit(1)
	}

	/* Open default browser */
	err = keygen.Run("http://" + keygen.GlobalConfig.ManagerSrvAddr() + ":" + strconv.Itoa(int(keygen.GlobalConfig.ManagerSrvPort())))
	if err != nil {
		fmt.Println("Browser cannot be started, please, open webctl page manually.")
	}

	/* Перехват CTRL+C для завершения приложения */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("\nReceived %v, shutdown procedure initiated.\n\n", sig)
			keygen.Quit <- 1
			keygen.WaitExit = true
		}
	}()


	//fmt.Println("TELEX KEY IS", keygen.CalcKey(34321, "RUB", keygen.SeqCnt, false, 0))
	//fmt.Printf("%s", keygen.CalcLog)

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

