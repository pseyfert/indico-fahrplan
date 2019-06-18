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

package fahrplanoutput

import (
	"encoding/xml"
	"time"
)

func (s *Schedule) Finalize() {
	if s.Version == "" {
		s.Version = "HEPWorldGomination"
	}
}

type Schedule struct {
	XMLName    xml.Name   `xml:"schedule"`
	Conference Conference `xml:"conference"`
	Days       []Day      `xml:"day"`
	Version    string     `xml:"version"`
}

type Day struct {
	XMLName xml.Name  `xml:"day"`
	Date    string    `xml:"date,attr"` // ISO-9601 day
	Index   int       `xml:"index,attr"`
	Start   time.Time `xml:"start,attr"`
	End     time.Time `xml:"end,attr"`
	Rooms   []Room    `xml:"room"`
}

type Room struct {
	XMLName xml.Name `xml:"room"`
	Name    string   `xml:"name,attr"`
	Events  []Event  `xml:"event"`
}

type Event struct {
	XMLName     xml.Name    `xml:"event"`
	Guid        string      `xml:"guid,attr"`
	Id          int         `xml:"id,attr"`
	Date        time.Time   `xml:"date"`
	Start       string      `xml:"start"`
	Duration    string      `xml:"duration"` // hh:mm
	Title       string      `xml:"title"`
	Track       string      `xml:"track"`
	Abstract    string      `xml:"abstract"`
	Attachments Attachments `xml:"links"`
	// Attachments Attachments `xml:"attachments"`
	// subtitle, description
	Persons Persons `xml:"persons"`
}

type Persons struct {
	XMLName xml.Name `xml:"persons"`
	Persons []Person `xml:"person"`
}

type Person struct {
	XMLName xml.Name `xml:"person"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:",chardata"`
}

type Attachments struct {
	XMLName     xml.Name     `xml:"links"`
	Attachments []Attachment `xml:"link"`
	// XMLName     xml.Name     `xml:"attachments"`
	// Attachments []Attachment `xml:"attachment"`
}

type Attachment struct {
	XMLName xml.Name `xml:"link"`
	// XMLName xml.Name `xml:"attachment"`
	Href string `xml:"href,attr"`
	Name string `xml:",chardata"`
}

type Conference struct {
	XMLName      xml.Name `xml:"conference"`
	Acronym      string   `xml:"acronym"`
	Title        string   `xml:"title"`
	Start        string   `xml:"start"` // ISO 8601 day
	End          string   `xml:"end"`   // ISO 8601 day
	Days         int      `xml:"days"`
	SlotDuration string   `xml:"timeslot_duration"` // hh:mm
	Url          string   `xml:"base_url"`          // https://indico.cern.ch/event/ID
}
