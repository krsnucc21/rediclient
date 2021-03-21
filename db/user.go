package db

import (
   "fmt"
   "strconv"
   "github.com/go-redis/redis/v8"
)

type User struct {
   Cellname string `json:"cellname" bining:"required"`
   Username string `json:"username" binding:"required"`
   Rsrp     int    `json:"rsrp" binding:"required"`
   Rank     int    `json:"rank"`
}

func (db *Database) SaveUser(user *User) error {
   member := &redis.Z{
      Score: float64(user.Rsrp),
      Member: user.Username,
   }
   pipe := db.Client.TxPipeline()
   pipe.ZAdd(Ctx, user.Cellname, member)
   rank := pipe.ZRank(Ctx, user.Cellname, user.Username)
   _, err := pipe.Exec(Ctx)
   if err != nil {
      return err
   }
   fmt.Println(rank.Val(), err)
   user.Rank = int(rank.Val())
   return nil
}

type Userlist struct {
   Count int `json:"count"`
   Users []*User
}

func (db *Database) GetCellUser(cellname string) (*Userlist, error) {
   scores := db.Client.ZRangeWithScores(Ctx, cellname, 0, -1)
   if scores == nil {
      return nil, ErrNil
   }
   count := len(scores.Val())
   users := make([]*User, count)
   for idx, member := range scores.Val() {
      users[idx] = &User{
         Cellname: cellname,
         Username: member.Member.(string),
         Rsrp: int(member.Score),
         Rank: idx,
      }
   }
   userlist := &Userlist{
      Count: count,
      Users: users,
   }
   return userlist, nil
}

func (db *Database) GetNumLoop() int {
   value, err := db.Client.Get(Ctx, "numloop").Result()
   if err != nil {
      return 0
   }
   num, err := strconv.Atoi(value)
   if err != nil {
      return 0
   }
   return num
}

func (db *Database) SetNumLoop(num int) error {
   value := strconv.Itoa(num)
   err := db.Client.Set(Ctx, "numloop", value, 0).Err()
   if err != nil {
      return err
   }
   return nil
}
