// src/services/https.ts

//Bam
import axios from "axios";
import type { AxiosResponse, AxiosError } from "axios";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export const getToken = (): string | null => {
  // ถ้า Login ใช้ sessionStorage ก็ดึงแค่จากที่นั่น
  const token = sessionStorage.getItem("token") || localStorage.getItem("token");

  if (!token) {
    console.warn("Token not found in storage");
    return null;
  }

  return token;
};

const getTokenType = (): string =>
  sessionStorage.getItem("token_type") ||
  localStorage.getItem("token_type") ||
  "Bearer";

// Helper to read cookies
const getCookie = (name: string): string | null => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift() || null;
  return null;
};

// รับ data เข้ามาเพื่อเช็คว่าเป็น FormData หรือไม่
const getConfig = (data?: any) => {
  const token = getToken();
  const headers: Record<string, string> = {};

  if (!(data instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  // Legacy Header Auth
  if (token) headers.Authorization = `${getTokenType()} ${token}`;

  // CSRF Token Header
  const csrfToken = getCookie("csrf_token");
  if (csrfToken) {
    headers["x-CSRF-Token"] = csrfToken;
  }

  return {
    headers,
    withCredentials: true // เพื่อให้ส่ง/รับ Cookie ได้
  };
};

// ทำเหมือนกันกับ getConfigWithoutAuth
const getConfigWithoutAuth = (data?: any) => {
  const headers: Record<string, string> = {};

  if (!(data instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  // CSRF for public endpoints (like login/register if needed)
  const csrfToken = getCookie("csrf_token");
  if (csrfToken) {
    headers["x-CSRF-Token"] = csrfToken;
  }

  return {
    headers,
    withCredentials: true
  };
};

export const Post = async (
  url: string,
  data: any,
  requireAuth: boolean = true
): Promise<AxiosResponse | any> => {
  // ส่ง data เข้าไปเช็ค config
  const config = requireAuth ? getConfig(data) : getConfigWithoutAuth(data);
  return await axios
    .post(`${API_URL}${url}`, data, config)
    .then((res) => res.data)
    .catch((error: AxiosError) => {
      if (error?.response?.status === 401) {
        try { sessionStorage.clear(); } catch { }
        try { localStorage.clear(); } catch { }
        //window.location.reload();
      }
      return error.response;
    });
};

export const Get = async (
  url: string,
  requireAuth: boolean = true
): Promise<AxiosResponse | any> => {
  const config = requireAuth ? getConfig() : getConfigWithoutAuth();
  return await axios
    .get(`${API_URL}${url}`, config)
    .then((res) => res.data)
    .catch((error: AxiosError) => {
      if (error?.message === "Network Error") {
        return error.response;
      }
      if (error?.response?.status === 401) {
        try { sessionStorage.clear(); } catch { }
        try { localStorage.clear(); } catch { }
        window.location.reload();
      }
      return error.response;
    });
};

export const Put = async (
  url: string,
  data: any,
  requireAuth: boolean = true
): Promise<AxiosResponse | any> => {
  // ส่ง data เข้าไปเช็ค config
  const config = requireAuth ? getConfig(data) : getConfigWithoutAuth(data);
  return await axios
    .put(`${API_URL}${url}`, data, config)
    .then((res) => res.data)
    .catch((error: AxiosError) => {
      if (error?.response?.status === 401) {
        try { sessionStorage.clear(); } catch { }
        try { localStorage.clear(); } catch { }
        window.location.reload();
      }
      return error.response;
    });
};

export const Patch = async (
  url: string,
  data: any,
  requireAuth: boolean = true
): Promise<AxiosResponse | any> => {
  const config = requireAuth ? getConfig(data) : getConfigWithoutAuth(data);
  return await axios
    .patch(`${API_URL}${url}`, data, config)
    .then((res) => res.data)
    .catch((error: AxiosError) => {
      if (error?.response?.status === 401) {
        try { sessionStorage.clear(); } catch { }
        try { localStorage.clear(); } catch { }
        window.location.reload();
      }
      return error.response;
    });
}

export const Delete = async (
  url: string,
  requireAuth: boolean = true
): Promise<AxiosResponse | any> => {
  const config = requireAuth ? getConfig() : getConfigWithoutAuth();
  return await axios
    .delete(`${API_URL}${url}`, config)
    .then((res) => res.data)
    .catch((error: AxiosError) => {
      if (error?.response?.status === 401) {
        try { sessionStorage.clear(); } catch { }
        try { localStorage.clear(); } catch { }
        window.location.reload();
      }
      return error.response;
    });
};

export const axiosInstance = axios.create({
  baseURL: API_URL,
  withCredentials: true, // เพื่อให้ส่ง/รับ Cookie ได้
});

// Add CSRF token to all requests
axiosInstance.interceptors.request.use((config) => {
  const csrfToken = getCookie("csrf_token");
  if (csrfToken) {
    config.headers["x-CSRF-Token"] = csrfToken;
  }
  // Add Authorization header if token exists in storage (Legacy support)
  const token = getToken();
  if (token) {
    config.headers.Authorization = `${getTokenType()} ${token}`;
  }
  return config;
});

export const https = {
  get: axiosInstance.get,
  post: axiosInstance.post,
  put: axiosInstance.put,
  delete: axiosInstance.delete,
  patch: axiosInstance.patch,
};