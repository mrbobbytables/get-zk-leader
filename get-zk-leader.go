package main

import (
    "flag"
    "fmt"
    "log"
    "regexp"
    "strconv"
    "strings"
    "time"
    "github.com/gogo/protobuf/proto"
     mesos "github.com/mesos/mesos-go/mesosproto"
    "github.com/samuel/go-zookeeper/zk"
)

var (
    zkServer = flag.String("zk", "", "zookeeper1[:port1],zookeeper[:port2]...")
    nodePath = flag.String("path", "", "/mesos, /chronos/state/candidate etc.")
    getMesosLeader = flag.Bool("mesos", false  , "Default to false, set to True if getting mesos leader.")
    debugMode = flag.Bool("debug", false, "Enable 'panic' for errors.")
)


func main() {

    flag.Parse()


    if *zkServer == "" || len(strings.Split(*zkServer,",")) == 0 {
        log.Fatalf("[ -zk ] - No server(s) defined.")
    }

    zkServer_arr := strings.Split(*zkServer, ",")

    if *nodePath == "" {
        log.Fatalf("[ -path ] - No path specified.")
    }

    conn := connect(zkServer_arr)

    children, _, err := conn.Children(*nodePath)
    tryHard(err, "ERROR: Either the Path could not be validated, or no Children could be enumerated.")


    leaderPath := getLeaderZnode(*nodePath, children)
    data, _, err  := conn.Get(leaderPath)
    tryHard(err, "ERROR: Could not get data for znode")

    if *getMesosLeader == true  {
        leaderInfo := new(mesos.MasterInfo)
        err = proto.Unmarshal(data, leaderInfo)
        tryHard(err, "ERROR: Problem encountered while unmarshing MasterInfo from Zookeeper.")
        rawIp := leaderInfo.GetIp()
        fmt.Printf("%d.%d.%d.%d:%d\n", byte(rawIp), byte(rawIp>>8), byte(rawIp>>16), byte(rawIp>>24), leaderInfo.GetPort())
    } else {
        fmt.Printf("%+v\n", string(data))
    }
    conn.Close()
}

func tryHard(err error, errorText string) {
    if err != nil {
        if *debugMode {
            log.Print(errorText)
	    log.Panic(err)
        } else {
            log.Fatalf(errorText)
        }
    }
}

func connect(zks []string) *zk.Conn {
    conn, _, err := zk.Connect(zks, time.Second*5)
    tryHard(err, "ERROR: Could not connect to Zookeeper")
    return conn
}


func getLeaderZnode(path string, children []string) string {

//    zookeeper ephemeral znodes all end with a 10 digit#, the lowest # tends to be the leader.
//    A sequenence # can also never be larger than 2147483647. See:
//    http://zookeeper.apache.org/doc/r3.2.1/zookeeperProgrammers.html#Sequence+Nodes+--+Unique+Naming

    seqRegExp, _ := regexp.Compile("[0-9]{10}$")
    leaderSeqNum := 2147483647
    leaderName := ""
    for _, name := range children {
        if seqRegExp.MatchString(name) {
            nodeSeqNum, _ := strconv.Atoi(seqRegExp.FindString(name))
            if nodeSeqNum < leaderSeqNum {
                leaderSeqNum = nodeSeqNum
                leaderName = name
            }
        }
    }

    leaderPath := path + "/" + leaderName
    return leaderPath
}
