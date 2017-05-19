package main

import (
    "fmt"
    "os"
    "log"
    "strings"
    "golang.org/x/sys/unix"
)

/*
var str_b string = "change@/devices/platform/axp22_board/axp22-supplyer.20/power_supply/batteryACTION=changeDEVPATH=/devices/platform/axp22_board/axp22-supplyer.20/power_supply/batterySUBSYSTEM=power_supplyPOWER_SUPPLY_NAME=batteryPOWER_SUPPLY_MODEL_NAME=batteryPOWER_SUPPLY_STATUS=DischargingPOWER_SUPPLY_PRESENT=1POWER_SUPPLY_ONLINE=0POWER_SUPPLY_HEALTH=GoodPOWER_SUPPLY_TECHNOLOGY=LiFePOWER_SUPPLY_VOLTAGE_MAX_DESIGN=4200000POWER_SUPPLY_VOLTAGE_MIN_DESIGN=3300POWER_SUPPLY_VOLTAGE_NOW=3990000POWER_SUPPLY_CURRENT_NOW=425000POWER_SUPPLY_ENERGY_FULL_DESIGN=1800POWER_SUPPLY_CAPACITY=100POWER_SUPPLY_TEMP=300SEQNUM=851"
*/

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

func BatteryStatusHandle(str_b string) {
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

func Init() {
    // To communicate with netlink, a netlink socket must be opened. This is done using the socket() system call:
    fmt.Println("Hello, this is netlink socket")
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

    fmt.Println("unix.NETLINK_KOBJECT_UEVENT:", unix.NETLINK_KOBJECT_UEVENT)
    fmt.Println("fd is:", fd)

    if socket_err != nil {
        fmt.Println("Socket_err is:" + socket_err.Error())
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
        fmt.Println("Bind_err is:" + bind_err.Error())
    }

    bstore := make([]byte, os.Getpagesize())
    fmt.Println("length b is:", len(bstore))
    for {
        fmt.Println("Start reading2 ...")
        for {
            // Peek at the buffer to see how many bytes are available.
            n, _, _ := unix.Recvfrom(fd, bstore, 4096)
            fmt.Println("Length of data is:", n)
            // http://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
            fmt.Println(string(bstore[:n]))
            // Break when we can read all messages.
            if n < len(bstore) {
                BatteryStatusHandle(bstore[:n])
                break
            }
            // Double in size if not enough bytes.
            bstore = make([]byte, len(bstore)*2)
        }
        fmt.Println("Start reading3 ...")
    }

}
