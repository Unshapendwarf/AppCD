package main

import (
	"encoding/json"
	"fmt"
	"github.com/rbxorkt12/applink/pkg/argocd"
	"github.com/rbxorkt12/applink/pkg/config"
	appli "github.com/rbxorkt12/applink/pkg/application"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)


func main(){
	firstArg := os.Args[1]
	switch firstArg {
	case "convert" :
		config,err:=Readstdinandunmarshalconfig()
		if err!=nil{
			os.Exit(111)
			log.Println(err)
		}
		apps:=config.ConvertApp()
		for _, app := range apps {
			if os.Args[2]=="auto" {
				app.Meta.Annotations.AppCDoption = "Auto"
				app.Spec.Sync = &appli.Syncpolicy{Automated: &appli.SyncPolicyAutomated{}}
			} else {
				app.Meta.Annotations.AppCDoption = "Manual"
			}
		}
		json_byte,err:=json.MarshalIndent(apps,"","    ")
		if err!=nil{
			os.Exit(111)
			log.Println(err)}
		fmt.Println(string(json_byte))

	case "diff":
		before := os.Args[2]
		after := os.Args[3]
		flag:=Exists(before)
		if flag == false {
			fmt.Errorf("no file %s",before)
			os.Exit(123)
		}
		flag=Exists(after)
		if flag == false {
			fmt.Errorf("no file %s",after)
			os.Exit(123)
		}
		fi,err:=os.Stat(before)
		len:=fi.Size()
		fi2,err:=os.Stat(after)
		len2:=fi2.Size()
		dat, err := ioutil.ReadFile(before)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		var beforeitems []appli.Item
                dat2, err := ioutil.ReadFile(after)
                var afteritems []appli.Item
		if len<=10{
			fmt.Println("sss")
			if len2<=10{
				fmt.Println("both emptry")
				os.Exit(123)
			}else{
				err=json.Unmarshal(dat2,&afteritems)
				writeitems(afteritems,"/diff/CREATE")
				os.Create("/diff/DELETE")
				os.Create("/diff/UPDATE")
				return
			}
		}else{
			if len2<=10{
				err=json.Unmarshal(dat,&beforeitems)
				writeitems(beforeitems,"/diff/DELETE")
				os.Create("/diff/CREATE")
				os.Create("/diff/UPDATE")
				return
			}
		}
		err=json.Unmarshal(dat,&beforeitems)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		err=json.Unmarshal(dat2,&afteritems)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(123)
		}
		create,delete,update:=appli.Appdiff(beforeitems,afteritems)
		writeitems(create,"/diff/CREATE")
		writeitems(delete,"/diff/DELETE")
		writeitems(update,"/diff/UPDATE")
	case "argoinfo" :
		argoinfo,err:=argocd.ArgocdSet(os.Args[2],os.Args[3])
		if err!=nil {
			log.Fatalln(err)
			os.Exit(111)
		}
		byte,err:=json.Marshal(argoinfo)
		if err!=nil {
			log.Fatalln(err)
			os.Exit(111)
		}
		fmt.Println(string(byte))
	default :
		log.Println("That is not implemented")

	}


}

func Readstdinandunmarshalconfig() (config.Appoconfig,error){
	data,err:= ioutil.ReadAll(os.Stdin)
	if err!=nil { return config.Appoconfig{}, err}
	config:= config.Appoconfig{}
	err = yaml.Unmarshal(data,&config)
	if err!=nil{ return config,err}
	return config,nil
}


func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func writeitems(qq []appli.Item,dest string){
	f, err := os.Create(dest)
	defer f.Close()
	if err!=nil {
		log.Fatalln(err)
		os.Exit(123)
	}
	byte,err:=json.MarshalIndent(qq,"","   ")
	if err!=nil {
		log.Fatalln(err)
		os.Exit(123)
	}
	f.WriteString(string(byte))
	
}

