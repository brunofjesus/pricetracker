package nulltype

import "database/sql"

func Float64PtrToSqlNullFloat64(p *float64) sql.NullFloat64 {
	if p == nil {
		return sql.NullFloat64{}
	}

	return sql.NullFloat64{
		Float64: *p,
		Valid:   true,
	}
}

func StringToSqlNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func BooleanPrtToSqlNullBool(p *bool) sql.NullBool {
	if p == nil {
		return sql.NullBool{}
	}

	return sql.NullBool{
		Bool:  *p,
		Valid: true,
	}
}

func IntPtrToSqlNullInt32(p *int) sql.NullInt32 {
	if p == nil {
		return sql.NullInt32{}
	}

	return sql.NullInt32{
		Int32: int32(*p),
		Valid: true,
	}
}
