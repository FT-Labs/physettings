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
    ROFI_COLOR            = "ROFI_COLOR"
    POWERMENU_TYPE        = "POWERMENU_TYPE"
    POWERMENU_STYLE       = "POWERMENU_STYLE"
    POWERMENU_CONFIRM     = "POWERMENU_CONFIRM"

    PICOM_EXPERIMENTAL    = "PICOM_EXPERIMENTAL"

    POS_MAKE_BAR          = "pOS-make-bar"
    POS_GRUB_CHOOSE_THEME = "pOS-grub-choose-theme"
    POS_SDDM_CHOOSE_THEME = "pOS-sddm-choose-theme"

    PLYMOUTH              = "PLYMOUTH"
)

var settingsPath string
// GLOBAL VARIABLES
var Attrs      map[string]string
var RofiColors []string
var PowerMenuTypes  []string
var PowerMenuStyles []string
var ScriptInfo map[string]string

func appendAttribute(attribute string) error {
    cmd := fmt.Sprintf("echo %s >> %s", attribute, settingsPath)
    return exec.Command("/bin/bash", "-c", cmd).Run()
}

func ChangeAttribute(attribute, value string) {
    cmd := fmt.Sprintf("sed -i '/%s/c\\%s=%s' %s", attribute, attribute, value, settingsPath)
    err := exec.Command("/bin/bash", "-c", cmd).Run()

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

func fetchPowerMenuTypes() ([]string) {
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


func fetchRofiColors() []string {
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

func fetchScriptInfo() map[string]string {
    m := make(map[string]string)
    out, err := exec.Command("manfilter", "phyos", "CUSTOMIZATION", "SCRIPTS").Output()

    if err != nil {
        panic("Can't fetch script data")
    }

    scriptInfo := strings.Split(string(out), ";")

    for i := 0; i < len(scriptInfo) - 1; i += 2 {
        m[scriptInfo[i]] = scriptInfo[i + 1]
    }
    return m
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
    settingsDefaultPath := fmt.Sprintf("%s/.config/phyos/phyos.conf.default", home)
    f, err := os.Open(settingsPath)

    if err != nil {
        panic("Can't open user settings file")
    }

    sc := bufio.NewScanner(f)

    for sc.Scan() {
        l := strings.ReplaceAll(sc.Text(), "\n", "")

        if strings.Contains(l, "=") {
            arr := strings.Split(l, "=")
            if len(arr) == 1 {
                Attrs[arr[0]] = ""
            } else {
                Attrs[arr[0]] = arr[1]
            }
        }
    }
    if err := sc.Err(); err != nil {
        panic(err)
    }
    f.Close()
    f, err = os.Open(settingsDefaultPath)

    sc = bufio.NewScanner(f)

    if err != nil {
        panic("Can't open default settings file")
    }

    for sc.Scan() {
        l := strings.ReplaceAll(sc.Text(), "\n", "")

        if strings.Contains(l, "=") {
            arr := strings.Split(l, "=")
            if _, ok := Attrs[arr[0]]; !ok {
                err := appendAttribute(l)
                if len(arr) > 1 {
                    Attrs[arr[0]] = arr[1]
                }
                if err != nil {
                    fmt.Fprintf(os.Stderr, err.Error())
                }
            }
        }
    }

    f.Close()

    PowerMenuTypes = fetchPowerMenuTypes()
    RofiColors = fetchRofiColors()
    PowerMenuStyles = append(PowerMenuStyles, "style-1", "style-2", "style-3", "style-4", "style-5")
    ScriptInfo = fetchScriptInfo()

    for i := range PowerMenuStyles {
        if PowerMenuStyles[i] == Attrs[POWERMENU_STYLE] {
            PowerMenuStyles[0], PowerMenuStyles[i] = PowerMenuStyles[i], PowerMenuStyles[0]
            break
        }
    }
}
