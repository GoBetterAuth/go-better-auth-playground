import { cookies } from "next/headers";

import { createClient } from "go-better-auth";
import {
  CSRFPlugin,
  EmailPasswordPlugin,
  OAuth2Plugin,
} from "go-better-auth/plugins";

import { ENV_CONFIG } from "@/constants/env-config";

export const goBetterAuthClientServer = createClient({
  url: ENV_CONFIG.gobetterauth.url,
  plugins: [
    new EmailPasswordPlugin(),
    new OAuth2Plugin(),
    new CSRFPlugin({
      cookieName: "gobetterauth_csrf_token",
      headerName: "x-gobetterauth-csrf-token",
    }),
  ],
  cookies,
});
