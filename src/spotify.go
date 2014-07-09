package main

import (
    "fmt"
    "os/exec"
    "os"
    "bytes"
    "unsafe"
    "strings"
    "strconv"
)

func execute(keyCommand string) {
    fullCommand := "tell Application \"Spotify\"" + keyCommand
    c := exec.Command("/usr/bin/osascript", "-e", fullCommand)
    defer c.Wait()
    if err := c.Start(); err != nil {
        fmt.Println(keyCommand, "not available")
    }
}

// Volume needs a seperate command because it needs a conversion from
// io.Reader to string. Also we need the out, _ command so I'll seperate it
// for now.
func getVolume() string{
    keyCommand := "to sound volume as integer"
    fullCommand := "tell Application \"Spotify\"" + keyCommand
    c := exec.Command("/usr/bin/osascript", "-e", fullCommand)
    defer c.Wait()
    out, _ := c.StdoutPipe()
    if err := c.Start(); err != nil {
        fmt.Println(keyCommand, "not available")
    }

    // Conversion from io.Reader to string
    buf := new(bytes.Buffer)
    buf.ReadFrom(out)
    b := buf.Bytes()
    s := *(*string)(unsafe.Pointer(&b))
    return s
}

func changeVolume(commands map[string]string, volumeAmount int) {
    words := strings.Fields(getVolume())
    volume := words[0]
    inputValue, err := strconv.Atoi(volume)
    if err != nil {
        fmt.Println(err)
    }
    outputValue := strconv.Itoa(inputValue + volumeAmount)
    execute(format(commands["volumeUp"], outputValue))
}

func format(command, key string) string {
    return fmt.Sprintf(command, key)
}

func main() {
    commands := make(map[string]string)
    commands["play"] = "to play"
    commands["nextTrack"] = "to next track"
    commands["previousTrack"] = "to previous track"
    commands["pause"] = "to pause"
    commands["playPause"] = "to playpause"
    commands["playTrack"] = "to play track \"spotify:track:%s\""
    commands["playPlaylist"] = "to play track \"spotify:user:ni_co:playlist:%s\""
    commands["repeatOn"] = "to set repeating to true"
    commands["repeatOff"] = "to set repeating to false"
    commands["shuffleOn"] = "to set shuffling to true"
    commands["shuffleOff"] = "to set shuffling to false"
    commands["volumeUp"] = "to set sound volume to %s"
    commands["quit"] = "to quit"

    if(len(os.Args) > 1) {
        if(os.Args[1] == "play"){
            if(len(os.Args) == 2){
                execute(commands["play"])
            } else {
                execute(format(commands["playTrack"], os.Args[2]))
            }
        } else if(os.Args[1] == "pause") {
            execute(commands["pause"])
        } else if(os.Args[1] == "playlist"){
            execute(format(commands["playPlaylist"], os.Args[2]))
        } else if(os.Args[1] == "next"){
            execute(commands["nextTrack"])
        } else if(os.Args[1] == "previous"){
            execute(commands["previousTrack"])
        } else if(os.Args[1] == "volume"){
            if(len(os.Args) == 3){
                inputValue, err := strconv.Atoi(os.Args[2])
                if err != nil {
                    fmt.Println(err)
                }
                outputValue := strconv.Itoa(inputValue)
                execute(format(commands["volumeUp"], outputValue))
            } else {
                fmt.Print(getVolume())
            }
        } else if(os.Args[1] == "up") {
            changeVolume(commands, 10)
        } else if(os.Args[1] == "down") {
            changeVolume(commands, -10)
        } else if(os.Args[1] == "shuffle"){
            if(len(os.Args) == 3){
                if(os.Args[2] == "on"){
                    execute(commands["shuffleOn"])
                } else if (os.Args[2] == "off"){
                    execute(commands["shuffleOff"])
                }
            }
        } else if(os.Args[1] == "repeat"){
            if(len(os.Args) == 3){
                if(os.Args[2] == "on"){
                    execute(commands["repeatOn"])
                } else if (os.Args[2] == "off"){
                    execute(commands["repeatOff"])
                }
            }
        } else {
            fmt.Println("Command not found")
        }
    } else {
        fmt.Println("Spotify Options")
        fmt.Println("   play             = Start playing Spotify")
        fmt.Println("   play <uri>       = Start playing specified Spotify URI")
        fmt.Println("   playlist <uri>   = Start playing playlist Spotify URI")
        fmt.Println("   pause            = Pause Spotify")
        fmt.Println("   next             = Play next song")
        fmt.Println("   previous         = Play previous song")
        fmt.Println("   shuffle <on/off> = Shuffle on or off?")
        fmt.Println("   repeat <on/off>  = Repeat on or off?")
        fmt.Println("   volume           = Get volume of Spotify")
        fmt.Println("   volume <amount>  = Set volume by Amount")
        fmt.Println("   up               = Increase volume by 10%")
        fmt.Println("   down             = Decrease volume by 10%")
    }
}