package main

import (
	"fmt"
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

func main() {
	sqlcon := config.NewDatabaseConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FSeoul",
		sqlcon.User, sqlcon.Password, sqlcon.Host, sqlcon.Port, sqlcon.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}

	// 자동 마이그레이션
	if err := db.AutoMigrate(&model.Team{}, &model.Hist{}); err != nil {
		log.Fatal("테이블 마이그레이션 실패:", err)
	}

	// Gin 라우터 초기화
	r := gin.Default()

	teamService := service.NewTeamService(db)
	// 템플릿 로드
	r.LoadHTMLGlob("templates/*")

	// 정적 파일 제공
	r.Static("/static", "./static")

	// 랜딩 페이지
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "포인트 계산기",
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

	// 팀 API 엔드포인트
	api := r.Group("/api")
	{
		// 팀 생성
		api.POST("/teams", func(c *gin.Context) {
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
		api.PUT("/teams/:id", func(c *gin.Context) {
			id := c.Param("id")
			var team model.Team

			if err := c.ShouldBindJSON(&team); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			team.ID, _ = strconv.Atoi(id)

			if _, err := teamService.UpdateTeam(team); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "팀 수정에 실패했습니다."})
				return
			}

			c.JSON(http.StatusOK, team)
		})

		// 팀 삭제
		api.DELETE("/teams/:id", func(c *gin.Context) {
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

	// 서버 시작
	r.Run(":8080")
}
