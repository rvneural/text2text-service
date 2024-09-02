package server

import (
	"Text2TextService/text2text"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	SERVER = "127.0.0.1"
	PORT   = "45680"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("New Request:", r.RemoteAddr)
	data, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error while reading body:", err)
		var ans Response
		ans.OldText = err.Error()
		ans.NewText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	var Req Request
	err = json.Unmarshal(data, &Req)

	if err != nil {
		log.Println("Error while unmarshalling body:", err)
		var ans Response
		ans.OldText = err.Error()
		ans.NewText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	model := Req.Model
	promt := Req.Promt
	text := Req.Text
	temperature := Req.Temperature

	result, err := text2text.ProccessText(model, promt, text, temperature)

	if err != nil {
		log.Println("Error while processing text:", err)
		var ans Response
		ans.OldText = err.Error()
		ans.NewText = err.Error()
		byteAns, _ := json.Marshal(ans)
		w.WriteHeader(200)
		w.Write(byteAns)
		return
	}

	var Res Response
	Res.OldText = text
	Res.NewText = result

	byteAns, _ := json.Marshal(Res)
	w.WriteHeader(200)
	w.Write(byteAns)
}

func StartServer() {
	http.HandleFunc("/", handleRequest)
	log.Printf("Server started at http://%s:%s/\n", SERVER, PORT)
	err := http.ListenAndServe(SERVER+":"+PORT, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
