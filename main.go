// Copyright 2010 The Go Authors. All rights reserved.

// Use of this source code is governed by a BSD-style

// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"

	"log"

	"net/http"

	"regexp"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HealthInfo struct {
	Version string `json:"version"`
	Ready   bool   `json:"ready"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	info := &HealthInfo{
		Version: "v0.0.0-unset",
		Ready:   true,
	}

	content, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Write(content)

}

func podHandler(w http.ResponseWriter, r *http.Request) {

	fakePod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-pod",
			Namespace: "my-namespace",
		},
	}

	content, err := json.Marshal(fakePod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Write(content)

}

var validPath = regexp.MustCompile("^/(health|pods)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {

			http.NotFound(w, r)

			return

		}

		fn(w, r)

	}

}

func main() {

	http.HandleFunc("/health", makeHandler(healthHandler))

	http.HandleFunc("/pods", makeHandler(podHandler))

	fmt.Println("Server is running at 0.0.0.0:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
