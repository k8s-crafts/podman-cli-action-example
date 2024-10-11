// MIT License

// Copyright (c) 2024 k8s-crafts Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
