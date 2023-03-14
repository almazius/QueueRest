package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
)

var queueMap = make(map[string][]string, 0)

func main() {
	port := flag.String("port", "127.0.0.1:4000", "Your port")
	flag.Parse()

	router := http.NewServeMux() // инициализация сервера

	router.HandleFunc("/", http.HandlerFunc(listenHandlers)) // прослушка необходимого адреса

	err := http.ListenAndServe(*port, router)
	if err != nil {
		log.Println(err)
	}
}

func listenHandlers(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		//timeOut, err := strconv.Atoi(r.URL.Query().Get("timeout"))
		//if err != nil {
		//
		//}
		//
		//timer := time.NewTimer(time.Duration(timeOut) * time.Second)

		queueName := r.URL.Path
		if res, err := getInQueue(queueName); err != nil {
			w.WriteHeader(404)
			//w.Write([]byte("not found"))
		} else {
			w.Header().Set("Content-Type", "application/json")
			b, _ := json.Marshal(res)
			if _, err := w.Write(b); err != nil {
				log.Println(err)
			}
		}
	} else if r.Method == http.MethodPut {
		param := r.URL.Query().Get("v")
		if param == "" {
			w.WriteHeader(400)
		} else {
			queueName := r.URL.Path
			putInQueue(queueName, param)
		}
	}
}

func putInQueue(queueName, param string) {
	if _, isBe := queueMap[queueName]; isBe {
		queueMap[queueName] = append(queueMap[queueName], param)
	} else {
		queueMap[queueName] = make([]string, 0)
		queueMap[queueName] = append(queueMap[queueName], param)
	}
}

func getInQueue(queueName string) (string, error) {
	if len(queueMap[queueName]) > 0 {
		res := queueMap[queueName][0]
		queueMap[queueName] = queueMap[queueName][1:]
		return res, nil
	} else {
		return "", errors.New("queue empty")
	}
}
