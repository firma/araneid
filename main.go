package main

import (
	"github.com/astaxie/beego"
	_ "github.com/beatrice950201/araneid/extend/begin"
	_ "github.com/beatrice950201/araneid/routers"
)

func main() {
	beego.Run()
}
