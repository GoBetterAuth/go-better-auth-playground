import { OAuth2ProviderType } from '@/models';
import { ENV_CONFIG } from './env-config';

function getCSRFCookieValue(): string {
  // Simple client-side cookie parser
  const match = document.cookie.match(/(?:^|; )gobetterauth_csrf=([^;]*)/);
  return match ? decodeURIComponent(match[1]) : "";
}

type WrappedFetchOptions = {
  method: "GET" | "POST";
  body?: Record<string, unknown>;
  callbackUrl?: string;
  includeCSRF?: boolean;
};

async function wrappedFetch(
  endpoint: string,
  options: WrappedFetchOptions,
): Promise<unknown> {
  const url = `${ENV_CONFIG.gobetterauth.url}${endpoint}`;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (options.includeCSRF) {
    const csrf = getCSRFCookieValue();
    if (csrf) {
      headers["X-GOBETTERAUTH-CSRF-TOKEN"] = csrf;
    }
  }

  const bodyData = options.body || {};
  if (options.callbackUrl) {
    bodyData.callback_url = options.callbackUrl;
  }

  const body =
    options.method === "POST" && Object.keys(bodyData).length > 0
      ? JSON.stringify(bodyData)
      : undefined;

  const response = await fetch(url, {
    method: options.method,
    headers,
    credentials: "include",
    body,
  });
  if (response.ok) {
    return await response.json();
  }

  const data = await response.json();
  return Promise.reject(new Error(data.message || response.statusText));
}

// Simulating GoBetterAuth Node.js Client SDK (coming soon)
export const goBetterAuthClient = {
  signUp: {
    email: async (name: string, email: string, password: string) =>
      wrappedFetch("/sign-up/email", {
        method: "POST",
        body: { name, email, password },
      }),
  },
  signIn: {
    email: async (email: string, password: string) =>
      wrappedFetch("/sign-in/email", {
        method: "POST",
        body: { email, password },
      }),
  },
  social: (provider: OAuth2ProviderType, redirectTo?: string) => {
    const redirectUrl = encodeURIComponent(redirectTo || "/");
    return `${ENV_CONFIG.gobetterauth.url}/oauth2/${provider}/login?redirect_to=${redirectUrl}`;
  },
  sendEmailVerification: async (callbackUrl?: string) =>
    wrappedFetch("/email-verification", {
      method: "POST",
      body: {},
      callbackUrl,
      includeCSRF: true,
    }),
  resetPassword: async (email: string, callbackUrl?: string) =>
    wrappedFetch("/reset-password", {
      method: "POST",
      body: { email },
      callbackUrl,
    }),
  changePassword: async (token: string, newPassword: string) =>
    wrappedFetch("/change-password", {
      method: "POST",
      body: { token, new_password: newPassword },
    }),
  emailChange: async (newEmail: string, callbackUrl?: string) =>
    wrappedFetch("/email-change", {
      method: "POST",
      body: { email: newEmail },
      callbackUrl,
    }),
  signOut: async () =>
    wrappedFetch("/sign-out", {
      method: "POST",
      includeCSRF: true,
    }),
  getSession: async () =>
    wrappedFetch("/me", {
      method: "GET",
    }),
};
