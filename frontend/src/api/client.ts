import axios from "axios";

function getCookie(name: string) {
  const parts = document.cookie.split(";").map((value) => value.trim());
  const match = parts.find((part) => part.startsWith(`${name}=`));
  return match ? decodeURIComponent(match.split("=").slice(1).join("=")) : "";
}

export const api = axios.create({
  baseURL: "/api",
  withCredentials: true
});

api.interceptors.request.use((config) => {
  if (config.method && !["get", "head", "options"].includes(config.method.toLowerCase())) {
    const csrf = getCookie("shieldpanel_csrf");
    if (csrf) {
      config.headers["X-CSRF-Token"] = csrf;
    }
  }
  return config;
});

export async function unwrap<T>(request: Promise<{ data: { data: T } }>): Promise<T> {
  const response = await request;
  return response.data.data;
}

export function getApiErrorMessage(error: unknown): string | null {
  if (axios.isAxiosError(error)) {
    const message = error.response?.data?.message;
    if (typeof message === "string" && message.trim() !== "") {
      return message;
    }
  }
  if (error instanceof Error && error.message.trim() !== "") {
    return error.message;
  }
  return null;
}
