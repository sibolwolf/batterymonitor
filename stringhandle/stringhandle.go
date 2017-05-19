package stringhandle

import (
	"log"
    "strings"
)

var str_b string = "change@/devices/platform/axp22_board/axp22-supplyer.20/power_supply/batteryACTION=changeDEVPATH=/devices/platform/axp22_board/axp22-supplyer.20/power_supply/batterySUBSYSTEM=power_supplyPOWER_SUPPLY_NAME=batteryPOWER_SUPPLY_MODEL_NAME=batteryPOWER_SUPPLY_STATUS=DischargingPOWER_SUPPLY_PRESENT=1POWER_SUPPLY_ONLINE=0POWER_SUPPLY_HEALTH=GoodPOWER_SUPPLY_TECHNOLOGY=LiFePOWER_SUPPLY_VOLTAGE_MAX_DESIGN=4200000POWER_SUPPLY_VOLTAGE_MIN_DESIGN=3300POWER_SUPPLY_VOLTAGE_NOW=3990000POWER_SUPPLY_CURRENT_NOW=425000POWER_SUPPLY_ENERGY_FULL_DESIGN=1800POWER_SUPPLY_CAPACITY=100POWER_SUPPLY_TEMP=300SEQNUM=851"

var str_f string = "/devices/platform/axp22_board/axp22-supplyer.20/power_supply/battery"

var batterpar = []string {
    "ACTION",
    "DEVPATH",
    "SUBSYSTEM",
    "POWER_SUPPLY_NAME",
    "POWER_SUPPLY_MODEL_NAME",
    "POWER_SUPPLY_STATUS",
    "POWER_SUPPLY_PRESENT",
    "POWER_SUPPLY_ONLINE",
    "POWER_SUPPLY_HEALTH",
    "POWER_SUPPLY_TECHNOLOGY",
    "POWER_SUPPLY_VOLTAGE_MAX_DESIGN",
    "POWER_SUPPLY_VOLTAGE_MIN_DESIGN",
    "POWER_SUPPLY_VOLTAGE_NOW",
    "POWER_SUPPLY_CURRENT_NOW",
    "POWER_SUPPLY_ENERGY_FULL_DESIGN",
    "POWER_SUPPLY_CAPACITY",
    "POWER_SUPPLY_TEMP",
    "SEQNUM",
}

var batterystatus map[string]string

func main() {
	batterystatus := make(map[string]string)
    if strings.Contains(str_b, str_f) {
        log.Println("This is the string I want")
        var start, stop int
        var par []string
        for i := 0; i < len(batterpar)-1; i++ {
            start = strings.Index(str_b, batterpar[i])
            stop  = strings.Index(str_b, batterpar[i+1])
            par = strings.Split(str_b[start:stop], "=")
            log.Println("start", start)
            log.Println("stop", stop)
            log.Println("par", par)
            batterystatus[par[0]] = par[1]
        }

        chargestatus := batterystatus[batterpar[7]]
        batterycapacity := batterystatus[batterypar[15]]

        log.Println("chargestatus", chargestatus)
        log.Println("batterycapacity", batterycapacity)

    }
}
