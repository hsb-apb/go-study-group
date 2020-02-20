package chapter3

// 課題1
// 以下のstructにgetterとsetterを実装してください。
// Getterの関数名ID, Name
// Setterの関数名SetID, SetName
type Kadai1 struct {
	id   int
	name string
}

func (k *Kadai1) Name() string {
	return k.name
}
func (k *Kadai1) SetName(name string) {
	k.name = name
}

func (k *Kadai1) ID() int {
	return k.id
}
func (k *Kadai1) SetID(id int) {
	k.id = id
}
