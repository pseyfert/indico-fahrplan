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
	"github.com/pseyfert/indico-fahrplan/indicoinput"
	"io"
	"time"
)

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

func main() {
	fmt.Println("vim-go")
}
