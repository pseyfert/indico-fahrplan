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

type Results struct {
	XMLName    xml.Name   `xml:"results"`
	Conference Conference `xml:"conference"`
}

type Conference struct {
	XMLName          xml.Name      `xml:"conference"`
	Id               int           `xml:"id,attr"`
	Contributions    Contributions `xml:"contributions"`
	Start            string        `xml:"startDate"` // datetime
	StartTime        time.Time
	End              string `xml:"endDate"` // datetime
	EndTime          time.Time
	Timezone         string `xml:"timezone"`
	TimezoneLocation *time.Location
	Room             string `xml:"room"` // FIXME: validate with acat
	Url              string `xml:"url"`  // FIXME: use and validate with acat
}

func (data *Apiresult) Parse() {
	data.Results.Conference.parse()
}

func (c *Contribution) parse() {
	t, err := time.Parse(time.RFC3339, c.Start)
	if err == nil {
		c.StartTime = t
		c.EndTime = t.Add(time.Duration(c.Duration) * time.Minute)
		c.TimeLess = false
	} else {
		c.TimeLess = true
	}
	if c.RoomFullname != "" {
		c.CombinedLocation = c.Location + ": " + c.RoomFullname
	} else {
		c.CombinedLocation = c.Location
	}
}

func (c *Conference) parse() {
	t, err := time.Parse(time.RFC3339, c.Start)
	if err == nil {
		c.StartTime = t
	}
	t, err = time.Parse(time.RFC3339, c.End)
	if err == nil {
		c.EndTime = t
	}
	for i, _ := range c.Contributions.Contributions {
		c.Contributions.Contributions[i].parse()
		if c.Contributions.Contributions[i].CombinedLocation == "" {
			c.Contributions.Contributions[i].CombinedLocation = c.Room
		}
	}
	location, err := time.LoadLocation(c.Timezone)
	if err == nil {
		c.TimezoneLocation = location
	}
}

type Contributions struct {
	XMLName       xml.Name       `xml:"contributions"`
	Contributions []Contribution `xml:"contribution"`
}

type Contribution struct {
	XMLName xml.Name `xml:"contribution"`
	Id      int      `xml:"id,attr"`
	// Start       time.Time `xml:"startDate"` // datetime
	Start            string `xml:"startDate"` // datetime
	StartTime        time.Time
	EndTime          time.Time
	TimeLess         bool
	Duration         int    `xml:"duration"`
	Title            string `xml:"title"`
	Location         string `xml:"location"`
	RoomFullname     string `xml:"roomFullname"`
	CombinedLocation string
	Description      string   `xml:"description"`
	Speakers         Speakers `xml:"speakers"`
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
