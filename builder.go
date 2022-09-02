package main

import "fmt"

type breakfast struct {
	coffee bool
	eggs   bool
	bacon  bool
	extra  string
}

type breakfastBuilder struct {
	bfast breakfast
}

func NewBreakfastBuilder() *breakfastBuilder {
	return &breakfastBuilder{
		bfast: breakfast{},
	}
}

func (builder *breakfastBuilder) addCoffee() *breakfastBuilder {
	fmt.Println("adding coffee to breakfast")
	builder.bfast.coffee = true
	return builder
}

func (builder *breakfastBuilder) addBacon() *breakfastBuilder {
	fmt.Println("adding bacon to breakfast")
	builder.bfast.bacon = true
	return builder
}

func (builder *breakfastBuilder) addEggs() *breakfastBuilder {
	fmt.Println("adding eggs to breakfast")
	builder.bfast.eggs = true
	return builder
}

func (builder *breakfastBuilder) addExtra(extra string) *breakfastBuilder {
	fmt.Println(fmt.Sprintf("adding '%s' as extra to breakfast", extra))
	builder.bfast.extra = extra
	return builder
}

func (builder *breakfastBuilder) build() breakfast {
	fmt.Println("breakfast is ready!")
	return builder.bfast
}

func main() {
	_ = NewBreakfastBuilder().
		addCoffee().
		addEggs().
		addBacon().
		addExtra("bagel").
		build()
}
