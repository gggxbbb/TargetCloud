package main

import (
	"bytes"
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

import _ "embed"

type TargetSchool struct {
	gorm.Model
	Name  string
	Value int
}

type SubmitRecord struct {
	gorm.Model
	TargetName string
	FromIP     string
	FromUA     string
	FromRef    string
	FromName   string
}

//go:embed templates/*
var templateFiles embed.FS

func main() {

	db, err := gorm.Open(sqlite.Open("targets.db"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&TargetSchool{}, &SubmitRecord{})
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	templates, err := template.ParseFS(templateFiles, "templates/*.html")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(templates)

	// 主页面
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// 提交记录
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
		db.Create(&SubmitRecord{
			TargetName: name,
			FromIP:     context.ClientIP(),
			FromUA:     context.Request.UserAgent(),
			FromRef:    context.Request.Referer(),
		})
		context.Redirect(http.StatusMovedPermanently, "/")
	})

	// 词云 JSON API
	r.GET("/cloud", func(context *gin.Context) {

		var targetSchools []TargetSchool

		re := db.Find(&targetSchools)
		if re.Error != nil {
			//goland:noinspection GoUnhandledErrorResult
			context.Error(err)
			panic(re.Error)
		}

		var wcData = map[string]interface{}{}

		for v := range targetSchools {
			wcData[targetSchools[v].Name] = targetSchools[v].Value
		}

		/*
			wc := charts.NewWordCloud()

			width := context.DefaultQuery("width", "800")
			height := context.DefaultQuery("height", "400")

			wc.SetGlobalOptions(
				charts.TitleOpts{Title: "", Subtitle: ""},
				charts.InitOpts{PageTitle: "目标高校词云", Width: width + "px", Height: height + "px"},
				charts.ToolboxOpts{Show: false},
			)
			wc.Add("", wcData, charts.WordCloudOpts{Shape: "star"})
			err := wc.Render(context.Writer)

			if err != nil {
				return
			}
		*/

		context.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"size": len(wcData),
			"data": wcData,
		})
	})

	// 导出 xlsx 文件
	r.GET("/xlsx", func(context *gin.Context) {

		file := xlsx.NewFile()
		sheet, err := file.AddSheet("目标高校")
		if err != nil {
			panic(err)
		}

		header := sheet.AddRow()
		header.AddCell().Value = "目标高校名称"
		header.AddCell().Value = "提交人数"

		var targetSchools []TargetSchool
		re := db.Find(&targetSchools)
		if re.Error != nil {
			//goland:noinspection ALL
			context.Error(re.Error)
			panic(re.Error)
		}

		for v := range targetSchools {
			row := sheet.AddRow()
			row.AddCell().Value = targetSchools[v].Name
			row.AddCell().Value = strconv.Itoa(targetSchools[v].Value)
		}

		// 根据时间生成文件名
		now := time.Now()
		filename := "target" + now.Format("20060102150405") + ".xlsx"

		buf := new(bytes.Buffer)

		err = file.Write(buf)
		if err != nil {
			panic(err)
		}

		context.Header("Content-Disposition", "attachment; filename="+filename)
		context.Data(http.StatusOK, "application/vnd.ms-excel", buf.Bytes())

	})

	err = r.Run("0.0.0.0:8054")
	if err != nil {
		return
	}
}
