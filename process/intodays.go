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
	"fmt"
	"time"

	dayvider "github.com/pseyfert/go-dayvider"
	"github.com/pseyfert/indico-fahrplan/indicoinput"
)

func IndicoToDays(data indicoinput.Apiresult) (retval [][]indicoinput.Contribution, err error) {
	data.Parse()
	dayviderevent := IndicoToDayvider(data)
	blocks := dayviderevent.Blockify()
	refdate, err := dayvider.EndOfFirstDay(blocks)
	refdate = refdate.Add(-24 * time.Hour)

	if err != nil {
		return
	}

	days := 1 + int(blocks[len(blocks)-1].End.Sub(refdate)/(24*time.Hour))

	retval = make([][]indicoinput.Contribution, 0, days)
	for i := 0; i < days; i += 1 {
		var currentday []indicoinput.Contribution
		currentday = make([]indicoinput.Contribution, 0)
		retval = append(retval, currentday)
	}

	for _, contrib := range data.Results.Conference.Contributions.Contributions {
		if contrib.TimeLess {
			continue
		}
		cdate := int(contrib.StartTime.Sub(refdate) / (24 * time.Hour))

		if cdate < 0 || cdate >= days {
			err = fmt.Errorf("could not assign contribution to day\n"+
				"The derived day is %d\n"+
				"The event has %d days\n",
				"%s\n",
				cdate+1, days, blockinfo(blocks, contrib))
			return
		}

		retval[cdate] = append(retval[cdate], contrib)
	}

	return
}
