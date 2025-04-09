// In app/_components/AuthContext.tsx

// In app/_components/AuthContext.tsx

"use client";
import { createContext, useContext } from "react";
import {
  IUser,
  IOwner,
  ILoginRequestDTO,
  IRegisterRequestDTO,
  IAuthResponse,
  IOwnerLoginRequestDTO,
  IOwnerAuthResponse,
  IOwnerRegisterRequestDTO,
} from "../_interfaces/auth";

// In AuthContext.tsx

interface AuthContextType {
  user: IUser | null;
  owner: IOwner | null;
  accessToken: string | null;
  ownerAccessToken: string | null;
  loading: boolean;
  login: (data: ILoginRequestDTO) => Promise<IAuthResponse>;
  register: (data: IRegisterRequestDTO) => Promise<IAuthResponse>;
  ownerLogin: (data: IOwnerLoginRequestDTO) => Promise<IOwnerAuthResponse>;
  ownerRegister: (
    data: IOwnerRegisterRequestDTO
  ) => Promise<IOwnerAuthResponse>;
  logout: () => void;
  ownerLogout: () => void;
  isAuthenticated: () => boolean;
  refreshToken: () => Promise<void>;
  refreshOwnerToken: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
