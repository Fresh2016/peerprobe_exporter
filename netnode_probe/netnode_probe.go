package main
import (
        "flag"
        "fmt"
        "bufio"
        "io"
        "io/ioutil"
        "net"
        "os"
        "strings"
        "strconv"
        "time"

)
var (
        GitTag    = "2017.01.01.release"
        BuildTime = "2017-01-01T00:00:00+0800"
)


//func getPeerConnectEntries() (error) {
  //      file, err := os.Open("/etc/hostsmap")
    //    if err != nil {
      //          return err
       // }
       // defer file.Close()

       // err = parsePeerConnectEntries(file)
       // if err != nil {
       //         return err
       // }

      //  return nil
//}
//func parsePeerConnectEntries(data io.Reader) (error) {
func getPeerConnectEntries() (int,int,error){
       hostmapfile, err := os.Open("/etc/hostsmap")
        if err != nil {
                return 0, 0, err
        }
        defer hostmapfile.Close()


        scanner := bufio.NewScanner(hostmapfile)
        //entries := make(map[string]uint32)
        //entries := make(map[string]float64)
        var resultfile    *os.File
//        var err   error
        var resultstr,tempStr string
        var titleStr = "IP               TYPE           PORT        TIMEDELAY       TESTTIME\n" 
        var setDialTimeoutVal = 3
        var intervalCheckRtn = 10
        var obsoleteTimeRtn  = 5
        var  tempScannerText string
            
        for scanner.Scan() {
                  tempScannerText = scanner.Text()
                 if (tempScannerText == "") {
               //    fmt.Println("Leave out the blank line:")
                     continue
                  }
                  columns := strings.Fields(tempScannerText)
                 if(len(columns) == 0){
               // if (columns == "") {
                 //   fmt.Println("Leave out the spacebar and tab key line:")
                     continue
                  }
                // columns := strings.Fields(scanner.Text())
                //if (strings.Contains(columns[0], "::") == false) && (strings.Contains(columns[0], "#") == false) {
                  if ((strings.Contains(columns[0], "IP") == false)&&(strings.Contains(columns[0], "App") == false)&&(strings.Contains(columns[0], "###") == false)) {
                     // if(strings.Contains(columns[0], "TCP") == true){
//                 if (strings.Contains(columns[0], "IP") == false) { 
                    if(strings.Contains(columns[0], "CheckInterval") == true) {
                      tempintervalCheck,err := strconv.Atoi(columns[1])
               //       fmt.Println("CheckInterval;"+columns[1])
               //       fmt.Println("Checkobsolete;"+columns[3])
               //       fmt.Println("CheckDialTimeout:;"+columns[5])
                      if err != nil {
                      return 0,0,err
                      }
                      tempobsolete,err :=strconv.Atoi(columns[3])
                      tempdialtimeout,err := strconv.Atoi(columns[5])
                      setDialTimeoutVal = tempdialtimeout
                      intervalCheckRtn = tempintervalCheck
                      obsoleteTimeRtn  = tempobsolete
                      } else {
                      // fmt.Println("CheckReturn;"+columns[0])
                       resultstr =resultstr+columns[0]+"     "+columns[1]+"     "+columns[2]+"     "
                       columns[0] = columns[0] + ":"+columns[2]
                        t1 := time.Now()
                        tempStr = strconv.Itoa((int)(t1.Unix()))
                        //conn, err := net.DialTimeout("tcp", columns[0], time.Second*3)
                        conn, err := net.DialTimeout(columns[1], columns[0], (time.Second)*(time.Duration)(setDialTimeoutVal))
                        //conn, err := net.DialTimeout(columns[1], columns[0], (time.Second)*"30")
                        t2 := time.Now()
                        d := t2.Sub(t1)
                    
                        if err != nil {
                                resultstr = resultstr + "---     "+tempStr+"\n"
                        } else {
                                resultstr = resultstr +strconv.FormatFloat(d.Seconds(), 'f', -1, 64)+"     "+tempStr+"\n" 
                                conn.Close()
                        }
                    }
                }
        }
               
           // printDialTime := fmt.Sprintf("%d",setDialTimeoutVal)
           // printInterval := fmt.Sprintf("%d",intervalCheckRtn)
           // printObsolete := fmt.Sprintf("%d",obsoleteTimeRtn)
           // fmt.Println("setDialTimeoutVal: " +printDialTime)
           // fmt.Println("intervalCheckRtn: " + printInterval)
           // fmt.Println("obsoleteTimeRtn: " + printObsolete)


         _, err = os.Stat("/etc/proberesultnew") 
     
         if (err == nil) {
            resultfile,err = os.OpenFile("/etc/proberesultnew",os.O_CREATE|os.O_APPEND|os.O_RDWR,0666)
            if(err != nil){
                    return 0, 0, err
            }
          } else {
           if os.IsNotExist(err) {
              resultfile, err = os.Create("/etc/proberesultnew")
             if err != nil{
               return 0, 0, err
              }else {
                io.WriteString(resultfile,titleStr)
             }
          }
         }

        // resultfile,err := os.OpenFile("/etc/proberesultnew",os.O_CREATE|os.O_APPEND|os.O_RDWR,0666)
        // if(err != nil){
        //            return err
         // }

          n, _ := resultfile.Seek(0, os.SEEK_END)
          _, err = resultfile.WriteAt([]byte(resultstr), n)
          defer resultfile.Close() 

        if err = scanner.Err(); err != nil {
                return 0, 0, fmt.Errorf("failed to parse connect info: %s", err)
        }       
        
        return intervalCheckRtn, obsoleteTimeRtn, nil
}

func updateResultFile(iObsoletetime int) (error) {
       //os.Rename("/etc/proberesult", "/etc/proberesult.bak") 
       var resultstr string
       var newtime,timegap int
       file, err := os.Open("/etc/proberesultnew")
        if err != nil {
                return err
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        var detectfinished= false

       // buf, err := ioutil.ReadFile("/etc/proberesult")
       // if err != nil {
                //err
       //         return err
        //}
       // content := string(buf)

       for scanner.Scan() {
             if(detectfinished == false){
                columns := strings.Fields(scanner.Text())
                if(strings.Contains(columns[0], "IP") == false){
                     t1 := time.Now()
                     oldtime, err:= strconv.Atoi(columns[4])
                    if err != nil {
                      return err
                     }
                     newtime = (int)(t1.Unix())
                     timegap =  newtime - oldtime
                       if ((timegap/60)> iObsoletetime){
                       resultstr = resultstr+columns[0]+"     "+columns[1]+"     "+columns[2]+"     "+columns[3]+"     "+columns[4]+"\n"
                       // resultstr = columns[0]+"     "+columns[1]+"     "+columns[2]+"     "+columns[3]+"\n"
                       // newContent := strings.Replace(content, resultstr, "", -1)
                       // content = newContent
                       // fmt.Println("should update:"+resultstr)
                       // fmt.Println("After update:"+newContent)
                     } else {
                       detectfinished = true
                     }
                }
         }
     }
        
        buf, err := ioutil.ReadFile("/etc/proberesultnew")
	if err != nil {
		return err
	}
	content := string(buf)
      //  fmt.Println("should update:"+resultstr)
        newContent := strings.Replace(content, resultstr, "", -1)
        ioutil.WriteFile("/etc/proberesultnew", []byte(newContent), 0)
        return nil
}

func main() {
        version := flag.Bool("v", false, "version")
        flag.Parse()
        if *version {
                fmt.Println("Git Tag: " + GitTag)
                fmt.Println("Build Time: " + BuildTime)
        }

       
        for {
           tempintervalTime,tempObsoleteTime,err:= getPeerConnectEntries()
           if err != nil {
            fmt.Println(err)
            return
           }
           updateResultFile(tempObsoleteTime)
           time.Sleep((time.Duration)(tempintervalTime)*time.Second)
         }
}
