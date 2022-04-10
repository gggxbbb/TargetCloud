package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/charts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

type TargetSchool struct {
	gorm.Model
	Name  string
	Value int
}

func main() {

	db, err := gorm.Open(sqlite.Open("targets.db"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&TargetSchool{})
	if err != nil {
		return
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.POST("/add", func(context *gin.Context) {
		name := context.PostForm("target")
		//If the name is already in the database, update the value
		var target TargetSchool
		db.Where("name = ?", name).First(&target)
		if target.Name != "" {
			target.Value += 1
			db.Save(&target)
		} else {
			db.Create(&TargetSchool{Name: name, Value: 1})
		}
		context.Redirect(http.StatusMovedPermanently, "/")
	})
	r.GET("/cloud", func(context *gin.Context) {

		var targetSchools []TargetSchool

		re := db.Find(&targetSchools)
		if re.Error != nil {
			panic(re.Error)
		}

		var wcData = map[string]interface{}{}

		for v := range targetSchools {
			wcData[targetSchools[v].Name] = targetSchools[v].Value
		}

		wc := charts.NewWordCloud()

		width := context.DefaultQuery("width", "800")
		height := context.DefaultQuery("height", "400")

		wc.SetGlobalOptions(
			charts.TitleOpts{Title: "", Subtitle: ""},
			charts.InitOpts{PageTitle: "目标高校词云", Width: width + "px", Height: height + "px"},
			charts.ToolboxOpts{Show: false},
		)
		wc.Add("", wcData)
		err := wc.Render(context.Writer)
		if err != nil {
			return
		}
	})
	err = r.Run()
	if err != nil {
		return
	}
}
