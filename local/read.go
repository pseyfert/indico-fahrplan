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
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/pseyfert/go-http-redirect-resolve/resolve"

	"github.com/pseyfert/indico-fahrplan/indicoinput"
	"github.com/pseyfert/indico-fahrplan/process"
)

func main() {
	var eventid int
	var short string
	var local bool
	var xmlFileName string
	var apikey string
	var secret string
	flag.IntVar(&eventid, "i", 708041, "indico.cern.ch event id")
	flag.BoolVar(&local, "l", false, "run from local file")
	flag.StringVar(&xmlFileName, "f", "./samples/indico.event.detail.contributions.xml", "default filename")
	flag.StringVar(&apikey, "apikey", "", "indico api key (copy and paste from website)")
	flag.StringVar(&secret, "secretkey", "", "indico secret key (copy and paste from website)")
	flag.StringVar(&short, "shortened", "", "url shortened event name as in https://indico.cern.ch/e/<shortened>")
	flag.Parse()

	var data indicoinput.Apiresult
	var err error
	querymap := make(map[string]string)
	if apikey != "" {
		querymap["apikey"] = apikey
	}
	if short != "" {
		resolved, err := resolve.Resolve(fmt.Sprintf("https://indico.cern.ch/e/%s", short))
		if err != nil {
			log.Fatalf("%v", err)
		}
		re := regexp.MustCompile("^https://indico.cern.ch/event/([0-9].*)/$")
		if !re.MatchString(resolved) {
			// error
		}
		eventid64, err := strconv.ParseInt(re.ReplaceAllString(resolved, "$1"), 10, 64)
		if err != nil {
			// error
		}
		eventid = int(eventid64)
	}
	if !local {
		data, err = indicoinput.Cernquery(eventid, querymap, secret)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		data, err = indicoinput.ReadFile(xmlFileName)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	process.DumpFahrplan(data, os.Stdout)
}
