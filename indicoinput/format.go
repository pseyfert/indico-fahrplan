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
	"time"
)

type Apiresult struct {
	XMLName xml.Name `xml:"httpapiresult"`
	Count   int      `xml:"count"`
	Results Results  `xml:"results"`
}

func (data *Apiresult) FillTimes() {
	for i, contrib := range data.Results.Conference.Contributions.Contributions {
		t, err := time.Parse(time.RFC3339, contrib.Start)
		if err == nil {
			data.Results.Conference.Contributions.Contributions[i].StartTime = t
			data.Results.Conference.Contributions.Contributions[i].TimeLess = false
		} else {
			data.Results.Conference.Contributions.Contributions[i].TimeLess = true
		}
	}
}

type Results struct {
	XMLName    xml.Name   `xml:"results"`
	Conference Conference `xml:"conference"`
}

type Conference struct {
	XMLName       xml.Name      `xml:"conference"`
	Id            int           `xml:"id,attr"`
	Contributions Contributions `xml:"contributions"`
}

type Contributions struct {
	XMLName       xml.Name       `xml:"contributions"`
	Contributions []Contribution `xml:"contribution"`
}

type Contribution struct {
	XMLName xml.Name `xml:"contribution"`
	Id      int      `xml:"id,attr"`
	// Start       time.Time `xml:"startDate"` // datetime
	Start       string `xml:"startDate"` // datetime
	StartTime   time.Time
	TimeLess    bool
	Duration    int      `xml:"duration"`
	Title       string   `xml:"title"`
	Location    string   `xml:"location"`
	Description string   `xml:"description"`
	Speakers    Speakers `xml:"speakers"`
	// speakers/[]contributionparticipation
	// primaryauthors/[]contributionparticipation
	// folders/[]folder/attachments/[]attachment
}

type Speakers struct {
	XMLName                   xml.Name                    `xml:"speakers"`
	Contributionparticipation []Contributionparticipation `xml:"contributionparticipation"`
}
type Contributionparticipation struct {
	XMLName    xml.Name `xml:"contributionparticipation"`
	Last_name  string
	First_name string
	FullName   string
}
