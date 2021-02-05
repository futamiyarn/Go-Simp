package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/backend/fanart/bilibili"
	"github.com/JustHumanz/Go-Simp/service/backend/fanart/twitter"
	"github.com/JustHumanz/Go-Simp/service/backend/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func main() {
	Twitter := flag.Bool("TwitterFanart", false, "Enable twitter fanart module")
	BiliBili := flag.Bool("BiliBiliFanart", false, "Enable bilibili fanart module")
	flag.Parse()

	conf, err := config.ReadConfig("../../../config.toml")
	if err != nil {
		log.Panic(err)
	}

	Bot, err := discordgo.New("Bot " + config.BotConf.Discord)
	if err != nil {
		log.Error(err)
	}
	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}

	database.Start(conf.CheckSQL())
	engine.Start()

	c := cron.New()
	c.Start()

	if *Twitter {
		twitter.Start(Bot, c)
		database.ModuleInfo("TwitterFanart")
	}

	if *BiliBili {
		bilibili.Start(Bot, c)
		database.ModuleInfo("BiliBiliFanart")
	}
	runfunc.Run(Bot)
}