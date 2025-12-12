import { axiosInstance } from "./httpsTest";
import type { LoginUserRequest, CreateUserRequest } from "../../interfaces/user";

const login = async (credentials: LoginUserRequest) => {
  // backend จะเซ็ต cookies: access_token (HttpOnly) + csrf_token
  const res = await axiosInstance.post("/login", credentials, {
    headers: { "Content-Type": "application/json" },
    withCredentials: true,
  });

  return res.data;
};

const signUp = async (userData: CreateUserRequest) => {
  const res = await axiosInstance.post("/register", userData, {
    headers: { "Content-Type": "application/json" },
    withCredentials: true,
  });

  return res.data;
};

export const authAPITest = {
  login,
  signUp,
};
