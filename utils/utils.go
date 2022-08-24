package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

const(
    POWERMENU_TYPE     = "POWERMENU_TYPE"
    POWERMENU_STYLE    = "POWERMENU_STYLE"
    PICOM_EXPERIMENTAL = "PICOM_EXPERIMENTAL"
    ROFI_COLOR         = "ROFI_COLOR"
)

var settingsPath string
// GLOBAL VARIABLES
var Attrs      map[string]string
var RofiColors []string
var PowerMenuTypes  []string
var PowerMenuStyles []string

func ChangeAttribute(attribute, value string) {
    s := fmt.Sprintf("sed -i '/%s/c\\%s=%s' %s", attribute, attribute, value, settingsPath)
    err := exec.Command("/bin/bash", "-c", s).Run()

    if err != nil {
        panic("Error occurred changing attribute")
    }

    if attribute == "CONKY_WIDGETS" {
        exec.Command("/usr/bin/nohup", "pOS-conky").Start()
    } else if attribute == "PICOM_EXPERIMENTAL" {
        exec.Command("killall", "-9", "picom").Run()
        time.Sleep(time.Millisecond * 200)
        exec.Command("/usr/bin/nohup", "pOS-compositor").Start()
    }
}

func FetchPowerMenuTypes() ([]string) {
    f, _ := ioutil.ReadDir("/usr/share/phyos/config/rofi/powermenu")
    types := []string{"Default"}
    for i := 1; i < len(f); i++ {
        types = append(types, fmt.Sprintf("type-%d", i))
        if types[i] == Attrs[POWERMENU_TYPE] {
            types[0], types[i] = types[i], types[0]
        }
    }
    return types
}


func FetchRofiColors() []string {
    var s string
    const path string = "/usr/share/phyos/config/rofi/colors"
    cmd := "file ~/.config/rofi/colors.rasi | tr \"/.\" \" \" | awk '{print $(NF-1)}'"
    out, err := exec.Command("/bin/bash", "-c", cmd).Output()

    if err != nil {
        s = "None"
    } else {
        s = strings.Trim(string(out)," \n")
    }
    cmd = "ls /usr/share/phyos/config/rofi/colors/ | sed -e 's/\\.rasi$//'"
    out, _ = exec.Command("/bin/bash", "-c", cmd).Output()
    colors := strings.Split(string(out), "\n")
    for i:= range colors {
        if colors[i] == s {
            colors[0], colors[i] = colors[i], colors[0]
            break
        }
    }
    return colors[:len(colors) - 1]
}

func SetAttribute(attribute, value string) error {
    if key, ok := Attrs[attribute]; ok {
        ChangeAttribute(attribute, value)
        Attrs[key] = attribute
        return nil
    }
    return errors.New("Can't get attribute")
}

func SetRofiColor(c string) {
    cmd := fmt.Sprintf("ln -sf /usr/share/phyos/config/rofi/colors/%s.rasi ~/.config/rofi/colors.rasi", c)
    exec.Command("/bin/bash", "-c", cmd).Start()
}

func RunScript(c string) error {
    const TERM_TITLE = "physet-run"
    const GEOM       = "80x30"
    err := exec.Command("st", "-n", TERM_TITLE, "-g", GEOM, "-e", c).Run()

    if err != nil {
        return err
    }
    return nil
}

func FetchAttributes() {
    Attrs = make(map[string]string)
    home, err := os.UserHomeDir()

    if err != nil {
        panic(err)
    }
    settingsPath = fmt.Sprintf("%s/.config/phyos/phyos.conf", home)
    f, err := os.Open(settingsPath)

    if err != nil {
        panic("Can't open file")
    }
    defer f.Close()

    sc := bufio.NewScanner(f)

    for sc.Scan() {
        l := strings.ReplaceAll(sc.Text(), "\n", "")

        if strings.Contains(l, "=") {
            arr := strings.Split(l, "=")
            if arr[1] == "" {
                Attrs[arr[0]] = "None"
            } else {
                Attrs[arr[0]] = arr[1]
            }
        }
    }
    if err := sc.Err(); err != nil {
        panic(err)
    }
    PowerMenuTypes = FetchPowerMenuTypes()
    RofiColors = FetchRofiColors()
    PowerMenuStyles = append(PowerMenuStyles, "style-1", "style-2", "style-3", "style-4", "style-5")

    for i := range PowerMenuStyles {
        if PowerMenuStyles[i] == Attrs[POWERMENU_STYLE] {
            PowerMenuStyles[0], PowerMenuStyles[i] = PowerMenuStyles[i], PowerMenuStyles[0]
            break
        }
    }
}