package msg

// ==================
// Implement Eqblame
// ==================

// ToProto exposes the underlying protobuf implementation of EqBlame
func (bl *EquivocationBlame) ToProto() *EquivocationBlame {
	return bl
}

// ==================
// Implement Npblame
// ==================

// ToProto exposes the underlying protobuf implementation of EqBlame
func (bl *NoProgressBlame) ToProto() *NoProgressBlame {
	return bl
}
