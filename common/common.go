/*
	Copyright (C) 2021-2025  The YNOproject Developers
	Copyright (C) 2025       Collective Unconscious Developers

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package common

import (
	"time"
)

var (
	GameNames                 = []string{"2kki", "amillusion", "braingirl", "deepdreams", "flow", "genie", "if", "mikan", "muma", "nostalgic", "oversomnia", "prayers", "sheawaits", "someday", "tsushin", "ultraviolet", "unaccomplished", "unconscious", "unevendream", "yume"}
	GameRankingCategories     = make(map[string][]*RankingCategory)
	CurrentEventPeriodOrdinal = -1
)

type EventPeriod struct {
	PeriodOrdinal int       `json:"periodOrdinal"`
	EndDate       time.Time `json:"endDate"`
	EnableVms     bool      `json:"enableVms"`
}

type RankingCategory struct {
	CategoryId     string               `json:"categoryId"`
	Game           string               `json:"game"`
	SubCategories  []RankingSubCategory `json:"subCategories"`
	Periodic       bool                 `json:"periodic"`
	SeparateByGame bool                 `json:"-"`
}

type RankingSubCategory struct {
	SubCategoryId string `json:"subCategoryId"`
	Game          string `json:"game"`
	PageCount     int    `json:"pageCount"`
}

type Ranking struct {
	Position   int     `json:"position"`
	Name       string  `json:"name"`
	Rank       int     `json:"rank"`
	Badge      string  `json:"badge"`
	SystemName string  `json:"systemName"`
	Medals     [5]int  `json:"medals"`
	ValueInt   int     `json:"valueInt"`
	ValueFloat float32 `json:"valueFloat"`
}

type RankingEntry struct {
	CategoryId     string
	SubCategoryId  string
	Position       int
	ActualPosition int
	Uuid           string
	ValueInt       int
	ValueFloat     float32
	Timestamp      time.Time
}
