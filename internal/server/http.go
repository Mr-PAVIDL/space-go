package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"space-go/internal/model"
)

type httpServer struct {
	Server
}

func (s *httpServer) getUniverse(w http.ResponseWriter, _ *http.Request) {
	data, err := json.Marshal(s.Universe())
	if err != nil {
		fmt.Printf("failed to serve get universe: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = io.WriteString(w, string(data))
	if err != nil {
		fmt.Printf("failed to serve get universe: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) postTravel(w http.ResponseWriter, r *http.Request) {
	var req model.TravelRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("failed to serve post travel: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		fmt.Printf("failed to serve post travel: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.Travel(req)
	if err != nil {
		fmt.Printf("failed to serve post travel: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("failed to serve post travel: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = io.WriteString(w, string(data))
	if err != nil {
		fmt.Printf("failed to serve post travel: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) postCollect(w http.ResponseWriter, r *http.Request) {
	var req model.CollectRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("failed to serve post collect: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		fmt.Printf("failed to serve post collect: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.Collect(req)
	if err != nil {
		fmt.Printf("failed to serve post collect: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("failed to serve post collect: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = io.WriteString(w, string(data))
	if err != nil {
		fmt.Printf("failed to serve post collect: %s", err)
		io.WriteString(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) deleteReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

func (s *httpServer) getRounds(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

func Run(server Server) error {
	s := httpServer{server}

	http.HandleFunc("/player/universe", s.getUniverse)
	http.HandleFunc("/player/travel", s.postTravel)
	http.HandleFunc("/player/collect", s.postCollect)
	http.HandleFunc("/player/reset", s.deleteReset)
	http.HandleFunc("/player/rounds", s.getRounds)

	err := http.ListenAndServe(":3333", nil)
	return err
}
