package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"net/http"
	"time"

	"github.com/itsjamie/gin-cors"
	"gopkg.in/appleboy/gin-jwt.v2"
)

type User struct {
	Nome    string `json:"nome"`
	Email   string `json:"email"`
	Login   string `json:"login"`
	Senha   string `json:"senha"`
	Empresa string `json:"empresa"`
	cnpj    string `json:"cnpj"`
}

type Usuario struct {
	gorm.Model
	NomeDeUsuario    string `json:"usuario"`
	Nome             string `json:"nome"`
	Sobrenome        string `json:"sobrenome"`
	Cpf              string `json:"cpf"`
	Cnpj             string `json:"cnpj"`
	Email            string `json:"qtd_pedidos"`
	DataNascimento   string `json:"img_url"`
	Senha            string `json:""`
	Telefone         string `json:"telefone"`
	Celular          string `json:"celular"`
	Empresa          string `json:"empresa"`
	Profissao        string `json:"profissao"`
	Cidade           string `json:"cidade"`
	Estado           string `json:"estado"`
	EstadoCivil      string `json:"estado_civil"`
	Endereco         string `json:"endereco"`
	Genero           string `json:"genero"`
	Blog             string `json:"blog"`
	Site             string `json:"site"`
	Idade            string `json:"idade"`
	Pais             string `json:"pais"`
	InteressesFilme  string `json:"interesse_filme"`
	InteressesMusica string `json:"interesse_musica"`
	InteressesEvento string `json:"interesse_evento"`
	InteressesCompra string `json:"interesse_compra"`
	InteressesLocal  string `json:"interesse_local"`
}

type Produto struct {
	gorm.Model
	Nome       string `json:"nome"`
	Descricao  string `json:"descricao"`
	Codigo     string `json:"codigo"`
	QtdPedidos string `json:"qtdpedidos"`
	Logo       string `json:"imgurl"`
}

type Pedido struct {
	gorm.Model
	Nome    string `json:"nome"`
	Codigo  string `json:"codigo"`
	Usuario string `json:"usuario"`
}

func main() {

	/* NOTE: See we're using = to assign the global var
	instead of := which would assign it only in this function
	*/
	db, err = gorm.Open("sqlite3", "./gorm.db")
	//db, _ = gorm.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/database?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.AutoMigrate(&Produto{})
	db.AutoMigrate(&Pedido{})
	db.AutoMigrate(&User{})

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	//
	//
	//
	//
	//
	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			//if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
			var user User
			if err := db.Where("email = ? and senha = ?", userId,password).First(&user).Error; err != nil {
				c.AbortWithStatus(404)
				fmt.Println(err)
			} else {
				return userId, true
				c.JSON(200, userId)
			}
			return userId, false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
	}

	r.POST("/login", authMiddleware.LoginHandler)
	//  r.POST("/login", LoginAuth)
	r.POST("/registro", Cadastrar)

	//
	//
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", HelloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	r.GET("/produto", GetAllProduto)
	r.GET("/produto/:id", GetProduto)
	r.POST("/produto", CreateProduto)
	r.PUT("/produto/:id", UpdateProduto)
	r.DELETE("/produto/:id", DeleteProduto)

	r.GET("/pedido", GetAllPedido)
	r.GET("/pedido/:id", GetPedido)
	r.POST("/pedido", CreatePedido)
	r.DELETE("/pedido/:id", DeletePedido)

	//r.Run("0.0.0.0:4000")
	r.Run(":" + os.Getenv("PORT"))
}

func DeleteProduto(c *gin.Context) {
	id := c.Params.ByName("id")
	var produto Produto
	d := db.Where("id = ?", id).Delete(&produto)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func DeletePedido(c *gin.Context) {
	id := c.Params.ByName("id")
	var pedido Pedido
	d := db.Where("id = ?", id).Delete(&pedido)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func UpdateProduto(c *gin.Context) {

	var produto Produto
	id := c.Params.ByName("id")

	if err := db.Where("id = ?", id).First(&produto).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&produto)

	db.Save(&produto)
	c.JSON(200, produto)

}

func CreateProduto(c *gin.Context) {

	var produto Produto
	c.BindJSON(&produto)

	db.Create(&produto)
	c.JSON(200, produto)
}

func CreatePedido(c *gin.Context) {

	var pedido Pedido
	c.BindJSON(&pedido)

	db.Create(&pedido)
	c.JSON(200, pedido)
}

func GetProduto(c *gin.Context) {
	id := c.Params.ByName("id")
	var produto Produto
	if err := db.Where("id = ?", id).First(&produto).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, produto)
	}
}

func GetPedido(c *gin.Context) {
	id := c.Params.ByName("id")
	var pedido Pedido
	if err := db.Where("id = ?", id).First(&pedido).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, pedido)
	}
}

func GetAllProduto(c *gin.Context) {
	var produto []Produto

	//codigo := c.Request.URL.Query().Get("codigo")
	//nome := c.Request.URL.Query().Get("nome")
	query := c.Request.URL.Query().Get("q")

	//if err := db.Where("codigo LIKE ? AND nome LIKE ?", "%"+codigo+"%", "%"+nome+"%").Find(&produto).Error; err != nil {
	if err := db.Where("codigo LIKE ?", "%"+query+"%").Or("nome LIKE ?", "%"+query+"%").Find(&produto).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, produto)
	}
}

func GetAllPedido(c *gin.Context) {
	var pedido []Pedido

	codigo := c.Request.URL.Query().Get("codigo")
	nome := c.Request.URL.Query().Get("nome")

	if err := db.Where("codigo LIKE ? AND nome LIKE ?", "%"+codigo+"%", "%"+nome+"%").Find(&pedido).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, pedido)
	}

}


func LoginAuth(c *gin.Context) {
	var json User
	if c.BindJSON(&json) == nil {
		result := db.Where("Login = ?", json.Login).Find(&json)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusOK, gin.H{"status": result})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"status": result})
		}
	}
}

func HelloHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"text": "Hello World.",
	})
}

var db *gorm.DB
var err error

/*
r.GET("/confirm", func(c *gin.Context) {
    token := c.Request.URL.Query().Get("token")
  }

## http://blog.gaku.net/get-url-parameter-with-gingolang/ ##
*/


func Cadastrar(c *gin.Context) {

	var user User
	c.BindJSON(&user)

	db.Create(&user)
	c.JSON(200, user)
}

