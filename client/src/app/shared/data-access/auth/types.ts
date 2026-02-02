export interface LoginPayload {
  email: string;
  password: string;
}

export interface LoginResponse {
  data: {
    token: string;
  };
}

export interface RegisterPayload {
  email: string;
  password: string;
  full_name: string;
  organization_name: string;
}

export interface RegisterResponse {
  data: {
    message: string;
  };
}
