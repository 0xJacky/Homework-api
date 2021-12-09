package model

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
)

type TemplateItem struct {
	DataIndex uint `json:"dataIndex"`
	Score     uint `json:"score,string"`
	Type      uint `json:"type"`
}

type Template []TemplateItem

type AnswerSheet map[uint]interface{}

func (a *Assign) Judge() (score uint, err error) {
	var homework Homework
	err = db.First(&homework, a.HomeworkId).Error
	if err != nil {
		return
	}
	var templateBytes []byte
	log.Println(homework.Template)
	templateBytes, err = homework.Template.MarshalJSON()
	if err != nil {
		return
	}
	var template Template
	err = json.Unmarshal(templateBytes, &template)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("template", template)
	var ansBytes []byte
	ansBytes, err = homework.Answer.MarshalJSON()
	if err != nil {
		return
	}
	ans := make(AnswerSheet)
	var assignAnsBytes []byte
	assignAnsBytes, err = a.Answer.MarshalJSON()
	if err != nil {
		return
	}
	err = json.Unmarshal(ansBytes, &ans)
	if err != nil {
		return
	}
	assignAns := make(AnswerSheet)
	err = json.Unmarshal(assignAnsBytes, &assignAns)
	if err != nil {
		return
	}

	for i := range template {
		log.Println("dataIndex", template[i].DataIndex)
		log.Println("score", template[i].Score)
		log.Println("type", template[i].Type)
		switch template[i].Type {
		case 0, 2, 3:
			// 单选，填空，对错判断
			if ans[template[i].DataIndex] == assignAns[template[i].DataIndex] {
				score += template[i].Score
			}
		case 1:
			// 多选
			_ansList := ans[template[i].DataIndex].(map[string]interface{})["checkedList"].([]interface{})
			_assignList := assignAns[template[i].DataIndex].(map[string]interface{})["checkedList"].([]interface{})
			ansList := make([]string, len(_ansList))

			for j, v := range _ansList {
				ansList[j] = fmt.Sprint(v)
			}
			assignList := make([]string, len(_assignList))

			for j, v := range _assignList {
				assignList[j] = fmt.Sprint(v)
			}

			sort.Strings(ansList)
			sort.Strings(assignList)
			ok := true
			for j := range ansList {
				if ansList[j] != ansList[j] {
					ok = false
					break
				}
			}
			if ok {
				score += template[i].Score
			}
		case 4:
			// 主观题
		}

	}
	return
}
