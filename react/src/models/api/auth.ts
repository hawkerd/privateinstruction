export interface SignUpRequest {
  username: string;
  password: string;
  email: string;
}

export interface SignInRequest {
  username: string;
  password: string;
  email: string;
}

export interface SignInResponse {
  accessToken: string;
}

export interface UpdatePasswordRequest {
  old_password: string;
  new_password: string;
}
