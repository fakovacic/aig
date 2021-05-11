autonomous-indoor-greenhouse

- controllers connect to wifi
- send/register ip address to raspberry bi
- keep track of location of controllers
- receive stats from controllers

Dashboard:
- display all controllers
- send commands to controllers

Commands:
- water on/off
- vent on/off
- lights on/off

******
Ping
 - ping controllers for stats

Config - add config field
 - when something will be done
 - depends on stats

Water  from  soil      to soil
Vent   from: temp      to temp
Lights from: time.Time to time.Time

Track config
 - go proccess to send auto commands to controllers

Camera
- track camera live
- make screenshot every hour 