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
	"log"
	"os"

	"github.com/pseyfert/indico-fahrplan/indicoinput"
	"github.com/pseyfert/indico-fahrplan/process"
)

func main() {
	var eventid int
	var local bool
	var xmlFileName string
	var apikey string
	var signature string
	flag.IntVar(&eventid, "i", 708041, "indico.cern.ch event id")
	flag.BoolVar(&local, "l", false, "run from local file")
	flag.StringVar(&xmlFileName, "f", "./samples/indico.event.detail.contributions.xml", "default filename")
	flag.StringVar(&apikey, "apikey", "", "indico api key (copy and paste from website)")
	flag.StringVar(&signature, "signature", "", "indico signature (copy and paste from website)")
	flag.Parse()

	var data indicoinput.Apiresult
	var err error
	querymap := make(map[string]string)
	if signature != "" {
		querymap["signature"] = signature
	}
	if apikey != "" {
		querymap["apikey"] = apikey
	}
	if !local {
		data, err = indicoinput.Cernquery(eventid, querymap)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		data, err = indicoinput.ReadFile(xmlFileName)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	process.Dump(data, os.Stdout)
}
