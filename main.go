package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"PointCalculator/config"
	"PointCalculator/model"
	"PointCalculator/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 매치 결과 처리 API
func handleMatchResult(c *gin.Context, db *gorm.DB) {
	// URL에서 매치 ID 추출
	matchId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid match ID",
		})
		return
	}

	// 요청 본문 파싱
	var request struct {
		WinnerTeamId int `json:"winnerTeamId"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// 매치 결과 처리
	matchService := service.NewMatchService(db)
	err = matchService.ProcessMatchResult(matchId, request.WinnerTeamId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 성공 응답
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Match result processed successfully",
	})
}

func main() {
	sqlcon := config.NewDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FSeoul",
		sqlcon.User, sqlcon.Password, sqlcon.Host, sqlcon.Port, sqlcon.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}

	// 자동 마이그레이션
	if err := db.AutoMigrate(&model.Team{}, &model.Game{}, &model.Hist{}, &model.Bet{}, &model.Match{}); err != nil {
		log.Fatal("테이블 마이그레이션 실패:", err)
	}

	// Gin 라우터 초기화
	r := gin.Default()

	teamService := service.NewTeamService(db)
	gameService := service.NewGameService(db)
	matchService := service.NewMatchService(db)
	betService := service.NewBetService(db)

	// 템플릿 로드
	r.LoadHTMLGlob("templates/*")

	// 정적 파일 제공
	r.Static("/static", "./static")

	// 랜딩 페이지
	r.GET("/", func(c *gin.Context) {
		// 팀 목록과 포인트 현황 조회
		teams, err := teamService.GetTeamList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "팀 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		// 각 팀의 진행중인 베팅 수 조회
		for i := range teams {
			activeBets, err := betService.GetActiveBetCountByTeamId(teams[i].ID)
			if err != nil {
				teams[i].ActiveBets = 0
			} else {
				teams[i].ActiveBets = activeBets
			}
		}

		// 진행중인 매치 목록 조회
		matches, err := matchService.GetActiveMatches()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "매치 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		// 각 매치의 베팅 수 조회
		for i := range matches {
			betCount, err := betService.GetBetCountByMatchId(matches[i]["ID"].(int))
			if err != nil {
				matches[i]["BetCount"] = 0
			} else {
				matches[i]["BetCount"] = betCount
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   "홈",
			"teams":   teams,
			"matches": matches,
		})
	})

	// 팀 관리 페이지
	r.GET("/teams", func(c *gin.Context) {
		result, err := teamService.GetTeamList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "팀 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		c.HTML(http.StatusOK, "teams.html", gin.H{
			"title": "팀 관리",
			"teams": result,
		})
	})

	// 게임 관리 페이지
	r.GET("/games", func(c *gin.Context) {
		result, err := gameService.GetGameList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "게임 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		c.HTML(http.StatusOK, "games.html", gin.H{
			"title": "게임 관리",
			"games": result,
		})
	})

	// 매치 관리 페이지
	r.GET("/matches", func(c *gin.Context) {
		// 매치 목록 조회
		matches, err := matchService.GetMatchList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "매치 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		// 게임 목록 조회
		games, err := gameService.GetGameList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "게임 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		// 팀 목록 조회
		teams, err := teamService.GetTeamList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "팀 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		c.HTML(http.StatusOK, "matches.html", gin.H{
			"title":   "매치 관리",
			"matches": matches,
			"games":   games,
			"teams":   teams,
		})
	})

	// 베팅 관리 페이지
	r.GET("/bets", func(c *gin.Context) {
		// 베팅 목록 조회
		bets, err := betService.GetBetList()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"error": "베팅 목록을 불러오는데 실패했습니다.",
			})
			return
		}

		// 베팅 정보에 팀 이름과 게임 이름 추가
		var betsWithDetails []gin.H
		for _, bet := range bets {
			// 베팅 팀 정보 조회
			team, err := teamService.GetTeam(bet.TeamId)
			if err != nil {
				continue
			}

			// 베팅 대상 팀 정보 조회
			targetTeam, err := teamService.GetTeam(bet.TargetTeamId)
			if err != nil {
				continue
			}

			// 매치 정보 조회
			match, err := matchService.GetMatch(bet.MatchID)
			if err != nil {
				continue
			}

			// 게임 정보 조회
			game, err := gameService.GetGame(match.GameId)
			if err != nil {
				continue
			}

			betInfo := gin.H{
				"ID":             bet.ID,
				"MatchID":        bet.MatchID,
				"GameName":       game.Name,
				"TeamName":       team.Name,
				"TargetTeamName": targetTeam.Name,
				"BetType":        bet.BetType,
				"BettingPoint":   bet.BettingPoint,
				"Status":         bet.Status,
			}
			betsWithDetails = append(betsWithDetails, betInfo)
		}

		c.HTML(http.StatusOK, "bets.html", gin.H{
			"title": "베팅 목록",
			"bets":  betsWithDetails,
		})
	})

	// 베팅 API 엔드포인트
	betApi := r.Group("/api")
	{
		// 베팅하기
		betApi.POST("/bets", func(c *gin.Context) {
			var bet model.Bet
			if err := c.ShouldBindJSON(&bet); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if _, err := betService.CreateBet(bet); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, bet)
		})

		// 베팅 정보조회(단일)
		betApi.GET("/bets/:id", func(c *gin.Context) {
			id := c.Param("id")
			betID, _ := strconv.Atoi(id)
			bet, err := betService.GetBet(betID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, bet)
		})

		// 베팅 삭제
		betApi.DELETE("/bets/:id", func(c *gin.Context) {
			id := c.Param("id")
			betID, _ := strconv.Atoi(id)

			if err := betService.DeleteBet(betID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "베팅이 삭제되었습니다."})
		})
	}

	// 매치 API 엔드포인트
	matchApi := r.Group("/api")
	{
		// 매치별 베팅 목록 조회
		matchApi.GET("/matches/:id/bets", func(c *gin.Context) {
			matchId, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 매치 ID입니다."})
				return
			}

			bets, err := betService.GetBetsByMatchId(matchId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// 팀 정보 조회
			var betsWithTeamNames []gin.H
			for _, bet := range bets {
				team, err := teamService.GetTeam(bet.TeamId)
				if err != nil {
					continue
				}
				targetTeam, err := teamService.GetTeam(bet.TargetTeamId)
				if err != nil {
					continue
				}

				betInfo := gin.H{
					"id":               bet.ID,
					"match_id":         bet.MatchID,
					"team_id":          bet.TeamId,
					"team_name":        team.Name,
					"target_team_id":   bet.TargetTeamId,
					"target_team_name": targetTeam.Name,
					"betting_point":    bet.BettingPoint,
					"status":           bet.Status,
					"created_at":       bet.CreatedAt,
					"updated_at":       bet.UpdatedAt,
				}
				betsWithTeamNames = append(betsWithTeamNames, betInfo)
			}

			c.JSON(http.StatusOK, betsWithTeamNames)
		})

		// 매치 정보조회(단일)
		matchApi.GET("/matches/:id", func(c *gin.Context) {
			id := c.Param("id")
			matchID, _ := strconv.Atoi(id)
			match, err := matchService.GetMatch(matchID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, match)
		})

		// 매치 생성
		matchApi.POST("/matches", func(c *gin.Context) {
			var match model.Match

			// 요청 바디를 직접 읽어서 로깅
			body, err := c.GetRawData()
			if err != nil {
				fmt.Printf("요청 바디 읽기 실패: %v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}
			fmt.Printf("원본 요청 데이터: %s\n", string(body))

			// 바디를 다시 설정
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

			if err := c.ShouldBindJSON(&match); err != nil {
				fmt.Printf("매치 생성 요청 바인딩 실패: %v\n", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			fmt.Printf("매치 생성 요청 데이터: %+v\n", match)

			if _, err := matchService.CreateMatch(match); err != nil {
				fmt.Printf("매치 생성 실패: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, match)
		})

		// 매치 수정
		matchApi.PUT("/matches/:id", func(c *gin.Context) {
			id := c.Param("id")
			var match model.Match

			if err := c.ShouldBindJSON(&match); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			match.ID, _ = strconv.Atoi(id)

			if _, err := matchService.UpdateMatch(match); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, match)
		})

		// 매치 삭제
		matchApi.DELETE("/matches/:id", func(c *gin.Context) {
			id := c.Param("id")
			var match model.Match

			match.ID, _ = strconv.Atoi(id)

			if _, err := matchService.DeleteMatch(match); err != nil {
				fmt.Printf("매치 삭제 실패: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "매치가 삭제되었습니다."})
		})
	}

	// 게임 API 엔드포인트
	gameApi := r.Group("/api")
	{
		// 게임 생성
		gameApi.POST("/games", func(c *gin.Context) {
			var game model.Game
			if err := c.ShouldBindJSON(&game); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if _, err := gameService.CreateGame(game.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, game)
		})

		// 게임 수정
		gameApi.PUT("/games/:id", func(c *gin.Context) {
			id := c.Param("id")
			var game model.Game

			if err := c.ShouldBindJSON(&game); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			game.ID, _ = strconv.Atoi(id)

			if _, err := gameService.UpdateGame(game); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "게임 수정에 실패했습니다."})
				return
			}

			c.JSON(http.StatusOK, game)
		})

		// 게임 삭제
		gameApi.DELETE("/games/:id", func(c *gin.Context) {
			id := c.Param("id")
			var game model.Game

			game.ID, _ = strconv.Atoi(id)

			if _, err := gameService.DeleteGame(game); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "게임 삭제에 실패했습니다."})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "게임 삭제 완료"})
		})
	}

	// 팀 API 엔드포인트
	teamApi := r.Group("/api")
	{
		// 팀 정보 조회(단일)
		teamApi.GET("/teams/:id", func(c *gin.Context) {
			id := c.Param("id")

			teamID, _ := strconv.Atoi(id)
			team, err := teamService.GetTeam(teamID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, team)
		})

		// 팀 정보 조회
		teamApi.GET("/teams", func(c *gin.Context) {
			teams, err := teamService.GetTeamList()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, teams)
		})

		// 팀 생성
		teamApi.POST("/teams", func(c *gin.Context) {
			var team model.Team
			if err := c.ShouldBindJSON(&team); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if _, err := teamService.CreateTeam(team.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, team)
		})

		// 팀 수정
		teamApi.PUT("/teams/:id", func(c *gin.Context) {
			id := c.Param("id")
			var team model.Team

			if err := c.ShouldBindJSON(&team); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			team.ID, _ = strconv.Atoi(id)

			if _, err := teamService.UpdateTeam(team); err != nil {
				if err.Error() == "team already exists" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "중복된 팀 이름입니다."})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, team)
		})

		// 팀 삭제
		teamApi.DELETE("/teams/:id", func(c *gin.Context) {
			id := c.Param("id")
			var team model.Team

			team.ID, _ = strconv.Atoi(id)

			if _, err := teamService.DeleteTeam(team); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "팀 삭제에 실패했습니다."})
				return
			}

			c.JSON(http.StatusOK, team)
		})
	}

	// 히스토리 페이지
	r.GET("/history", func(c *gin.Context) {
		c.HTML(http.StatusOK, "history.html", gin.H{
			"title": "히스토리",
		})
	})

	// API 라우트 추가
	r.POST("/api/matches/:id/result", func(c *gin.Context) {
		handleMatchResult(c, db)
	})

	// 서버 시작
	r.Run(":8080")
}
