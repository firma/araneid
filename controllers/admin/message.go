package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beatrice950201/araneid/controllers"
	table "github.com/beatrice950201/araneid/extend/func"
	"github.com/beatrice950201/araneid/extend/model/index"
	"gopkg.in/go-playground/validator.v9"
)

/** 新闻中心 **/
type Message struct{ Main }

// @router /message/index [get,post]
func (c *Message) Index() {
	if c.IsAjax() {
		length, _ := c.GetInt("length", table.WebPageSize())
		draw, _ := c.GetInt("draw", 0)
		search := c.GetString("search[value]")
		result := c.PageListItems(length, draw, c.PageNumber(), search)
		result["data"] = c.tableBuilder.Field(true, result["data"].([]orm.ParamsList), c.TableColumnsType(), c.TableButtonsType())
		c.Succeed(&controllers.ResultJson{Data: result})
	}
	dataTableCore := c.tableBuilder.TableOrder([]string{"1", "desc"})
	for _, v := range c.DataTableColumns() {
		dataTableCore = dataTableCore.TableColumns(v["title"].(string), v["name"].(string), v["className"].(string), v["order"].(bool))
	}
	for _, v := range c.DataTableButtons() {
		dataTableCore = dataTableCore.TableButtons(v)
	}
	c.TableColumnsRender(dataTableCore.ColumnsItemsMaps, dataTableCore.OrderItemsMaps, dataTableCore.ButtonsItemsMaps, table.WebPageSize())
}

// @router /message/edit [get,post]
func (c *Message) Edit() {
	id := c.GetMustInt(":id", "非法请求！")
	if c.IsAjax() {
		item, _ := c.Find(id)
		if err := c.ParseForm(&item); err != nil {
			c.Fail(&controllers.ResultJson{Message: "解析错误: " + error.Error(err)})
		}
		if message := c.verifyBase.Begin().Struct(item); message != nil {
			c.Fail(&controllers.ResultJson{
				Message: c.verifyBase.Translate(message.(validator.ValidationErrors)),
			})
		}
		if _, err := orm.NewOrm().Update(&item); err == nil {
			c.Succeed(&controllers.ResultJson{Message: "处理成功", Url: beego.URLFor("Message.Index")})
		} else {
			c.Fail(&controllers.ResultJson{Message: "处理失败，请稍后再试！"})
		}
	}
	c.Data["info"], _ = c.Find(id)
}

// @router /message/status [post]
func (c *Message) Status() {
	status, _ := c.GetInt8(":status", 0)
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.StatusArray(array, status); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "状态更新成功！马上返回中。。。",
			Url:     beego.URLFor("Message.Index"),
		})
	}
}

// @router /message/delete [post]
func (c *Message) Delete() {
	array := c.checkBoxIds(":ids[]", ":ids")
	if errorMessage := c.DeleteArray(array); errorMessage != nil {
		c.Fail(&controllers.ResultJson{
			Message: error.Error(errorMessage),
		})
	} else {
		c.Succeed(&controllers.ResultJson{
			Message: "删除成功！马上返回中。。。",
			Url:     beego.URLFor("Message.Index"),
		})
	}
}

/************************************************表格渲染机制 ************************************************************/

/** 获取一条数据 **/
func (c *Message) Find(id int) (index.Message, error) {
	item := index.Message{
		Id: id,
	}
	return item, orm.NewOrm().Read(&item)
}

/** 批量删除 **/
func (c *Message) DeleteArray(array []int) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Delete(&index.Message{Id: v}); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 更新状态 **/
func (c *Message) StatusArray(array []int, status int8) (e error) {
	_ = orm.NewOrm().Begin()
	for _, v := range array {
		if _, e = orm.NewOrm().Update(&index.Message{Id: v, Status: status}, "Status"); e != nil {
			_ = orm.NewOrm().Rollback()
			break
		}
	}
	if e == nil {
		_ = orm.NewOrm().Commit()
	}
	return e
}

/** 获取需要渲染的Column **/
func (c *Message) DataTableColumns() []map[string]interface{} {
	var maps []map[string]interface{}
	maps = append(maps, map[string]interface{}{"title": "", "name": "_checkbox_", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "ID", "name": "id", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "姓名", "name": "name", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "邮箱", "name": "email", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "标题", "name": "title", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "简介", "name": "context", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "状态", "name": "status", "className": "text-center data_table_btn_style", "order": false})
	maps = append(maps, map[string]interface{}{"title": "时间", "name": "create_time", "className": "text-center", "order": false})
	maps = append(maps, map[string]interface{}{"title": "操作", "name": "button", "className": "text-center data_table_btn_style", "order": false})
	return maps
}

/** 获取需要渲染的按钮组 **/
func (c *Message) DataTableButtons() []*table.TableButtons {
	var array []*table.TableButtons
	array = append(array, &table.TableButtons{
		Text:      "删除",
		ClassName: "btn btn-sm btn-alt-danger mt-1 ids_deletes",
		Attribute: map[string]string{"data-action": beego.URLFor("Message.Delete")},
	})
	array = append(array, &table.TableButtons{
		Text:      "标记",
		ClassName: "btn btn-sm btn-alt-primary mt-1 ids_enables",
		Attribute: map[string]string{"data-action": beego.URLFor("Message.Status"), "data-field": "status"},
	})
	return array
}

/** 处理分页 **/
func (c *Message) PageListItems(length, draw, page int, search string) map[string]interface{} {
	var lists []orm.ParamsList
	qs := orm.NewOrm().QueryTable(new(index.Message))
	if search != "" {
		qs = qs.Filter("title__icontains", search)
	}
	recordsTotal, _ := qs.Count()
	_, _ = qs.Limit(length, length*(page-1)).ValuesList(&lists, "id", "name", "email", "title", "context", "status", "create_time")
	for _, v := range lists {
		v[3] = c.substr2HtmlHref("javascript:void(0)", v[3].(string), 0, 15)
		v[4] = c.substr2HtmlHref("javascript:void(0)", v[4].(string), 0, 15)
		v[5] = c.Int2HtmlStatus(v[5], v[0], beego.URLFor("Message.Status"))
	}
	data := map[string]interface{}{
		"draw":            draw,         // 请求次数
		"recordsFiltered": recordsTotal, // 从多少条里面筛选
		"recordsTotal":    recordsTotal, // 总条数
		"data":            lists,        // 筛选结果
	}
	return data
}

/**  转为pop提示 **/
func (c *Message) substr2HtmlHref(u, s string, start, end int) string {
	html := fmt.Sprintf(`<a href="%s" target="_blank" class="badge badge-primary js-tooltip" data-placement="top" data-toggle="tooltip" data-original-title="%s">%s...</a>`, u, s, beego.Substr(s, start, end))
	return html
}

/** 将状态码转为html **/
func (c *Message) Int2HtmlStatus(val interface{}, id interface{}, url string) string {
	t := table.BuilderTable{}
	return t.AnalysisTypeSwitch(map[string]interface{}{"status": val, "id": id}, "status", url, map[int64]string{0: "待处理", 1: "已处理"})
}

/** 返回表单结构字段如何解析 **/
func (c *Message) TableColumnsType() map[string][]string {
	result := map[string][]string{
		"columns":   {"string", "string", "string", "string", "string", "string", "date"},
		"fieldName": {"id", "name", "email", "title", "context", "status", "create_time"},
		"action":    {"", "", "", "", "", "", ""},
	}
	return result
}

/** 返回右侧按钮数据结构 **/
func (c *Message) TableButtonsType() []*table.TableButtons {
	buttons := []*table.TableButtons{
		{
			Text:      "回复",
			ClassName: "btn btn-sm btn-alt-primary open_iframe",
			Attribute: map[string]string{
				"href":      beego.URLFor("Message.Edit", ":id", "__ID__", ":popup", 1),
				"data-area": "580px,200px",
			},
		},
		{
			Text:      "删除",
			ClassName: "btn btn-sm btn-alt-danger ids_delete",
			Attribute: map[string]string{
				"data-action": beego.URLFor("Message.Delete"),
				"data-ids":    "__ID__",
			},
		},
	}
	return buttons
}
