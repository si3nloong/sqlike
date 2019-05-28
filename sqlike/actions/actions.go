package actions

// FindOne :
func FindOne() SelectOneStatement {
	return &FindOneActions{}
}

// Find :
func Find() SelectStatement {
	return &FindActions{}
}

// UpdateOne :
func UpdateOne() UpdateOneStatement {
	return &UpdateOneActions{}
}

// Update :
func Update() UpdateStatement {
	return &UpdateActions{}
}

// Delete :
func Delete() DeleteStatement {
	return &DeleteActions{}
}
