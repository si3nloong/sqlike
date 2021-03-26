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

// Update :
func Update() UpdateStatement {
	return &UpdateActions{}
}

// DeleteOne :
func DeleteOne() DeleteOneStatement {
	return &DeleteOneActions{}
}

// Delete :
func Delete() DeleteStatement {
	return &DeleteActions{}
}
