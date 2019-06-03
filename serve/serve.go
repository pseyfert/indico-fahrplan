/*
 * Copyright (C) 2019 Paul Seyfert
 * Author: Paul Seyfert <pseyfert.mathphys@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"github.com/pseyfert/indico-fahrplan/indicoinput"
	"github.com/pseyfert/indico-fahrplan/process"
	"log"
	"net/http"
	"strconv"
)

func doit(w http.ResponseWriter, id int, querymap map[string]string, secret string) error {
	data, err := indicoinput.Cernquery(id, querymap, secret)
	if err != nil {
		return err
	}

	process.Dump(data, w)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	s := r.FormValue("event")
	key := r.FormValue("apikey")
	secret := r.FormValue("secret")
	querymap := make(map[string]string)
	if key != "" {
		querymap["apikey"] = key
	}
	if s == "" {
		http.Error(w, fmt.Errorf("event id missing").Error(), 500)
		return
	}

	id, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		http.Error(w, fmt.Errorf("couldn't parse event id: %v", err).Error(), 500)
		return
	}

	err = doit(w, int(id), querymap, secret)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "p", 8084, "port to listen to")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
