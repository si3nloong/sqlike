package actions

// FindOne :
func FindOne() SelectOneStatement {
	return &FindOneActions{}
}

// Find :
func Find() SelectStatement {
	return &FindActions{}
}

// Paginate :
func Paginate() PaginateStatement {
	return &PaginateActions{}
}

// UpdateOne :
func UpdateOne() UpdateOneStatement {
	return &UpdateOneActions{}
}

// UpdateMany :
func UpdateMany() UpdateStatement {
	return &UpdateActions{}
}

// DeleteOne :
func DeleteOne() DeleteOneStatement {
	return &DeleteOneActions{}
}

// DeleteMany :
func DeleteMany() DeleteStatement {
	return &DeleteActions{}
}

// Copy :
func Copy() CopyStatement {
	return &CopyActions{}
}
