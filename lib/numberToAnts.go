package lib

import "fmt"

func numberToAnts(numberOfAnts int, field *Field) error {
	if numberOfAnts < 1 {
		return fmt.Errorf("invalid number of ants")
	}

	for i := 0; i < numberOfAnts; i++ {
		ant := Ant{}
		ant.ID = i
		ant.IsFinished = false
		field.Ants = append(field.Ants, &ant)
	}

	return nil
}
