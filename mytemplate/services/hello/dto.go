package hello

func queryDTO(v Query) (res map[string]interface{}) {
	res = make(map[string]interface{})

	res["age = ?"] = v.Age
	res["email != ?"] = v.Email

	return
}
