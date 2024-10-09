package main

type Car struct {
	Model string
	Color string
	Year  int
	// Otros atributos...
}

// Creamos nuestra interfaz CarBuilder
type CarBuilder interface {
	SetModel(model string) CarBuilder
	SetColor(color string) CarBuilder
	SetYear(year int) CarBuilder
	Build() Car
}

// Implementamos un tipo concreto
type ConcreteCarBuilder struct {
	car Car
}

func NewCarBuilder() ConcreteCarBuilder {
	return ConcreteCarBuilder{}
}

func (b ConcreteCarBuilder) SetModel(model string) CarBuilder {
	b.car.Model = model
	return b
}

func (b ConcreteCarBuilder) SetColor(color string) CarBuilder {
	b.car.Color = color
	return b
}

func (b ConcreteCarBuilder) SetYear(year int) CarBuilder {
	b.car.Year = year
	return b
}

func (b ConcreteCarBuilder) Build() Car {
	return b.car
}

// Opcionalmente podemos apoyarnos del patron
// director para construir un tipo de carro
type CarDirector struct {
	builder CarBuilder
}

func NewCarDirector(builder CarBuilder) CarDirector {
	return CarDirector{builder: builder}
}

func (d CarDirector) ConstructSportsCar() Car {
	return d.builder.SetModel("Sports Model").
		SetColor("Red").
		SetYear(2023).
		Build()
}

// Finalmente podemos usar el director para construir
// nuestro deportivo
func main() {
	builder := NewCarBuilder()
	director := NewCarDirector(builder)

	sportsCar := director.ConstructSportsCar()
	println("Model:", sportsCar.Model)
	println("Color:", sportsCar.Color)
	println("Year:", sportsCar.Year)
}
