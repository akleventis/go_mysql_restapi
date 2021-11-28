package test

import (
	handlers "go_mysql/src/handlers"
)

var Cat1Mock = &handlers.Animal{
	ID:     1,
	Name:   "Cat1",
	Age:    "1",
	Color:  "White",
	Gender: "Female",
	Breed:  "Munchkin",
	Weight: "8",
}

var Cat2Mock = &handlers.Animal{
	ID:     2,
	Name:   "Cat2",
	Age:    "3",
	Color:  "Orange",
	Gender: "Male",
	Breed:  "Bengal",
	Weight: "4",
}

var Dog1Mock = &handlers.Animal{
	ID:     1,
	Name:   "Dog1",
	Age:    "5",
	Color:  "Brown",
	Gender: "Female",
	Breed:  "Husky",
	Weight: "34",
}

var Dog2Mock = &handlers.Animal{
	ID:     2,
	Name:   "Dog2",
	Age:    "3",
	Color:  "Grey",
	Gender: "Male",
	Breed:  "Boxer",
	Weight: "4",
}
