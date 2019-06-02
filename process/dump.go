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

package process

import (
	"encoding/xml"
	"fmt"
	dayvider "github.com/pseyfert/go-dayvider"
	"github.com/pseyfert/indico-fahrplan/indicoinput"
	"io"
	"time"
)

func DumpBlocks(data indicoinput.Apiresult, w io.Writer) {
	data.Parse()
	event := IndicoToDayvider(data)
	blocks := event.Blockify()

	for _, block := range blocks {
		fmt.Fprintf(w, "session block from %s to %s\n", block.Start.String(), block.End.String())
	}
	eod, err := dayvider.EndOfFirstDay(blocks)
	if err != nil {
		fmt.Fprintf(w, "can't establish day splitting on event\n")
		return
	}

	fmt.Fprintf(w, "\n\nFirst day should end at %s\n", eod.String())

	days, err := IndicoToDays(data)
	if err != nil {
		fmt.Fprintf(w, "unexpected error: %v\n", err)
		return
	}
	for i, day := range days {
		fmt.Fprintf(w, " === Day %d === \n", i+1)
		for _, contrib := range day {
			fmt.Fprintf(w, "%s - %s: %s\n", contrib.StartTime.Format("15:04"), contrib.EndTime.Format("15:04"), contrib.Title)
		}
	}
}

func Dump(data indicoinput.Apiresult, w io.Writer) {
	for _, contrib := range data.Results.Conference.Contributions.Contributions {
		t, err := time.Parse(time.RFC3339, contrib.Start)

		if err != nil {
			fmt.Fprintf(w, "YYYY-MM-DDTHH:MM:SS+HH:MM %s\n", contrib.Title)
		} else {
			fmt.Fprintf(w, "%s %s\n", t.String(), contrib.Title)
		}
	}
}

func DumpFahrplan(data indicoinput.Apiresult, w io.Writer) {
	fahrplan, err := FahrplanFromApi(data)
	if err != nil {
		fmt.Fprintf(w, "an error occured:\n%v\n", err)
		return
	}

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err := enc.Encode(fahrplan); err != nil {
		fmt.Fprintf(w, "encoding error: %v\n", err)
	}
	fmt.Fprintf(w, "\n")
}
