package main

import (
   "github.com/gin-gonic/gin"
   "rediclient/db"
   "net/http"
   "log"
   "os"
)

var (
   ListenAddr = ":" + os.Getenv("PORT")
   RedisAddr = os.Getenv("REDIS_ADDR")
)

func main() {
   database, err := db.NewDatabase(RedisAddr, 0)
   if err != nil {
      log.Fatalf("Failed to connect to redis: %s", err.Error())
   }

   router := initRouter(database)
   router.Run(ListenAddr)
}

var result int = 0
var curSum int = 0
var numGet int = 0
var numPost int = 0
var database *db.Database

func postFunction(c *gin.Context) {
   var userJson db.User
   if err := c.ShouldBindJSON(&userJson); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }
   err := database.SaveUser(&userJson)
   if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
   }
   c.JSON(http.StatusOK, gin.H{"user": userJson})

   num := database.GetNumLoop()
   for i := 0; i < num; i++ {
      for j := 0; j < num; j ++ {
         result += i * j
      }
      curSum++
   }
   numPost++
}

func getFunction(c *gin.Context) {
   cellname := c.Param("cell")
   userlist, err := database.GetCellUser(cellname)
   if err != nil {
      if err == db.ErrNil {
         c.JSON(http.StatusNotFound, gin.H{"error": "No record found for " + cellname})
         return
      }
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
   }

   num := database.GetNumLoop()
   c.JSON(http.StatusOK, gin.H{"cell": userlist, "numLoop": num})

   for i := 0; i < num; i++ {
      for j := 0; j < num; j ++ {
         result += i * j
      }
      curSum++
   }
   numGet++
}

func getWorkFunction(c *gin.Context) {
   numLoop := database.GetNumLoop()
   c.JSON(http.StatusOK, gin.H{"load": numLoop, "sum": curSum, "numPost": numPost, "numGet": numGet, "result": result})
}

func initRouter(dbase *db.Database) *gin.Engine {
   r := gin.Default()

   database = dbase
   r.POST("/rsrp", postFunction)
   r.GET("/rsrp/:cell", getFunction)
   r.GET("/work", getWorkFunction)

   return r
}
