import { Post } from './index';
import type { LoginUserRequest, CreateUserRequest } from '../../interfaces/user';

const login = (credentials: LoginUserRequest) => {
  return Post('/login', credentials, false);
};

const signUp = (userData: CreateUserRequest) => {
  return Post('/register', userData, false);
};

export const authAPI = {
  login,
  signUp,
};
