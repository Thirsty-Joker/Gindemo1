package service

import (
	"demo1/src/model"
	"github.com/spf13/cast"
	"math/rand"
	"strings"
)

func Hit(form *model.Client) (string, string, string, string, string) {
	//todo test
	//if hit,
	//return downloadUrl, updateVersionCode, md5, title, updateTips

	var hitRule *model.Rule

	//var rules *[]model.Rule
	rules := model.GetRules(form.DeviceId)
	
	//order by updateVersionCode, use quickSort O(nlogn)
	quickSort(rules, 0, len(*rules)-1)
	
	//对于 rule, 显然 rule.UpdateVersionCode > rule.MaxUpdateVersionCode >= rule.MinUpdateVersionCode
	//按照常理来说，对 rule.UpdateVersionCode 进行排序之后，rule.MaxUpdateVersionCode 和 rule.MinUpdateVersionCode 也是有序的
	// 所以进行二分，找到 form.UpdateVersionCode 符合的 rule 子集
	left, right := binarySearchLeft(rules, form), binarySearchRight(rules, form)
	
	//由于 rules 是按照从小到大的顺序排列，而有多个匹配的规则时，应该返回 UpdateVersion 最大的（即返回最新版本的安装包链接）
	//所以需要 从 right 到 left 进行匹配，取第一个命中的 rule
	for i := right; i >= left; i-- {
		if matchRule(&(*rules)[i], form) {
			hitRule = &(*rules)[i]
			break
		}
	}
	
	return getDownloadInfo(hitRule)
}

func binarySearchRight(rules *[]model.Rule, form *model.Client) int {
	l, r := 0, len(*rules)-1
	for l < r{
		m := l + ((r-l) >> 1 + 1) // 很怪，>> 比 r-l 先执行？？？
		//if form.UpdateVersionCode < (*rules)[m].MinUpdateVersionCode
		if compareUpdateVersionCode(form.UpdateVersionCode, (*rules)[m].MinUpdateVersionCode) == -1 {
			r = m-1
		} else {
			l = m
		}
	}
	return l
}


func binarySearchLeft(rules *[]model.Rule, form *model.Client) int {
	l, r := 0, len(*rules)-1
	for l < r{
		m := l + ((r-l) >> 1) // 很怪，>> 比 r-l 先执行？？？
		//if form.UpdateVersionCode > (*rules)[m].MaxUpdateVersionCode {
		if compareUpdateVersionCode(form.UpdateVersionCode, (*rules)[m].MaxUpdateVersionCode) == 1{
			l = m+1
		} else {
			r = m
		}
	}
	return l
}


func compareUpdateVersionCode(UpdateVersionCode1, UpdateVersionCode2 string) int {
	/*
	比较 UpdateVersionCode 大小的函数，UpdateVersionCode1 和 UpdateVersionCode2 格式相同(eg. "1.1.1.1")
	返回值解析 : 
	-1 : UpdateVersionCode1 < UpdateVersionCode2
	0 : UpdateVersionCode1 == UpdateVersionCode2
	1 : UpdateVersionCode1 > UpdateVersionCode2
	 */
	var res int = 0
	lis1 := strings.Split(UpdateVersionCode1, ".")
	lis2 := strings.Split(UpdateVersionCode2, ".")
	for index := 0; index < len(lis1); index++ {
		if cast.ToInt(lis1[index]) < cast.ToInt(lis2[index]) {
			res = -1
			break
		}
		if cast.ToInt(lis1[index]) > cast.ToInt(lis2[index]) {
			res = 1
			break
		}
	}
	return res 
}


func quickSort(rules *[]model.Rule, l, r int) {
	if l >= r {
		return
	}

	mid := l
	randIdx := rand.Intn(r-l+1) + l
	(*rules)[l], (*rules)[randIdx] = (*rules)[randIdx], (*rules)[l]
	target := (*rules)[l].UpdateVersionCode
	for i, rule := range (*rules)[l : r+1] {
		if rule.UpdateVersionCode <= target {
			(*rules)[mid], (*rules)[l+i] = (*rules)[l+i], (*rules)[mid]
			mid++
		}
	}
	//左侧 小于等于 target
	//右侧 大于 target
	quickSort(rules, l, mid-1)
	quickSort(rules, mid, r)
	return
}

func matchRule(rule *model.Rule, form *model.Client) bool {
	//model.Client.Version : 请求api版本，⽐如v1/v2
	//model.Client.version_code : 应⽤⼤版本，⽐如8.1.4
	//deviceIdList 白名单，model 里处理
	if form.DevicePlatform != rule.Platform {
		//设备平台
		return false
	}
	if form.Aid != rule.Aid {
		//app 是否相同
		return false
	}
	if (*rule).Channel != (*form).Channel {
		//渠道 是否相同
		return false
	}

	//是否符合 版本要求（应⽤⼩版本，⽐如8.1.4.01），将版本筛选放到前面了，此函数不再需要负责此部分工作了

	if (*form).OsApi < (*rule).MinOsApi || (*form).OsApi > (*rule).MaxOsApi {
		//系统 是否适配
		return false
	}
	if (*form).CpuArch != (*rule).CpuArch {
		return false
	}
	return true
}

func getDownloadInfo(rule *model.Rule) (string, string, string, string, string) {
	if nil == rule {
		return "", "", "", "", ""
	}
	return (*rule).DownloadUrl, (*rule).UpdateTips, (*rule).Md5, (*rule).Title, (*rule).UpdateTips
}
