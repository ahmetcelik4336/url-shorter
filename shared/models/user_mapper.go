package dto

import "2/ent"

func ToUserResponse(p *ent.User) *UserResponse {

	response := &UserResponse{
		Id:    p.ID,
		Name:  p.Ad,
		Email: p.Email,
	}

	/*if p.Edges.User != nil {
		response.UserID = p.Edges.User.ID
	}*/

	return response
}

func ToLoginResponse(token string) *LoginResponse {

	response := &LoginResponse{
		Token: token,
	}

	return response
}

func ToRegisterResponse(status bool, message string) *RegisterResponse {

	response := &RegisterResponse{
		Status:  status,
		Message: message,
	}

	return response
}

func ToUserResponseList(posts []*ent.User) []*UserResponse {

	var list []*UserResponse

	for _, p := range posts {
		list = append(list, ToUserResponse(p))
	}

	return list
}
