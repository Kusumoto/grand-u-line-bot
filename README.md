# grand-u-line-bot
Check the register mail & check balance and send to line bot for Grand Unity condominium  

(https://github.com/Kusumoto/grand-u-line-bot/raw/master/ss.png)

## How to use
- Write the configuration file (config.json)
```json
{
     "checkRegisterMailAPIUrl": "",
    "checkCheckBalanceAPIUrl": "",
    "phoneNumber": "",
    "unitID": "",
    "projectCode": "",
    "lineChannelSecret": "",
    "lineAccessToken": ""
}
```
- Run via docker use command

```
docker run -d -v <config.json path>:/root/config.json kusumoto/grand-u-line-bot
```
