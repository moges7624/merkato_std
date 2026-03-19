package main

//
// func (s *Server) HandleLogin() handlers.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) error {
// 		var req struct {
// 			Email    string `json:"email"`
// 			Password string `json:"password"`
// 		}
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			api.BadRequest(w, "Invalid JSON")
// 			return nil
// 		}
//
// 		user, err := s.userService.Authenticate(r.Context(), req.Email, req.Password)
// 		if err != nil {
// 			// Log but don't leak details
// 			slog.WarnContext(r.Context(), "login failed", "email", req.Email, "err", err)
// 			api.Unauthorized(w)
// 			return nil
// 		}
//
// 		accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
// 		if err != nil {
// 			api.InternalServerError(w)
// 			return nil
// 		}
//
// 		refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
// 		if err != nil {
// 			api.InternalServerError(w)
// 			return nil
// 		}
//
// 		// Store refresh token in DB (with expiration + user_id + revoked flag)
// 		if err := s.tokenRepo.StoreRefreshToken(r.Context(), user.ID, refreshToken, auth.RefreshTokenDuration); err != nil {
// 			api.InternalServerError(w)
// 			return nil
// 		}
//
// 		api.WriteJSON(w, http.StatusOK, auth.TokenPair{
// 			AccessToken:  accessToken,
// 			RefreshToken: refreshToken,
// 		})
// 		return nil
// 	}
// }
