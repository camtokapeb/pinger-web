package main

import (
	"fmt"
	"html/template"
)

func InitDefData(s *Session) {

	s.HTML.HTML_Devices = []string{"HUAWEI_MA5800", "ELTEX_MA4000", "HUAWEI_MA5800", "ELTEX_LTP-8X"}

	s.HTML.HTML_DeviceIP = map[string]string{
		"HUAWEI_MA5800": "10.228.14.116",
		"ELTEX_MA4000":  "10.228.200.200",
		"ELTEX_LTP-8X":  "11.228.14.11"}

	s.HTML.HTML_Slots = []string{"1", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
	s.HTML.HTML_Ports = []string{"15", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}
	s.HTML.HTML_Onts = []string{"30"}
	s.HTML.HTML_VlanTr069s = []string{"111", "68", "69", "111"}
	s.HTML.HTML_VlanPPPoEs = []string{"3506", "3502", "3506"}
	s.HTML.HTML_VlanIPTVs = []string{"3526", "3526", "999"}
	s.HTML.HTML_VlanIMSs = []string{"3542", "3542", "789"}
	s.HTML.HTML_VlanvIMSs = []string{"99", "99", "777"}

}

func HtmlSelect(items []string) string {

	var slect_item_html string
	if len(items) > 1 {
		for i := 1; i < len(items); i++ {
			if items[0] == items[i] {
				slect_item_html += fmt.Sprintf("<option selected value=%s>%s</option>", items[i], items[i])
			} else {
				slect_item_html += fmt.Sprintf("<option value=%s>%s</option>", items[i], items[i])
			}
		}
	} else {
		// Длина списка значений по умолчанию минимум 2, но если длина списка 1, то это простое поле input
		slect_item_html += items[0]
	}
	return slect_item_html
}

func UpdateHTML(s *Session) {
	s.HTML.HTML_Device = template.HTML(HtmlSelect(s.HTML.HTML_Devices))
	s.HTML.HTML_Slot = template.HTML(HtmlSelect(s.HTML.HTML_Slots))
	s.HTML.HTML_Port = template.HTML(HtmlSelect(s.HTML.HTML_Ports))
	s.HTML.HTML_Ont = template.HTML(HtmlSelect(s.HTML.HTML_Onts))
	s.HTML.HTML_VlanTr069 = template.HTML(HtmlSelect(s.HTML.HTML_VlanTr069s))
	s.HTML.HTML_VlanPPPoE = template.HTML(HtmlSelect(s.HTML.HTML_VlanPPPoEs))
	s.HTML.HTML_VlanIPTV = template.HTML(HtmlSelect(s.HTML.HTML_VlanIPTVs))
	s.HTML.HTML_VlanIMS = template.HTML(HtmlSelect(s.HTML.HTML_VlanIMSs))
}
