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

// UpdateMany :
func UpdateMany() UpdateStatement {
	return &UpdateActions{}
}

// Delete :
func Delete() DeleteStatement {
	return &DeleteActions{}
}

// Copy :
func Copy() CopyStatement {
	return &CopyActions{}
}
