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
	"time"

	"github.com/pseyfert/indico-fahrplan/indicoinput"
)

func main() {
	var eventid int
	var local bool
	var xmlFileName string
	flag.IntVar(&eventid, "i", 708041, "indico.cern.ch event id")
	flag.BoolVar(&local, "l", false, "run from local file")
	flag.StringVar(&xmlFileName, "f", "./samples/indico.event.detail.contributions.xml", "default filename")
	flag.Parse()

	var data indicoinput.Apiresult
	var err error
	if !local {
		data, err = indicoinput.Cernquery(eventid)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else {
		data, err = indicoinput.ReadFile(xmlFileName)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	for _, contrib := range data.Results.Conference.Contributions.Contributions {
		t, err := time.Parse(time.RFC3339, contrib.Start)

		if err != nil {
			fmt.Printf("YYYY-MM-DDTHH:MM:SS+HH:MM %s\n", contrib.Title)
		} else {
			fmt.Printf("%s %s\n", t.String(), contrib.Title)
		}
	}
}
