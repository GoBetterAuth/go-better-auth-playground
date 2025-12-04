import { cookies } from 'next/headers';

import { ENV_CONFIG } from './env-config';

async function applySetCookie(response: Response) {
  const cookieStore = await cookies();

  const setCookieHeader = response.headers.getSetCookie();

  if (setCookieHeader && setCookieHeader.length > 0) {
    setCookieHeader.forEach((cookieString) => {
      const [nameValue, ...attributes] = cookieString.split(";");
      const [name, value] = nameValue.split("=");

      const options: any = {};
      attributes.forEach((attr) => {
        const part = attr.trim();
        if (part.toLowerCase().startsWith("path="))
          options.path = part.split("=")[1];
        if (part.toLowerCase().startsWith("max-age="))
          options.maxAge = parseInt(part.split("=")[1]);
        if (part.toLowerCase().startsWith("expires="))
          options.expires = new Date(part.split("=")[1]);
        if (part.toLowerCase() === "httponly") options.httpOnly = true;
        if (part.toLowerCase() === "secure") options.secure = true;
        if (part.toLowerCase().startsWith("domain="))
          options.domain = part.split("=")[1];
        if (part.toLowerCase().startsWith("samesite="))
          options.sameSite = part.split("=")[1];
      });

      cookieStore.set(name, value, options);
    });
  }
}

export const goBetterAuthServer = {
  signUp: {
    email: async (
      name: string,
      email: string,
      password: string,
      callbackUrl?: string,
    ) => {
      const response = await fetch(
        `${ENV_CONFIG.gobetterauth.url}/sign-up/email`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            name,
            email,
            password,
            callback_url: callbackUrl,
          }),
          cache: "no-store",
        },
      );

      if (response.ok) {
        await applySetCookie(response);
        return await response.json();
      }

      const data = await response.json();
      throw new Error(data.message || response.statusText);
    },
  },
  signIn: {
    email: async (email: string, password: string, callbackUrl?: string) => {
      const response = await fetch(
        `${ENV_CONFIG.gobetterauth.url}/sign-in/email`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password, callback_url: callbackUrl }),
          cache: "no-store",
        },
      );

      if (response.ok) {
        await applySetCookie(response);
        return await response.json();
      }

      const data = await response.json();
      throw new Error(data.message || response.statusText);
    },
  },
  sendEmailVerification: async (callbackUrl?: string) => {
    const cookieStore = await cookies();
    const allCookies = cookieStore
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join("; ");

    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/email-verification`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Cookie: allCookies,
        },
        body: JSON.stringify({
          callback_url: callbackUrl,
        }),
        cache: "no-store",
      },
    );

    if (response.ok) {
      await applySetCookie(response);
      return await response.json();
    }

    const data = await response.json();
    throw new Error(data.message || response.statusText);
  },
  resetPassword: async (email: string, callbackUrl?: string) => {
    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/reset-password`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, callback_url: callbackUrl }),
        cache: "no-store",
      },
    );

    if (response.ok) {
      return await response.json();
    }

    const data = await response.json();
    throw new Error(data.message || response.statusText);
  },
  changePassword: async (token: string, newPassword: string) => {
    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/change-password`,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          token,
          new_password: newPassword,
        }),
        cache: "no-store",
      },
    );

    if (response.ok) {
      return await response.json();
    }

    const data = await response.json();
    throw new Error(data.message || response.statusText);
  },
  emailChange: async (newEmail: string, callbackUrl?: string) => {
    const cookieStore = await cookies();
    const allCookies = cookieStore
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join("; ");

    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/email-change`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Cookie: allCookies,
        },
        body: JSON.stringify({
          email: newEmail,
          callback_url: callbackUrl,
        }),
        cache: "no-store",
      },
    );

    if (response.ok) {
      return await response.json();
    }

    const data = await response.json();
    throw new Error(data.message || response.statusText);
  },
  signOut: async () => {
    const cookieStore = await cookies();
    const allCookies = cookieStore
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join("; ");

    const response = await fetch(`${ENV_CONFIG.gobetterauth.url}/sign-out`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Cookie: allCookies,
      },
      cache: "no-store",
    });

    if (response.ok) {
      await applySetCookie(response);
      return await response.json();
    }

    const data = await response.json();
    throw new Error(data.message || response.statusText);
  },
  getSession: async () => {
    const cookieStore = await cookies();
    const allCookies = cookieStore
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join("; ");

    const response = await fetch(`${ENV_CONFIG.gobetterauth.url}/me`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Cookie: allCookies,
      },
      cache: "no-store",
    });

    if (response.ok) {
      return await response.json();
    }

    return null;
  },
};
