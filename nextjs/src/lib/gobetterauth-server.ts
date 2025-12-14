import { cookies } from 'next/headers';

import { OAuth2ProviderType } from '@/models';
import { ENV_CONFIG } from './env-config';

async function applySetCookie(response: Response): Promise<void> {
  const cookieStore = await cookies();
  const setCookieHeaders = response.headers.getSetCookie();

  if (!setCookieHeaders || setCookieHeaders.length === 0) {
    return;
  }

  setCookieHeaders.forEach((cookieString) => {
    const [nameValue, ...attributes] = cookieString.split(";");
    const [name, value] = nameValue.trim().split("=");

    if (!name || !value) {
      console.warn("[WARN] Invalid cookie format:", nameValue);
      return;
    }

    const options: Parameters<typeof cookieStore.set>[2] = {};

    attributes.forEach((attr) => {
      const trimmedAttr = attr.trim();
      const [attrName, ...attrValueParts] = trimmedAttr.split("=");
      const attrValue = attrValueParts.join("=");

      switch (attrName.toLowerCase()) {
        case "path":
          options.path = attrValue;
          break;
        case "max-age":
          options.maxAge = parseInt(attrValue, 10);
          break;
        case "expires":
          options.expires = new Date(attrValue);
          break;
        case "httponly":
          options.httpOnly = true;
          break;
        case "secure":
          options.secure = true;
          break;
        case "domain":
          options.domain = attrValue;
          break;
        case "samesite":
          options.sameSite = attrValue as "lax" | "strict" | "none";
          break;
      }
    });

    cookieStore.set(name.trim(), decodeURIComponent(value.trim()), options);
  });
}

async function getCSRFCookieValue(): Promise<string> {
  const cookieStore = await cookies();
  const csrf = cookieStore.get("gobetterauth_csrf");
  return csrf?.value ?? "";
}

async function wrappedFetch(
  endpoint: string,
  options: {
    method: "GET" | "POST";
    body?: Record<string, unknown>;
    includeCookies?: boolean;
    includeCSRF?: boolean;
    callbackUrl?: string;
    applySetCookie?: boolean;
  } = { method: "GET", applySetCookie: false },
): Promise<unknown> {
  const url = `${ENV_CONFIG.gobetterauth.url}${endpoint}`;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (options.includeCookies) {
    const cookieStore = await cookies();
    const allCookies = cookieStore
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join("; ");
    headers.Cookie = allCookies;
  }

  if (options.includeCSRF) {
    const csrfToken = await getCSRFCookieValue();
    headers["X-GOBETTERAUTH-CSRF-TOKEN"] = csrfToken;
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
    body,
    cache: "no-store",
  });

  if (!response.ok) {
    const data = await response.json();
    throw new Error(data.message || response.statusText);
  }

  if (!!options.applySetCookie) {
    await applySetCookie(response);
  }

  const data = await response.json();
  return data;
}

// Simulating GoBetterAuth Node.js Server-Side SDK (coming soon)
export const goBetterAuthServer = {
  signUp: {
    email: async (
      name: string,
      email: string,
      password: string,
      callbackUrl?: string,
    ) => {
      return await wrappedFetch("/sign-up/email", {
        method: "POST",
        body: { name, email, password },
        callbackUrl,
        applySetCookie: true,
      });
    },
  },
  signIn: {
    email: async (email: string, password: string, callbackUrl?: string) => {
      return await wrappedFetch("/sign-in/email", {
        method: "POST",
        body: { email, password },
        callbackUrl,
        applySetCookie: true,
      });
    },
  },
  social: (provider: OAuth2ProviderType, redirectTo?: string) => {
    const redirectUrl = encodeURIComponent(redirectTo || "/");
    return `${ENV_CONFIG.gobetterauth.url}/oauth2/${provider}/login?redirect_to=${redirectUrl}`;
  },
  sendEmailVerification: async (callbackUrl?: string) => {
    return await wrappedFetch("/email-verification", {
      method: "POST",
      body: {},
      includeCookies: true,
      includeCSRF: true,
      callbackUrl,
    });
  },
  resetPassword: async (email: string, callbackUrl?: string) => {
    return await wrappedFetch("/reset-password", {
      method: "POST",
      body: { email },
      callbackUrl,
    });
  },
  changePassword: async (token: string, newPassword: string) => {
    return await wrappedFetch("/change-password", {
      method: "POST",
      body: { token, new_password: newPassword },
    });
  },
  emailChange: async (newEmail: string, callbackUrl?: string) => {
    return await wrappedFetch("/email-change", {
      method: "POST",
      body: { email: newEmail },
      includeCookies: true,
      callbackUrl,
    });
  },
  signOut: async () => {
    const result = await wrappedFetch("/sign-out", {
      method: "POST",
      includeCookies: true,
      includeCSRF: true,
      applySetCookie: true,
    });

    // Explicitly remove cookies on the frontend (Next.js cookies API does not auto-remove expired cookies)
    // TODO: the sdk should be able to allow the user to specify cookie names as it can be customised so that the correct cookies are removed.
    const cookieStore = await cookies();
    cookieStore.delete("gobetterauth.session_token");
    cookieStore.delete("gobetterauth_csrf");

    return result;
  },
  getSession: async () => {
    try {
      return await wrappedFetch("/me", {
        method: "GET",
        includeCookies: true,
        includeCSRF: true,
      });
    } catch {
      return null;
    }
  },
};
