package main

// package name smartconn.cc/sibolwolf/batterymonitor
import (
    "log"
    "os"
    "os/exec"
    "time"
    "strings"
    "strconv"
    "golang.org/x/sys/unix"

)
//SysSW   "smartconn.cc/sibolwolf/syssleepwake"
//SysSWWH "smartconn.cc/sibolwolf/syssleepwake/wakehandle"
/*
var str_b string = "change@/devices/platform/axp22_board/axp22-supplyer.20/power_supply/batteryACTION=changeDEVPATH=/devices/platform/axp22_board/axp22-supplyer.20/power_supply/batterySUBSYSTEM=power_supplyPOWER_SUPPLY_NAME=batteryPOWER_SUPPLY_MODEL_NAME=batteryPOWER_SUPPLY_STATUS=DischargingPOWER_SUPPLY_PRESENT=1POWER_SUPPLY_ONLINE=0POWER_SUPPLY_HEALTH=GoodPOWER_SUPPLY_TECHNOLOGY=LiFePOWER_SUPPLY_VOLTAGE_MAX_DESIGN=4200000POWER_SUPPLY_VOLTAGE_MIN_DESIGN=3300POWER_SUPPLY_VOLTAGE_NOW=3990000POWER_SUPPLY_CURRENT_NOW=425000POWER_SUPPLY_ENERGY_FULL_DESIGN=1800POWER_SUPPLY_CAPACITY=100POWER_SUPPLY_TEMP=300SEQNUM=851"
*/

var str_f string = "/devices/platform/axp22_board/axp22-supplyer.20/power_supply/battery"

var batterypar = []string {
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
var lastchargestatusint int
var currbatterycapacityint int
var lasttime int64
var currtime int64
var timeperiod int64 = 3 // second


func BatteryStatusHandle(str_b string) {
	batterystatus := make(map[string]string)
    if strings.Contains(str_b, str_f) {
        log.Println("This is the string I want")
        var start, stop int
        var par []string
        for i := 0; i < len(batterypar)-1; i++ {
            start = strings.Index(str_b, batterypar[i])
            stop  = strings.Index(str_b, batterypar[i+1])
            par = strings.Split(str_b[start:stop], "=")
            batterystatus[par[0]] = par[1]
        }

        chargestatus := batterystatus[batterypar[7]]
        batterycapacity := batterystatus[batterypar[15]]
        chargestatus = strings.Replace(chargestatus, "\x00", "", -1)
        batterycapacity = strings.Replace(batterycapacity, "\x00", "", -1)

        // charge status handle
        chargestatusint, chargestatus_err := strconv.Atoi(chargestatus)
        if chargestatus_err != nil {
            log.Println(chargestatus_err)
        }

        if chargestatusint != lastchargestatusint {
            log.Println("Current battery charge status is:", chargestatusint)
            currtime = time.Now().Unix()
            if (currtime - lasttime) > timeperiod {
                lasttime = currtime
                lastchargestatusint = chargestatusint
                //SysSW.UpdateLockStatus("chargelock", chargestatusint)
            } else {
                lasttime = currtime
                lastchargestatusint = chargestatusint
                //SysSW.UpdateLockStatus("chargelock", chargestatusint)
            }
        }

        if chargestatusint == 1 {
            log.Println("RA got an event for power supplier")
            // wake up handle
            //SysSWWH.WakeJudgment()
        }

        // battery capacity handle
        batterycapacityint, batterycapacityint_err := strconv.Atoi(batterycapacity)
        if batterycapacityint_err != nil {
            log.Println(batterycapacityint_err)
        } else {
            currbatterycapacityint = batterycapacityint
            log.Println("Current battery capacity is:", currbatterycapacityint)
        }


    }
}

func InitBatteryOnlineStatus() {
    // Init battery online status
    cmd := exec.Command("/bin/sh", "-c", "sysint getBatteryOnline")
    bytes, err := cmd.Output()
    if err != nil {
        log.Println("Init battery online status fault: " + err.Error())
    }

    lasttime = time.Now().Unix()
    chargestatus := strings.TrimSpace(string(bytes))
    chargestatusint, err := strconv.Atoi(chargestatus)
    if err != nil {
        log.Println(err)
    } else {
        lastchargestatusint = chargestatusint
        log.Println("Current battery charge status is:", chargestatusint)
        //SysSW.UpdateLockStatus("chargelock", lastchargestatusint)
    }
}

func InitBatteryCapacity() {
    // Init battery capacity_now
    cmd := exec.Command("/bin/sh", "-c", "sysint getBatteryCapacity")
    bytes, err := cmd.Output()
    if err != nil {
        log.Println("Init battery capacity now fault: " + err.Error())
    }

    batterycapacity := strings.TrimSpace(string(bytes))
    batterycapacityint, err := strconv.Atoi(batterycapacity)
    if err != nil {
        log.Println(err)
    } else {
        currbatterycapacityint = batterycapacityint
        log.Println("Current battery capacity is:", currbatterycapacityint)
    }

}

func BatteryNetlinkMonitor() {
    // To communicate with netlink, a netlink socket must be opened. This is done using the socket() system call:
    //log.Println("Hello, this is battery monitor")
    fd, socket_err := unix.Socket(
        // Always used when opening netlink sockets.
        unix.AF_NETLINK,
        // Seemingly used interchangeably with SOCK_DGRAM,
        // but it appears not to matter which is used.
        unix.SOCK_RAW,
        // The netlink family that the socket will communicate
        // with, such as NETLINK_ROUTE or NETLINK_GENERIC.
        unix.NETLINK_KOBJECT_UEVENT,
    )

    if socket_err != nil {
        log.Println("Socket_err is:" + socket_err.Error())
    }

    // Once the socket is created, bind() must be called to prepare it to send and receive messages.
    bind_err := unix.Bind(fd, &unix.SockaddrNetlink{
        // Always used when binding netlink sockets.
        Family: unix.AF_NETLINK,
        // A bitmask of multicast groups to join on bind.
        // Typically set to zero.
        Groups: 0xffffffff,
        // If you'd like, you can assign a PID for this socket
        // here, but in my experience, it's easier to leave
        // this set to zero and let netlink assign and manage
        // PIDs on its own.
        Pid: uint32(os.Getpid()),
    })

    if bind_err != nil {
        log.Println("Bind_err is:" + bind_err.Error())
    }

    bstore := make([]byte, os.Getpagesize())

    for {
        for {
            // Peek at the buffer to see how many bytes are available.
            n, _, _ := unix.Recvfrom(fd, bstore, 4096)
            //log.Println("Length of data is:", n)
            // http://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
            //log.Println(string(bstore[:n]))
            // Break when we can read all messages.
            if n < len(bstore) {
                BatteryStatusHandle(string(bstore[:n]))
                break
            }
            // Double in size if not enough bytes.
            bstore = make([]byte, len(bstore)*2)
        }
    }
}

func GetBatteryOnlineStatus() int {
    return lastchargestatusint
}

func GetBatteryCapacityStatus() int {
    return currbatterycapacityint
}

func main() {
    // Init battery online status
    InitBatteryOnlineStatus()
    InitBatteryCapacity()
    BatteryNetlinkMonitor()
}
