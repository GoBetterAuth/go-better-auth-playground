import { ENV_CONFIG } from './env-config';

// Simulating GoBetterAuth Node.js SDK (coming soon)
export const goBetterAuthClient = {
  signUp: {
    email: async (name: string, email: string, password: string) => {
      const response = await fetch(
        `${ENV_CONFIG.gobetterauth.url}/sign-up/email`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            name,
            email,
            password,
          }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        return data;
      }

      const data = await response.json();
      return Promise.reject(new Error(data.message || response.statusText));
    },
  },
  signIn: {
    email: async (email: string, password: string) => {
      const response = await fetch(
        `${ENV_CONFIG.gobetterauth.url}/sign-in/email`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            email,
            password,
          }),
        }
      );

      if (response.ok) {
        const data = await response.json();
        return data;
      }

      const data = await response.json();
      return Promise.reject(new Error(data.message || response.statusText));
    },
  },
  resetPassword: async (email: string, callbackUrl?: string) => {
    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/reset-password`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email, callback_url: callbackUrl }),
      }
    );

    if (response.ok) {
      const data = await response.json();
      return data;
    }

    const data = await response.json();
    return Promise.reject(new Error(data.message || response.statusText));
  },
  changePassword: async (token: string, newPassword: string) => {
    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/change-password`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ token, new_password: newPassword }),
      }
    );

    if (response.ok) {
      const data = await response.json();
      return data;
    }

    const data = await response.json();
    return Promise.reject(new Error(data.message || response.statusText));
  },
  emailChange: async (newEmail: string, callbackUrl?: string) => {
    const response = await fetch(
      `${ENV_CONFIG.gobetterauth.url}/email-change`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          email: newEmail,
          callback_url: callbackUrl,
        }),
      }
    );

    if (response.ok) {
      const data = await response.json();
      return data;
    }

    const data = await response.json();
    return Promise.reject(new Error(data.message || response.statusText));
  },
  signOut: async () => {
    const response = await fetch(`${ENV_CONFIG.gobetterauth.url}/sign-out`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (response.ok) {
      const data = await response.json();
      return data;
    }

    const data = await response.json();
    return Promise.reject(new Error(data.message || response.statusText));
  },
  getSession: async () => {
    const response = await fetch(`${ENV_CONFIG.gobetterauth.url}/me`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });

    if (response.ok) {
      const data = await response.json();
      return data;
    }

    return null;
  },
};
