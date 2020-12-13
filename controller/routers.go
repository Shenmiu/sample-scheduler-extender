package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	extender "k8s.io/kube-scheduler/extender/v1"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to sample-scheduler-extender!\n")
}

func Filter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs extender.ExtenderArgs
	var extenderFilterResult *extender.ExtenderFilterResult
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &extender.ExtenderFilterResult{
			Error: err.Error(),
		}
	} else {
		extenderFilterResult = filter(extenderArgs)
	}

	if response, err := json.Marshal(extenderFilterResult); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func Prioritize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	var extenderArgs extender.ExtenderArgs
	var hostPriorityList *extender.HostPriorityList
	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		hostPriorityList = &extender.HostPriorityList{}
	} else {
		hostPriorityList = prioritize(extenderArgs)
	}

	if response, err := json.Marshal(hostPriorityList); err != nil {
		log.Fatalln(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
