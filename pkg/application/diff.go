package app

import "reflect"

func Appdiff(before []Item,after []Item)(create []Item,delete []Item,update []Item){
	var flag int
	for _,app := range before{
		flag = 1
		for _,com := range after {
			if app.Meta.Annotations.Identifier != com.Meta.Annotations.Identifier{
				continue
			}
			if app.Meta.Name != com.Meta.Name {
				continue
			} else {
				flag = 0
				if !deepequalitem(app, com) {
					com.Meta.ResourceVersion = app.Meta.ResourceVersion //update 시에는 resourceversion 명시가 필요
					update = append(update,com)
				}
				break
			}
		}
		if flag ==1 {
			delete = append(delete,app)
		}
	}

	for _,app := range after{
		flag = 1
		for _,com := range before {
			if app.Meta.Annotations.Identifier != com.Meta.Annotations.Identifier{
				continue
			}
			if app.Meta.Name != com.Meta.Name {
				continue
			} else {
				flag = 0
				break
			}
		}
		if flag ==1 {
			create = append(create,app)
		}
	}
	return create,delete,update

}

func deepequalitem(item1 Item, item2 Item) bool{
	if item1.Meta.Annotations != item2.Meta.Annotations { return false}
	if !reflect.DeepEqual(item1.Spec.Dest,item2.Spec.Dest) {return false}
	if !reflect.DeepEqual(item1.Spec.Source,item2.Spec.Source) {return false}
	if !reflect.DeepEqual(item1.Spec.Project,item2.Spec.Project) {return false}
	return true
}