export const authBaseUrl =
  process.env.NEXT_PUBLIC_AUTH_URL ?? "http://localhost:8080";

export const googleLoginUrl = `${authBaseUrl}/auth/google/login`;
