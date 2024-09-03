package main

import (
	"github.com/sqzxcv/glog"
	"time"
	//_ "github.com/sqzxcv/glog/log"
)

func main() {
	logsDir := "./logs" //GetCurrentDirectory()
	glog.SetConsole(false)
	glog.SetLevel(0)
	glog.Info("日志目录:", logsDir)
	glog.SetRolling(logsDir, "ss_backend.log", false, 10, 1024*1, 1)
	//glog.SetRollingFile(logsDir, "ss_backend.log", 10, 1024*1, 1)
	//dd := fmt.Sprintf("test22:%s", "aa")

	//glog.Error("日志目录:", dd)
	glog.FError("1ddd3etest:%d", 23)
	glog.Debug("test", 232, "aaa")
	for {
		glog.Info("{\"session_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"messages\":[{\"session_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"website_id\":\"55b4f6ab-7419-425e-a95d-56e2a51b8fa6\",\"type\":\"text\",\"from\":\"operator\",\"origin\":\"urn:crisp.im:triggers:0\",\"content\":\"当前客服繁忙, 为了能够快速解决您的问题, 请问:\\n1. 您遇到问题的设备是什么:  **iPhone/iPad, Windows, Mac os, Android** ?\\n2. 请问您用的客户端app 叫什么名字, 是否是刚刚从网站上下载的?\\n3. _请一句话描述您遇到的问题_?\\n4. 为了能更快解决您的问题, 请提供一下出现问题时, 软件的截图\\n5. 如果是不能翻墙, 请问使用的是哪个节点, 有没有更换别的节点尝试\\n\\n当您完整提供以上信息以后, 我将把您的问题转交给技术人员, 您的问题将会很快被解决\",\"preview\":[],\"mentions\":[],\"read\":\"\",\"delivered\":\"\",\"fingerprint\":-1078189900,\"timestamp\":1725266031252,\"user\":{\"type\":\"website\",\"nickname\":\"Gofast\",\"avatar\":\"https://storage.crisp.chat/users/avatar/website/eaef2c5451607000/photo2021-09-18-175606_dna1fd.jpeg\"}},{\"session_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"website_id\":\"55b4f6ab-7419-425e-a95d-56e2a51b8fa6\",\"type\":\"text\",\"from\":\"user\",\"origin\":\"chat\",\"content\":\"收不到邮箱验证码\",\"preview\":[],\"mentions\":[],\"read\":\"chat\",\"delivered\":\"\",\"fingerprint\":172526603025313,\"timestamp\":1725266031754,\"user\":{\"user_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"nickname\":\"visitor250153\"}},{\"session_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"website_id\":\"55b4f6ab-7419-425e-a95d-56e2a51b8fa6\",\"type\":\"text\",\"from\":\"operator\",\"origin\":\"chat\",\"content\":\"看看垃圾箱\",\"preview\":[],\"mentions\":[],\"read\":\"chat\",\"delivered\":\"chat\",\"fingerprint\":172526645086267,\"timestamp\":1725266451324,\"user\":{\"user_id\":\"1cc725a3-e639-4823-820f-695df5faf168\",\"nickname\":\"nina JK\"}},{\"session_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"website_id\":\"55b4f6ab-7419-425e-a95d-56e2a51b8fa6\",\"type\":\"text\",\"from\":\"user\",\"origin\":\"chat\",\"content\":\"也没有\",\"preview\":[],\"mentions\":[],\"read\":\"chat\",\"delivered\":\"\",\"fingerprint\":172526647438614,\"timestamp\":1725266475403,\"user\":{\"user_id\":\"session_086de023-71cf-4243-91b5-8c9d57e256eb\",\"nickname\":\"1807413921\"}}],\"account\":\"1807413921@qq.com\"}\n阿斯蒂芬斯蒂芬是的防守打法四大四大法司法送达发啥打法是都发啥打法")
		time.Sleep(1 * time.Second)
	}
}
