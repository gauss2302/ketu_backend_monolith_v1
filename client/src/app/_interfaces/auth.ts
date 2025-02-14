// app/_interfaces/auth.ts

// --- Request Interfaces ---
export interface ILoginRequestDTO {
  email: string;
  password: string;
}

export interface IRegisterRequestDTO {
  username: string;
  name: string;
  email: string;
  password: string;
}

export interface IOwnerLoginRequestDTO {
  email: string;
  password: string;
}

export interface IOwnerRegisterRequestDTO {
  name: string;
  email: string;
  phone: string;
  password: string;
}

// --- Response Interfaces ---
export interface IUser {
  id: number;
  username: string;
  email: string;
  name: string;
}

export interface IOwner {
  owner_id: number;
  name: string;
  email: string;
  phone: string;
  created_at: string; // Or Date
}

export interface IAuthResponse {
  user: IUser | null;
  accessToken: string;
  expiresIn: number;
}

export interface IOwnerAuthResponse {
  owner: IOwner | null;
  access_token: string;
  expires_in: number;
}

export interface ITokenRefreshResponse {
  accessToken: string;
  expiresIn: number;
}
