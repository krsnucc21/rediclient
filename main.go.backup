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

func initRouter(database *db.Database) *gin.Engine {
   r := gin.Default()

   r.POST("/rsrp", func (c *gin.Context) {
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
   })

   r.GET("/rsrp/:cell", func (c *gin.Context) {
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
      c.JSON(http.StatusOK, gin.H{"cell": userlist})
   })

   return r
}
