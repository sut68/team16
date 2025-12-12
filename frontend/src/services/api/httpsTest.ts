import axios, { AxiosError } from "axios";

const API_URL = import.meta.env.VITE_API_URL || "/api";

const getCookie = (name: string): string | null => {
  const match = document.cookie.match(new RegExp('(?:^|; )' + name + '=([^;]*)'));
  if (!match || match[1] === undefined) return null;
  return decodeURIComponent(match[1]);
};

export const getToken = (): string | null => {
  return sessionStorage.getItem("token") || localStorage.getItem("token") || null
}

const getTokenType = (): string => (sessionStorage.getItem("token_type") || localStorage.getItem("token_type") || "Bearer");

export const axiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

axios.defaults.withCredentials = true;

axiosInstance.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers = config.headers || {};
    const ttype = getTokenType();
    config.headers["Authorization"] = `${ttype} ${token}`;
  }

  const method = (config.method || "").toLowerCase();
  if (["post", "put", "patch", "delete"].includes(method)) {
    const csrf = getCookie("csrf_token");
    if (csrf) {
      config.headers = config.headers || {};
      config.headers["X-CSRF-Token"] = csrf;
    }
  }
  return config;
}, (error) => Promise.reject(error));

const getConfigWithoutAuth = () => ({
  headers: {
    "Content-Type": "application/json",
  },
});


export const Post = async (url: string, data: any): Promise<any> => {
  try {
    const res = await axiosInstance.post(url, data, getConfigWithoutAuth());
    return res.data;
  } catch (error: any) {
    handleAuthError(error);
    return error.response;
  }
};

export const Get = async (url: string): Promise<any> => {
  try {
    const res = await axiosInstance.get(url, getConfigWithoutAuth());
    return res.data;
  } catch (error: any) {
    handleAuthError(error);
    return error.response;
  }
};

export const Put = async (url: string, data: any): Promise<any> => {
  try {
    const res = await axiosInstance.put(url, data, getConfigWithoutAuth());
    return res.data;
  } catch (error: any) {
    handleAuthError(error);
    return error.response;
  }
};

export const Delete = async (url: string): Promise<any> => {
  try {
    const res = await axiosInstance.delete(url, getConfigWithoutAuth());
    return res.data;
  } catch (error: any) {
    handleAuthError(error);
    return error.response;
  }
};

function handleAuthError(error: AxiosError | any) {
  if (error?.response?.status === 401) {
    try { sessionStorage.clear(); } catch {}
    try { localStorage.clear(); } catch {}
    window.location.reload();
  }
  return;
}

export const httpsTest = {
  get: axiosInstance.get,
  post: axiosInstance.post,
  put: axiosInstance.put,
  delete: axiosInstance.delete,
  patch: axiosInstance.patch,
};