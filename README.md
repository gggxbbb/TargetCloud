# TargetCloud

目标高校词云，起源于听班会课时的脑子一抽，另类的 Golang Hello World

## 使用

```bash
git clone https://github.com/gggxbbb/TargetCloud.git
cd TargetCloud
go mod download
go build -o targetcloud
./targetcloud
```

*大概吧*

## 项目依赖

### 后端

* [Golang](https://golang.org/)
* [Git](https://git-scm.com/)
* [Gin](https://gin-gonic.com/)
* [Gorm](https://gorm.io/)
* [xlsx](https://github.com/tealeg/xlsx)

### 前端

* [picocss](https://picocss.com/)
* [wordcloud2.js](https://github.com/timdream/wordcloud2.js/)

*出于某些原因，前端引用的 CSS 和 JS 来自我个人的 `static.evax.top`，而非公共 CDN*

## 无了