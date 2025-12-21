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

package rankings

import (
	"log"
	"strconv"
	"time"

	"github.com/ynoproject/ynorankings/common"
	"github.com/ynoproject/ynorankings/database"

	"github.com/go-co-op/gocron"
)

var (
	scheduler = gocron.NewScheduler(time.UTC)
)

func Init() {
	common.CurrentEventPeriodOrdinal, _ = database.GetCurrentEventPeriodOrdinal()

	for _, gameName := range common.GameNames {
		var rankingCategories []*common.RankingCategory

		bpCategory := &common.RankingCategory{CategoryId: "bp"}
		rankingCategories = append(rankingCategories, bpCategory)

		badgeCountCategory := &common.RankingCategory{CategoryId: "badgeCount"}
		rankingCategories = append(rankingCategories, badgeCountCategory)

		bpCategory.SubCategories = append(bpCategory.SubCategories, common.RankingSubCategory{SubCategoryId: "all"})
		badgeCountCategory.SubCategories = append(badgeCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: "all"})

		bpCategory.SubCategories = append(bpCategory.SubCategories, common.RankingSubCategory{SubCategoryId: gameName, Game: gameName})
		badgeCountCategory.SubCategories = append(badgeCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: gameName, Game: gameName})

		eventPeriods, err := database.GetEventPeriodData(gameName)
		if err != nil {
			log.Print("SERVER ", "exp", err.Error())
		} else if len(eventPeriods) > 0 {
			expCategory := &common.RankingCategory{CategoryId: "exp", Periodic: true}
			rankingCategories = append(rankingCategories, expCategory)

			if len(eventPeriods) > 1 {
				expCategory.SubCategories = append(expCategory.SubCategories, common.RankingSubCategory{SubCategoryId: "all"})
			}
			for _, eventPeriod := range eventPeriods {
				expCategory.SubCategories = append(expCategory.SubCategories, common.RankingSubCategory{SubCategoryId: strconv.Itoa(eventPeriod.PeriodOrdinal)})
			}

			eventLocationCountCategory := &common.RankingCategory{CategoryId: "eventLocationCount", Periodic: true}
			rankingCategories = append(rankingCategories, eventLocationCountCategory)

			if len(eventPeriods) > 1 {
				eventLocationCountCategory.SubCategories = append(eventLocationCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: "all"})
			}
			for _, eventPeriod := range eventPeriods {
				eventLocationCountCategory.SubCategories = append(eventLocationCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: strconv.Itoa(eventPeriod.PeriodOrdinal)})
			}

			freeEventLocationCountCategory := &common.RankingCategory{CategoryId: "freeEventLocationCount", Game: gameName, Periodic: true, SeparateByGame: true}
			rankingCategories = append(rankingCategories, freeEventLocationCountCategory)

			if len(eventPeriods) > 1 {
				freeEventLocationCountCategory.SubCategories = append(freeEventLocationCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: "all", Game: gameName})
			}
			for _, eventPeriod := range eventPeriods {
				freeEventLocationCountCategory.SubCategories = append(freeEventLocationCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: strconv.Itoa(eventPeriod.PeriodOrdinal), Game: gameName})
			}

			if gameName == "2kki" {
				eventLocationCompletionCategory := &common.RankingCategory{CategoryId: "eventLocationCompletion", Game: gameName, Periodic: true}
				rankingCategories = append(rankingCategories, eventLocationCompletionCategory)

				if len(eventPeriods) > 1 {
					eventLocationCompletionCategory.SubCategories = append(eventLocationCompletionCategory.SubCategories, common.RankingSubCategory{SubCategoryId: "all", Game: gameName})
				}
				for _, eventPeriod := range eventPeriods {
					eventLocationCompletionCategory.SubCategories = append(eventLocationCompletionCategory.SubCategories, common.RankingSubCategory{SubCategoryId: strconv.Itoa(eventPeriod.PeriodOrdinal), Game: gameName})
				}
			}

			eventVmCountCategory := &common.RankingCategory{CategoryId: "eventVmCount", Periodic: true}
			rankingCategories = append(rankingCategories, eventVmCountCategory)

			for _, eventPeriod := range eventPeriods {
				if eventPeriod.EnableVms {
					eventVmCountCategory.SubCategories = append(eventVmCountCategory.SubCategories, common.RankingSubCategory{SubCategoryId: strconv.Itoa(eventPeriod.PeriodOrdinal)})
				}
			}

			if len(eventVmCountCategory.SubCategories) > 1 {
				eventVmCountCategory.SubCategories = append([]common.RankingSubCategory{{SubCategoryId: "all"}}, eventVmCountCategory.SubCategories...)
			}
		}

		if gameName == "2kki" {
			timeTrialMapIds, err := database.GetTimeTrialMapIds()
			if err != nil {
				log.Print("SERVER ", "timeTrial", err.Error())
			} else if len(timeTrialMapIds) > 0 {
				timeTrialCategory := &common.RankingCategory{CategoryId: "timeTrial", Game: gameName}
				rankingCategories = append(rankingCategories, timeTrialCategory)

				for _, mapId := range timeTrialMapIds {
					timeTrialCategory.SubCategories = append(timeTrialCategory.SubCategories, common.RankingSubCategory{SubCategoryId: strconv.Itoa(mapId), Game: gameName})
				}
			}
		}

		gameMinigameIds, err := database.GetGameMinigameIds(gameName)
		if err != nil {
			log.Print("SERVER ", "minigame", err.Error())
		} else {
			minigameCategory := &common.RankingCategory{CategoryId: "minigame", Game: gameName}
			rankingCategories = append(rankingCategories, minigameCategory)

			for _, minigameId := range gameMinigameIds {
				minigameCategory.SubCategories = append(minigameCategory.SubCategories, common.RankingSubCategory{SubCategoryId: minigameId, Game: gameName})
			}
		}

		for c, category := range rankingCategories {
			categoryId := category.CategoryId
			if category.SeparateByGame {
				categoryId += "_" + category.Game
			} else if category.Periodic && category.Game == "" && gameName != "2kki" {
				continue
			}
			err := database.WriteRankingCategory(categoryId, category.Game, c)
			if err != nil {
				log.Print("SERVER ", categoryId, err.Error())
				continue
			}
			for sc, subCategory := range category.SubCategories {
				err = database.WriteRankingSubCategory(categoryId, subCategory.SubCategoryId, subCategory.Game, sc)
				if err != nil {
					log.Print("SERVER ", categoryId+"/"+subCategory.SubCategoryId, err.Error())
				}
			}
		}

		common.GameRankingCategories[gameName] = rankingCategories
	}

	scheduler.Every(15).Minutes().SingletonMode().Do(func() {
		for _, gameName := range common.GameNames {
			for _, category := range common.GameRankingCategories[gameName] {
				categoryId := category.CategoryId
				if category.SeparateByGame {
					categoryId += "_" + category.Game
				}
				for _, subCategory := range category.SubCategories {
					// Use Yume 2kki server to update 'all' rankings
					if subCategory.SubCategoryId == "all" && !category.SeparateByGame && gameName != "2kki" {
						continue
					}
					if category.Periodic && subCategory.SubCategoryId != "all" {
						eventPeriodOrdinal, errconv := strconv.Atoi(subCategory.SubCategoryId)
						if errconv != nil || eventPeriodOrdinal != common.CurrentEventPeriodOrdinal {
							continue
						}
					}

					err := database.UpdateRankingEntries(categoryId, subCategory.SubCategoryId, subCategory.Game)
					if err != nil {
						log.Print("SERVER ", gameName+"/"+categoryId+"/"+subCategory.SubCategoryId, err.Error())
					}
				}
			}

			err := database.UpdatePlayerMedals(gameName)
			if err != nil {
				log.Print("SERVER ", "medals", err.Error())
			}
		}
	})

	scheduler.StartAsync()
}
