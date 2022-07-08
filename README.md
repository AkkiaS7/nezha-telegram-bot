# nezha-telegram-bot
A telegram bot for NeZha Monitor

## 功能

- 服务器状态简述
- 性能排名
- 自动撤回消息

## 运行方法

```bash
wget -O docker-compose.yaml https://github.com/AkkiaS7/nezha-telegram-bot/blob/master/docker-compose.yaml
mkdir data && wget -O data/config.ini https://raw.githubusercontent.com/AkkiaS7/nezha-telegram-bot/master/data/config.example.ini
cp data/config.example.ini data/config.ini
nano data/config.ini
```
修改完成后，运行 `docker-compose up -d`


## Credits

https://github.com/naiba/nezha
https://github.com/tucnak/telebot