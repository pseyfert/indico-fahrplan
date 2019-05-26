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

package indicoinput

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"
)

func Cernquery(eventid int, extravals map[string]string) (Apiresult, error) {
	var data Apiresult
	c := http.Client{
		Timeout: 600 * time.Second,
	}

	requrl := fmt.Sprintf("https://indico.cern.ch/export/event/%d.xml", eventid)
	req, err := http.NewRequest("GET", requrl, nil)

	q := req.URL.Query()
	q.Add("detail", "contributions")
	for k, v := range extravals {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return data, err
	}

	req.Header.Set("User-Agent", "fahrplan-app")
	response, err := c.Do(req)
	if err != nil {
		return data, err
	}

	// b, err := ioutil.ReadAll(response.Body)
	// fmt.Printf("%s\n", b)
	xmldecoder := xml.NewDecoder(response.Body)
	err = xmldecoder.Decode(&data)
	return data, err
}

func ReadFile(fname string) (Apiresult, error) {
	var data Apiresult
	xmlfile, err := os.Open(fname)
	if err != nil {
		return data, err
	}
	xmldecoder := xml.NewDecoder(xmlfile)
	err = xmldecoder.Decode(&data)
	return data, err
}
