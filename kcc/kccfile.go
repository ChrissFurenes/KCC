package kcc

func NewKccFile(CurrentFilePath string) *Item{
	k := &Item{
		Path: KccPath(),
	}
	k.KccInit()
	return k
}
func (i *Item) KccInit(){

}