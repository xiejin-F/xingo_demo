package main

type ST struct{
}

func (s *ST)Show(){
    println("ST")
}

func (s *ST)Show2(){
    println("ST:Show2()")
}

type ST2 struct{
    ST
    I int
}

func (s *ST2)Show(){
	s.ST.Show()
    println("ST2")
}

func main() {
    s := ST2{I:5}
    s.Show()
    s.Show2()
    println(s.I)
}